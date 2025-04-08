package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Constants for performance tuning
const (
	ShardCount          = 1024            // Power of 2 for optimal sharding
	EvictionCheckPeriod = 5 * time.Second // How often to check for expired keys
	ConnectionPoolSize  = 100000          // Pre-allocated connection pool size
	MaxConcurrentConns  = 500000          // Maximum concurrent connections
	TCPReadBufferSize   = 4 * 1024        // 4KB read buffer
	TCPWriteBufferSize  = 4 * 1024        // 4KB write buffer
)

// CacheEntry represents a value with its expiration time
type CacheEntry struct {
	Value    string
	ExpireAt int64 // Unix timestamp in nanoseconds
}

// CacheShard represents a single shard of the cache
type CacheShard struct {
	data map[string]CacheEntry
	mu   sync.RWMutex
}

// Cache represents our in-memory key-value store with sharding
type Cache struct {
	shards       []*CacheShard
	shardMask    uint64
	stats        CacheStats
	shutdownChan chan struct{}
}

// CacheStats holds cache statistics for monitoring
type CacheStats struct {
	Gets        uint64
	Sets        uint64
	Deletes     uint64
	Hits        uint64
	Misses      uint64
	Evictions   uint64
	ActiveConns int64
}

// NewCache initializes a new Cache with sharding
func NewCache() *Cache {
	// Ensure ShardCount is a power of 2
	shardMask := uint64(ShardCount - 1)

	cache := &Cache{
		shards:       make([]*CacheShard, ShardCount),
		shardMask:    shardMask,
		shutdownChan: make(chan struct{}),
	}

	// Initialize each shard
	for i := 0; i < ShardCount; i++ {
		cache.shards[i] = &CacheShard{
			data: make(map[string]CacheEntry),
		}
	}

	// Start background eviction worker
	go cache.evictionWorker()

	return cache
}

// getShard returns the appropriate shard for a given key
func (c *Cache) getShard(key string) *CacheShard {
	// FNV-1a hash for even distribution
	var h uint64
	for i := 0; i < len(key); i++ {
		h ^= uint64(key[i])
		h *= 0x100000001b3
	}
	return c.shards[h&c.shardMask]
}

// Set adds a key-value pair to the cache
func (c *Cache) Set(key, value string, ttl time.Duration) {
	shard := c.getShard(key)
	shard.mu.Lock()

	var expireAt int64
	if ttl > 0 {
		expireAt = time.Now().Add(ttl).UnixNano()
	}

	shard.data[key] = CacheEntry{
		Value:    value,
		ExpireAt: expireAt,
	}

	shard.mu.Unlock()
	atomic.AddUint64(&c.stats.Sets, 1)
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (string, bool) {
	shard := c.getShard(key)
	shard.mu.RLock()

	entry, exists := shard.data[key]

	// Check if key exists and is not expired
	if !exists {
		shard.mu.RUnlock()
		atomic.AddUint64(&c.stats.Gets, 1)
		atomic.AddUint64(&c.stats.Misses, 1)
		return "", false
	}

	// Check expiration
	if entry.ExpireAt > 0 && time.Now().UnixNano() > entry.ExpireAt {
		shard.mu.RUnlock()

		// Delete expired key with write lock
		shard.mu.Lock()
		delete(shard.data, key)
		shard.mu.Unlock()

		atomic.AddUint64(&c.stats.Gets, 1)
		atomic.AddUint64(&c.stats.Misses, 1)
		atomic.AddUint64(&c.stats.Evictions, 1)
		return "", false
	}

	value := entry.Value
	shard.mu.RUnlock()

	atomic.AddUint64(&c.stats.Gets, 1)
	atomic.AddUint64(&c.stats.Hits, 1)
	return value, true
}

// Delete removes a key-value pair from the cache
func (c *Cache) Delete(key string) bool {
	shard := c.getShard(key)
	shard.mu.Lock()

	_, exists := shard.data[key]
	if exists {
		delete(shard.data, key)
		shard.mu.Unlock()
		atomic.AddUint64(&c.stats.Deletes, 1)
		return true
	}

	shard.mu.Unlock()
	return false
}

// GetAll returns all non-expired keys and values in the cache
// Note: This is expensive and should be used for UI/admin only
func (c *Cache) GetAll() map[string]string {
	result := make(map[string]string)
	now := time.Now().UnixNano()

	for _, shard := range c.shards {
		shard.mu.RLock()
		for k, entry := range shard.data {
			if entry.ExpireAt == 0 || now < entry.ExpireAt {
				result[k] = entry.Value
			}
		}
		shard.mu.RUnlock()
	}

	return result
}

// evictionWorker periodically checks for and removes expired keys
func (c *Cache) evictionWorker() {
	ticker := time.NewTicker(EvictionCheckPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.evictExpired()
		case <-c.shutdownChan:
			return
		}
	}
}

// evictExpired checks all shards and removes expired keys
func (c *Cache) evictExpired() {
	now := time.Now().UnixNano()
	var evictionCount uint64

	for _, shard := range c.shards {
		var keysToDelete []string

		// First, identify expired keys with read lock
		shard.mu.RLock()
		for k, entry := range shard.data {
			if entry.ExpireAt > 0 && now > entry.ExpireAt {
				keysToDelete = append(keysToDelete, k)
			}
		}
		shard.mu.RUnlock()

		// Then delete them with write lock if any were found
		if len(keysToDelete) > 0 {
			shard.mu.Lock()
			for _, k := range keysToDelete {
				// Double-check expiration before deleting (it might have been updated)
				if entry, exists := shard.data[k]; exists {
					if entry.ExpireAt > 0 && now > entry.ExpireAt {
						delete(shard.data, k)
						evictionCount++
					}
				}
			}
			shard.mu.Unlock()
		}
	}

	if evictionCount > 0 {
		atomic.AddUint64(&c.stats.Evictions, evictionCount)
	}
}

// GetStats returns current cache statistics
func (c *Cache) GetStats() CacheStats {
	return CacheStats{
		Gets:        atomic.LoadUint64(&c.stats.Gets),
		Sets:        atomic.LoadUint64(&c.stats.Sets),
		Deletes:     atomic.LoadUint64(&c.stats.Deletes),
		Hits:        atomic.LoadUint64(&c.stats.Hits),
		Misses:      atomic.LoadUint64(&c.stats.Misses),
		Evictions:   atomic.LoadUint64(&c.stats.Evictions),
		ActiveConns: atomic.LoadInt64(&c.stats.ActiveConns),
	}
}

// Shutdown gracefully shuts down the cache
func (c *Cache) Shutdown() {
	close(c.shutdownChan)
}

// StartTCPServer starts a TCP server on the specified port
func StartTCPServer(cache *Cache, port string) {
	// Set system limits
	// In production, also set ulimit -n to a high value (1M+)

	// Create listener with optimized options
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP server listening on %s", port)

	//TODO: crate lg file for database records

	// Connection limiter using a semaphore
	connLimiter := make(chan struct{}, MaxConcurrentConns)

	// Pre-allocate worker goroutines
	wg := &sync.WaitGroup{}

	for {
		// Block if too many concurrent connections
		connLimiter <- struct{}{}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			<-connLimiter // Release the token
			continue
		}

		// Optimize TCP connection
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetNoDelay(true)   // Disable Nagle's algorithm
			tcpConn.SetKeepAlive(true) // Enable keep-alive
			tcpConn.SetKeepAlivePeriod(30 * time.Second)
		}

		wg.Add(1)
		atomic.AddInt64(&cache.stats.ActiveConns, 1)

		go func() {
			defer func() {
				conn.Close()
				<-connLimiter // Release the token
				atomic.AddInt64(&cache.stats.ActiveConns, -1)
				wg.Done()
			}()

			handleConnection(conn, cache)
		}()
	}
}

// handleConnection processes incoming TCP connections
func handleConnection(conn net.Conn, cache *Cache) {
	reader := bufio.NewReaderSize(conn, TCPReadBufferSize)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		cmd := strings.ToUpper(parts[0])
		switch cmd {
		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
				continue
			}
			key := parts[1]
			value := parts[2]
			var ttl time.Duration
			if len(parts) > 3 && parts[3] == "EX" && len(parts) > 4 {
				seconds, err := time.ParseDuration(parts[4] + "s")
				if err == nil {
					ttl = seconds
				}
			}
			cache.Set(key, value, ttl)
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(parts) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
				continue
			}
			key := parts[1]
			value, exists := cache.Get(key)
			if exists {
				conn.Write([]byte("+" + value + "\r\n"))
			} else {
				conn.Write([]byte("$-1\r\n"))
			}

		case "DEL":
			if len(parts) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'del' command\r\n"))
				continue
			}
			key := parts[1]
			if cache.Delete(key) {
				conn.Write([]byte(":1\r\n"))
			} else {
				conn.Write([]byte(":0\r\n"))
			}

		case "PING":
			conn.Write([]byte("+PONG\r\n"))

		case "QUIT":
			conn.Write([]byte("+OK\r\n"))
			return

		case "INFO":
			stats := cache.GetStats()
			info := fmt.Sprintf("+# Stats\r\n"+
				"gets:%d\r\n"+
				"sets:%d\r\n"+
				"deletes:%d\r\n"+
				"hits:%d\r\n"+
				"misses:%d\r\n"+
				"evictions:%d\r\n"+
				"active_connections:%d\r\n",
				stats.Gets, stats.Sets, stats.Deletes,
				stats.Hits, stats.Misses, stats.Evictions,
				stats.ActiveConns)
			conn.Write([]byte(info))

		default:
			conn.Write([]byte("-ERR unknown command '" + cmd + "'\r\n"))
		}
	}
}

// StartWebServer starts a web server for the GUI
func StartWebServer(cache *Cache, port string) {
	// Use a more efficient HTTP server setup
	server := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := cache.GetAll()
		stats := cache.GetStats()

		// Calculate hit rate
		hitRate := 0.0
		totalOps := stats.Hits + stats.Misses
		if totalOps > 0 {
			hitRate = float64(stats.Hits) / float64(totalOps) * 100
		}

		templateData := struct {
			Data    map[string]string
			Stats   CacheStats
			HitRate float64
		}{
			Data:    data,
			Stats:   stats,
			HitRate: hitRate,
		}

		tmpl := template.Must(template.New("index").Parse(dashboardTemplate))
		tmpl.Execute(w, templateData)
	})

	http.HandleFunc("/api/command", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cmd := r.FormValue("cmd")
		if cmd == "" {
			http.Error(w, `{"status":"error","message":"Empty command"}`, http.StatusBadRequest)
			return
		}

		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			http.Error(w, `{"status":"error","message":"Invalid command"}`, http.StatusBadRequest)
			return
		}

		command := strings.ToUpper(parts[0])
		switch command {
		case "SET":
			if len(parts) < 3 {
				http.Error(w, `{"status":"error","message":"Invalid SET command format"}`, http.StatusBadRequest)
				return
			}
			key := parts[1]
			value := parts[2]
			var ttl time.Duration
			if len(parts) > 3 && parts[3] == "EX" && len(parts) > 4 {
				seconds, err := time.ParseDuration(parts[4] + "s")
				if err == nil {
					ttl = seconds
				}
			}
			cache.Set(key, value, ttl)
			fmt.Fprintf(w, `{"status":"success"}`)

		case "DEL":
			if len(parts) != 2 {
				http.Error(w, `{"status":"error","message":"Invalid DEL command format"}`, http.StatusBadRequest)
				return
			}
			key := parts[1]
			if cache.Delete(key) {
				fmt.Fprintf(w, `{"status":"success"}`)
			} else {
				fmt.Fprintf(w, `{"status":"success","message":"Key not found"}`)
			}

		default:
			http.Error(w, `{"status":"error","message":"Unsupported command"}`, http.StatusBadRequest)
		}
	})

	log.Printf("Web server listening on %s", port)
	log.Fatal(server.ListenAndServe())
}

const dashboardTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>High-Performance Cache Dashboard</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f8f9fa;
            color: #343a40;
        }
        h1, h2, h3 {
            color: #212529;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            padding: 25px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.05);
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-bottom: 30px;
        }
        .stat-card {
            background-color: #f1f3f5;
            padding: 15px;
            border-radius: 6px;
            text-align: center;
        }
        .stat-value {
            font-size: 24px;
            font-weight: bold;
            margin: 10px 0;
            color: #1971c2;
        }
        .stat-label {
            color: #495057;
            font-size: 14px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
            font-size: 14px;
        }
        th, td {
            padding: 12px 15px;
            text-align: left;
            border-bottom: 1px solid #e9ecef;
        }
        th {
            background-color: #f8f9fa;
            font-weight: 600;
        }
        tr:hover {
            background-color: #f1f3f5;
        }
        .form-container {
            margin: 30px 0;
            padding: 20px;
            background-color: #f8f9fa;
            border-radius: 6px;
        }
        .form-row {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
            align-items: flex-end;
        }
        .form-group {
            flex: 1;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: 500;
            font-size: 14px;
        }
        input[type="text"] {
            width: 100%;
            padding: 10px;
            border: 1px solid #ced4da;
            border-radius: 4px;
            font-size: 14px;
        }
        button {
            padding: 10px 15px;
            background-color: #1971c2;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            transition: background-color 0.2s;
        }
        button:hover {
            background-color: #1864ab;
        }
        .delete-btn {
            background-color: #f03e3e;
        }
        .delete-btn:hover {
            background-color: #e03131;
        }
        .refresh-btn {
            background-color: #37b24d;
        }
        .refresh-btn:hover {
            background-color: #2f9e44;
        }
        .actions {
            margin-top: 20px;
            text-align: right;
        }
        .empty-message {
            text-align: center;
            padding: 30px;
            color: #868e96;
            font-style: italic;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>High-Performance Cache Dashboard</h1>
        
        <h2>Cache Statistics</h2>
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-value">{{.Stats.Gets}}</div>
                <div class="stat-label">GET Operations</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{.Stats.Sets}}</div>
                <div class="stat-label">SET Operations</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{.Stats.Deletes}}</div>
                <div class="stat-label">DELETE Operations</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{printf "%.1f" .HitRate}}%</div>
                <div class="stat-label">Hit Rate</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{.Stats.Evictions}}</div>
                <div class="stat-label">Evictions</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{.Stats.ActiveConns}}</div>
                <div class="stat-label">Active Connections</div>
            </div>
        </div>
        
        <div class="form-container">
            <h3>Add/Update Key</h3>
            <form id="cacheForm">
                <div class="form-row">
                    <div class="form-group">
                        <label for="key">Key</label>
                        <input type="text" id="key" placeholder="Key" required>
                    </div>
                    <div class="form-group">
                        <label for="value">Value</label>
                        <input type="text" id="value" placeholder="Value" required>
                    </div>
                    <div class="form-group">
                        <label for="ttl">TTL (seconds, optional)</label>
                        <input type="text" id="ttl" placeholder="e.g. 300">
                    </div>
                    <button type="submit">Set</button>
                </div>
            </form>
        </div>
        
        <h2>Cached Data</h2>
        <table>
            <thead>
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                {{if eq (len .Data) 0}}
                <tr>
                    <td colspan="3" class="empty-message">No data in cache</td>
                </tr>
                {{else}}
                {{range $key, $value := .Data}}
                <tr>
                    <td>{{$key}}</td>
                    <td>{{$value}}</td>
                    <td><button class="delete-btn" onclick="deleteKey('{{$key}}')">Delete</button></td>
                </tr>
                {{end}}
                {{end}}
            </tbody>
        </table>
        
        <div class="actions">
            <button class="refresh-btn" onclick="location.reload()">Refresh</button>
        </div>
    </div>

    <script>
        document.getElementById('cacheForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const key = document.getElementById('key').value;
            const value = document.getElementById('value').value;
            const ttl = document.getElementById('ttl').value;
            
            let cmd = 'SET ' + key + ' ' + value;
            if (ttl) {
                cmd += ' EX ' + ttl;
            }
            
            fetch('/api/command', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: 'cmd=' + encodeURIComponent(cmd)
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    location.reload();
                } else {
                    alert('Error: ' + data.message);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred while processing your request.');
            });
        });
        
        function deleteKey(key) {
            const cmd = 'DEL ' + key;
            
            fetch('/api/command', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: 'cmd=' + encodeURIComponent(cmd)
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    location.reload();
                } else {
                    alert('Error: ' + data.message);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred while processing your request.');
            });
        }
    </script>
</body>
</html>
`

func main() {
	// Set max CPU cores for parallelism
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create cache with optimized shard count
	cache := NewCache()
	/*
	   	// Log startup info
	   	fmt.Printf(`
	   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⡀⠀⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣀⡀⢀⡶⠋⠉⠉⠓⢦⠀⠀⠀
	   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⠏⠉⠈⢻⡞⠀⠀⠀⠀⠀⠈⣧⠀⠀
	   ⠀⠀⠀⠀⠀⠀⠀⢠⠾⠛⠻⡄⠀⠀⠀⣧⠀⠀⠀⠀⠀⠀⢸⠆⠀
	   ⠀⠀⠀⠀⠀⣠⡤⠿⣄⠀⠀⢹⠄⠀⠀⢻⡀⠀⠀⠀⠀⠀⣸⠀⠀
	   ⠀⠀⠀⠀⢰⠃⠀⠀⠸⡆⢠⣸⠤⠤⠤⠸⠧⠤⠤⢄⡀⠈⡇⠀⠀
	   ⠀⠀⠀⢀⣼⣀⣀⣀⣜⡱⠊⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣷⠀⠀
	   ⠀⠀⡼⠉⠈⢳⠤⠐⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⡄⠀
	   ⠀⢰⡇⠀⢠⠎⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣧⠀
	   ⠀⢸⡧⠔⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⠀
	   ⢀⡞⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⠇⠀
	   ⢸⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡾⠀⠀
	   ⢸⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⠃⠀⠀
	   ⠸⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡞⠀⠀⠀
	   ⠀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⠇⠀⠀⠀
	   ⠀⠸⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⠀⠀⠀⠀
	   ⠀⠀⢷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡟⠀⠀⠀⠀
	   ⠀⠀⠈⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀
	   ⠀⠀⠀⢹⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⡇⠀⠀⠀⠀
	   ⠀⠀⠀⠀⢷⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀
	   ⠀⠀⠀⠀⢸⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⡇⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⢸⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⠈⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⠀⠸⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡇⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⠀⠀⠹⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⠁⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⠀⠀⠀⠘⢦⡀⠀⠀⠀⠀⠀⠀⢀⡴⠃⠀⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠳⠤⣤⣤⡤⠴⠋⠁⠀⠀⠀⠀⠀⠀
	   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	   ⠀⠐⠆⣴⠲⢶⠰⢖⠒⠴⡶⠰⢶⠶⠆⠀⠶⠶⠄⠐⡖⠒⠒⠒⠲\t
	   `)

	*/
	fmt.Printf(`
██████╗ ██╗   ██╗███████╗████████╗██████╗ ██████╗ 
██╔══██╗██║   ██║██╔════╝╚══██╔══╝██╔══██╗██╔══██╗
██║  ██║██║   ██║███████╗   ██║   ██║  ██║██████╔╝
██║  ██║██║   ██║╚════██║   ██║   ██║  ██║██╔══██╗
██████╔╝╚██████╔╝███████║   ██║   ██████╔╝██████╔╝
╚═════╝  ╚═════╝ ╚══════╝   ╚═╝   ╚═════╝ ╚═════╝ 
                                                  

`)

	fmt.Printf("🌟 keep foots clean because dust is inevitable\n")

	log.Printf("Starting high-performance cache with %d shards", ShardCount)
	log.Printf("System has %d CPU cores", runtime.NumCPU())

	// Start the TCP server in a goroutine
	go StartTCPServer(cache, ":8989")

	// Start the web server (this will block)
	StartWebServer(cache, ":9090")
}
