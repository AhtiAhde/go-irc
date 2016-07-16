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

type Handler interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
}

type Connections struct {
    Id []uint64
    Address []Address
    MessageQueue [1024]Message
}

type Address struct {
    IP string
    port string
}

type Message struct {
    Recipient uint64
    Payload string
}

var clients Connections

func main() {
    // Listen for incoming connections.
    serverAddr := Address{IP: ip.GetIP(), port: CONN_PORT}
    l, err := net.Listen(CONN_TYPE, serverAddr.IP + ":" + serverAddr.port)
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
func handleRequest(conn Handler) {
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
    routeRequest(string(buf[:n]), conn)
    // Close the connection when you're done with it.
    conn.Close()
}

func routeRequest(request string, conn Handler) {
    requestSplit := strings.SplitN(request, ":", 2)
    action := requestSplit[0]
    body := ""
    if (len(requestSplit) > 1) {
        body = requestSplit[1]
        body = strings.TrimSuffix(body, "\n")
        fmt.Printf("ACTION IS: %s", action)
        switch {
            case action == "JOIN":
                handleClientJoinRequest(body, conn)
            case action == "PEOPLE":
                handlePeopleRequest(body, conn)
            case action == "MESSAGE":
                handleMessageRequest(body, conn)
        }
    }
}

func handleClientJoinRequest(body string, conn Handler) {
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

func handlePeopleRequest(body string, conn Handler) {
    // var request_id uint64
    request_id, err := strconv.ParseUint(body, 10, 64)
    fmt.Printf("Request id: %q, err %s", request_id, err)
    var ret_ids []string
    for _, id := range clients.Id {
        fmt.Printf("Request Id: %s, iter id: %s, body: %s", strconv.FormatUint(request_id, 10), strconv.FormatUint(id, 10), body)
        if (id != request_id) {
            ret_ids = append(ret_ids, strconv.FormatUint(id, 10))
        }
    }
    conn.Write([]byte (strings.Join(ret_ids,",")))
}

func handleMessageRequest(body string, conn Handler) {
    bodySplit := strings.SplitN(body, ":", 2)
    recipients := bodySplit[0]
    message := bodySplit[1]
    success := true
    for _, recipient := range strings.Split(recipients, ",") {
        recipientId, _ := strconv.ParseUint(recipient, 10, 64)
        if insertNewMessage(recipientId, message) == false {
            fmt.Println("Error: MessageQueue full")
            success = false
        }
    }
    if success == true {
        conn.Write([]byte ("Sent: \"" + message + "\" to users " + recipients))
    }
}

func insertNewMessage(recipient uint64, payload string) bool {
    for index, slot := range clients.MessageQueue {
        empty := Message{}
        if slot == empty {
            clients.MessageQueue[index] = Message{recipient, payload}
            return true
        }
    }
    return false
}