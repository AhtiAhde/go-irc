package main

import (
    "fmt"
    "net"
    "os"
    "github.com/ThatGuyFromFinland/utils"
    "strings"
    "strconv"
)

const (
    CONN_PORT = "50500"
    CONN_TYPE = "tcp"
)

type Connections struct {
    Id []uint64
    Address []Address
}

type Address struct {
    IP string
    port string
}

var clients Connections

func main() {
    // Listen for incoming connections.
    serverAddr := Address{IP: ip.GetIP(), port: CONN_PORT}
    l, err := net.Listen(CONN_TYPE, serverAddr.IP+":"+serverAddr.port)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Printf("Listening on %s:%s", serverAddr.IP, serverAddr.port)
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
    if (len(clients.Id) != len(clients.Address)) {
        fmt.Println("Error: client registry has been corrupted, aborting")
        os.Exit(1)
    }
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)
    // Read the incoming connection into the buffer.
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }
    // Send a response back to person contacting us.
    // message = "Message received: " + message
    request := strings.SplitN(string(buf[:n]), ":", 2)
    action := request[0]
    body := ""
    if (len(request) > 1) {
        body = request[1]
        switch {
            case action == "JOIN":
                handleClientJoinRequest(body, conn)
            case action == "PEOPLE":
                handlePeopleRequest(body, conn)
        }
    }
    // Close the connection when you're done with it.
    conn.Close()
}

func handleClientJoinRequest(body string, conn net.Conn) {
    id := uint64(len(clients.Id))
    fmt.Printf("Body: %s", body)
    addressParts := strings.SplitN(body, ":", 2)
    if (len(addressParts) == 1) {
        conn.Write([]byte ("Port missing!"))
    } else {
        clients.Id = append(clients.Id, id)
        clients.Address = append(clients.Address, Address{IP: addressParts[0], port: addressParts[1]})
        conn.Write([]byte ("Welcome! Your id is: " + strconv.Itoa(int(id)) + ", you address is: " + clients.Address[id].IP + ":" + clients.Address[id].port))
    }
}

func handlePeopleRequest(body string, conn net.Conn) {
    request_id, _ := strconv.ParseUint(body, 10, 64)
    var ret_ids []string
    for _, id := range clients.Id {
        if (id != request_id) {
            ret_ids = append(ret_ids, strconv.FormatUint(id, 10))
        }
    }
    conn.Write([]byte (strings.Join(ret_ids,",")))
}