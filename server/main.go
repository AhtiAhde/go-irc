package main

import (
    "fmt"
    "net"
    "os"
    "github.com/ThatGuyFromFinland/utils"
    "github.com/ThatGuyFromFinland/server/core"
    "time"
)

const (
    CONN_PORT = "50500"
    CONN_TYPE = "tcp"
)

var clients core.Connections
var router core.Router

func main() {
    // Listen for incoming connections.
    serverAddr := core.AddressEntry{IP: ip.GetIP(), Port: CONN_PORT}
    l, err := net.Listen(CONN_TYPE, serverAddr.IP + ":" + serverAddr.Port)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Printf("Listening on %s:%s", serverAddr.IP, serverAddr.Port)
    
    // Start goroutine for message buffer;
    // @todo: Refactor to a channel
    go handleMessageBuffer(contactClient)
    
    // Listen for incoming messages and route them for specific handlers
    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        
        go handleRequest(conn)
    }
}

func contactClient(address core.AddressEntry) core.Handler {
    conn, _ := net.Dial("tcp", address.IP + ":" + address.Port)
    return conn
}

// Very suboptimal solution, but this will do for now
func handleMessageBuffer(contact core.Dialer) {
    for {
        // This should perhaps save some processing resources? Not sure though
        time.Sleep(1 * time.Millisecond)
        clients.MessageQueue.HandleMessage(contact, &clients)
    }
}

// Handles incoming requests, duh
func handleRequest(conn core.Handler) {
    // Make a buffer to hold incoming data and read it
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }
    
    // Do something specific for the data and close the connection
    router.RouteRequest(string(buf[:n]), conn, &clients)
    conn.Close()
}
