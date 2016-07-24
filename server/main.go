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
    
    // Start goroutine for message delivery
    router.Init()
    go handleMessages(contactClient, &router)
    
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

// Handles incoming requests, duh
func handleRequest(conn core.Handler) {
    // Make a buffer to hold incoming data and read it
    buf := make([]byte, 1024) //Wondering if 1048577 would cause bad performance?
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }
    
    // Do something specific for the data and close the connection
    router.RouteRequest(string(buf[:n]), conn)
    conn.Close()
}

// Passing contact dialer feels stupid; there is a better way?
func handleMessages(contact core.Dialer, router *core.Router) {
    for {
        select {
            case message := <-router.MessageQueue:
                fmt.Println("Mesasge received, delivering...")
                conn := router.EstablishConnection(message.Recipient, contact)
                conn.Write([]byte (message.Payload))
            default:
                time.Sleep(1 * time.Millisecond)
        }
    }
    
}
