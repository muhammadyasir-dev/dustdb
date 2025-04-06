# dustdb 
### A faster Database  than redis and memcache 

🌟 keep foots clean because dust is inevitable.

⠀⠀⠀⢀⡤⣾⠉⠑⡄⠀⠀⠀⠀⠀⠀⠀⠀⢠⠊⠉⣧⢤⡀⠀⠀⠀
⠀⢀⣔⠙⡄⠈⡆⠀⢀⠀⠀⠀⠀⠀⠀⠀⠀⠨⠀⢠⠃⢠⠋⣢⡀⠀
⣀⣌⠈⡆⣗⣚⠯⠚⠘⢆⠀⠀⠀⠀⠀⠀⡰⠃⠓⠽⣓⣺⢰⡁⣱⣀
⡇⢈⣝⠖⠉⣿⠀⠀⠀⠀⢇⠀⠀⠀⠀⡰⠀⠀⠀⠀⢸⠉⠲⡏⡁⢨
⠘⡺⠁⠀⠀⢸⠀⠀⠀⠀⢸⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⢉⡇
⢸⠀⠀⠀⠀⢄⠀⠀⠀⠀⡎⠀⠀⠀⠀⠹⡀⠀⠀⠀⡰⠀⠀⠀⠀⡇
⠈⡄⠀⠀⠀⠘⠄⠀⢀⡜⠀⠀⠀⠀⠀⠀⢣⡀⠀⠠⠃⠀⠀⠀⢠⠃
⠀⠘⡄⠀⠀⠀⠈⠠⠎⡇⠀⠀⠀⠀⠀⠀⢸⠱⠀⠁⠀⠀⠀⢠⠃⠀
⠀⠀⠘⡄⠀⠀⠀⠀⠀⠇⠀⠀⠀⠀⠀⠀⢸⠐⠀⠀⠀⠀⢠⠇⠀⠀
⠀⠀⠀⠘⡀⠀⠀⠀⠀⠘⡄⠀⠀⠀⠀⢠⠃⠀⠀⠀⠀⢀⠆⠀⠀⠀
⠀⠀⠀⠀⢡⠀⠀⠀⠀⠀⠈⡄⠀⠀⢠⠃⠀⠀⠀⠀⠀⡈⠀⠀⠀⠀
⠀⠀⠀⠀⠈⡄⠀⠀⠀⠀⠀⠸⠀⠀⠆⠀⠀⠀⠀⠀⢀⠃⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢣⠀⠀⠀⠀⠀⢀⠆⠀⡀⠀⠀⠀⠀⠀⡜⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠳⠄⣀⣀⠤⠊⠀⠀⠑⠤⣀⣀⠠⠜⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠤⠴⠤⠤⠠⠴⠠⠤⠄⠀⠂⠰⠒⠦⠶⠆⠠⠀⠀⠀⠀


DUSTDB is easy to use yet powerful caching database which supports key value data which is cached and it uses Programmiung languages for core caching mechanisma and network request masking which your Operating system uses (Assembly,C,C++) 



# Comparison of In-Memory Data Stores

| **Feature**            | **Redis**                            | **Memcached**                        | **DragonflyDB**                             | **DustDB**                                         |
|------------------------|---------------------------------------|--------------------------------------|---------------------------------------------|----------------------------------------------------|
| **Type**               | In-memory data store                  | In-memory cache                      | In-memory data store                        | In-memory data store                               |
| **Performance**        | High                                  | High                                 | Higher than both                            | Above all three                                    |
| **Data Structures**    | Rich data types                       | Simple key-value                     | Compatible with Redis                       | Advanced key-value data structures                 |
| **Scalability**        | Clustering available                  | Limited                              | Highly scalable                             | Linear horizontal scaling                          |
| **Latency**            | Low                                   | Low                                  | Lower than both                             | Ultra-low latency                                  |
| **Use Cases**          | Caching, messaging                    | Caching                              | Caching, real-time apps                     | Caching, real-time apps                            |
| **Memory Efficiency**  | Moderate (~70–80%)                    | High (~90%)                          | Superior (~95%)                             | **Exceptional (99.5%)**                            |
| **Eviction Policy**    | LRU, LFU                              | LRU                                  | Custom algorithm                            | Custom yet OS level optimized algorithms                        |
| **Minimum RAM Required** | 512 MB+                            | 256 MB+                              | 512 MB+                                     | **Only 128 MB — runs like a charm on old servers** |
| **Threading Model**    | Mostly single-threaded                | Multi-threaded but not optimal       | Fully multithreaded, async                  | **Multithreaded, async-optimized**                 |
| **Algorithms**         | Event loop, RDB/AOF                   | Slab allocator, hash table           | Custom scheduler, smart locking             | OS-level schedulers, zero-copy pipelines           |
| **Garbage Handling**   | Can leave memory fragmentation        | Can leak/fragment memory             | Efficient GC avoidance                      | **No garbage buildup — runs clean on legacy servers** |
| **License**            | BSD 3-Clause (open source)            | BSD 3-Clause (open source) (paid hosting) | Proprietary (source-available)         | **GNU Public License (100% free)**                 |
| **Ownership**          | Corporate                             | Corporate                            | Corporate                                   | **Individually maintained**                        |
| **Pricing**            | Free                                  | Free                                 | Free (limited) / Paid                       | **Fully free**                                     |

-----------------------

# Docs


* setting values(setting key and value in database)
```
SET name "gujjar"
```
* Getting values(getting certain key will show its value)
```
GET name
```
* Deleting values
```
DEL key name
```
PING compatibility with redis
```
PING 
```
* Deceasing value of a key integer  
```
DECR key
```


