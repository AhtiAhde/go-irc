package main

import (
    "fmt"
    "net"
    "os"
    "github.com/ThatGuyFromFinland/utils"
)

const (
    CONN_HOST = "123.123.123.123"
    CONN_PORT = "50500"
    CONN_TYPE = "tcp"
)

type Connections struct {
    Id []uint64
    Address []Address
}

type Address struct {
    IP [4]uint8
    port uint16
}

func main() {
    // Listen for incoming connections.
    var ipAddress string = ip.GetIP()
    l, err := net.Listen(CONN_TYPE, ipAddress+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Printf("Listening on %s:" + CONN_PORT, ipAddress)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
  // Make a buffer to hold incoming data.
  buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  n, err := conn.Read(buf)
  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }
  // Send a response back to person contacting us.
  // message = "Message received: " + message
  fmt.Printf("Received: %s", string(buf[:n]))
  conn.Write([]byte(string(buf[:n])))
  // Close the connection when you're done with it.
  conn.Close()
}