/*
 *
 * This script was inspired / extended from the original SSH port forwarding solution
 * here https://stackoverflow.com/questions/21417223/simple-ssh-port-forward-in-golang 
 *
 */
package main

import (
  "fmt"
  "io"
  "net"
  "os"
  "os/signal"
  "strconv"
  "syscall"
)

// REMOTE_PORT will return the remote port binding defined in the
// environment, or default to 8081 value
func REMOTE_PORT() int {
  var port int = 8081
  if portFromEnv := os.Getenv("PORT_FORWARDER_LOCAL_LISTENER_PORT"); portFromEnv != "" {
    if _port, err := strconv.Atoi(portFromEnv); err == nil {
      port = _port
    }
  }
  return port
}

var (
  // Loads from the PORT_FORWARDER_LOCAL_LISTENER_PORT env var (optional)
  localAddress string = fmt.Sprintf("localhost:%d", REMOTE_PORT())

  // Loads from the PORT_FORWARDER_REMOTE_ADDRESS env var (required)
  remoteForwarder string = os.Getenv("PORT_FORWARDER_REMOTE_ADDRESS")
)

// handlePortForward will use the passed in local connection to perform 
// a copy of the data from the local to the remote, and performs a 
// TCP/IP dial to the remote location per handling (this could be 
// optimized to avoid this dial on each call in the future)
func handlePortForward(local net.Conn) {
  fmt.Printf("forwarding connection from local %s to remote %s\n", localAddress, remoteForwarder)

  // Establish a connection to the remote host location
  remoteConnectionForwarded, err := net.Dial("tcp", remoteForwarder)
  if err != nil {
    panic(err)
  }

  // copier is a closure defined in the handler func 
  copier := func(w io.Writer, r io.Reader) {
    written, err := io.copy(w, r)
    if err := nil {
      panic(err)
    }
    fmt.Printf("wrote %d bytes to the destination location\n", written)
  }

  go func() { copier(local, remoteConnectionForwarded) }()
  go func() { copier(remoteConnectionForwarded, local) }()
}

// signalHandler handles performing a graceful shutdown 
func signalHandler() chan os.Signal {
  c := make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGTERM, os.Interrupt)

  return c
}

// shutdown handles reading the signal channel to perform memory cleanup on exit
func shutdown(c chan os.Signal) {
  _ = <-c
  fmt.Println("received shutdown signal...")

  for _, idleConnection := range openRemotes {
    idleConnection.Close()
  }
  os.Exit(0)
}

func main() {
  fmt.Println("Start port forwarder service")

  // Handle panics raised from the server
  defer func() {
     if e := recover(); e != nil {
       fmt.Printf("[CRITICAL] encountered a critical error, recovering from panic, error trace: %v", e)
    }
  }()

  // Setup the graceful termination handler 
  go shutdown(signalHandler())

  if remoteForwarder == "" {
    panic(fmt.Errorf("PORT_FORWARDER_REMOTE_ADDRESS must be defined"))
  }

  listener, err := net.Listen("tcp", localAddress)
  if err != nil {
    panic(err)
  }

  defer listener.Close()

  // Handler listening func
  for {
    localConnection, err := listener.Accept()
    if err != nil {
      panic(err)
    }
    
    // Handle the actual forwarding to the remote
    go handlePortForward(localConnection)
  }
}
