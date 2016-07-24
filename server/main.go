package main

import (
    "fmt"
    "net"
    "os"
    "github.com/ThatGuyFromFinland/utils"
    "github.com/ThatGuyFromFinland/server/core"
    "strings"
    "strconv"
    "time"
)

const (
    CONN_PORT = "50500"
    CONN_TYPE = "tcp"
)

var clients core.Connections

func main() {
    // Listen for incoming connections.
    serverAddr := core.Address{IP: ip.GetIP(), Port: CONN_PORT}
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

func contactClient(address core.Address) core.Handler {
    conn, _ := net.Dial("tcp", address.IP + ":" + address.Port)
    return conn
}

// Very suboptimal solution, but this will do for now
func handleMessageBuffer(contact core.Dialer) {
    for {
        // This should perhaps save some processing resources? Not sure though
        time.Sleep(1 * time.Millisecond)
        clients.MessageQueue.HandleMessage(contact, clients)
    }
}

// Handles incoming requests, duh
func handleRequest(conn core.Handler) {
    // So basicly I am using the length of connection array as user id; lazy
    if (len(clients.Id) != len(clients.Address)) {
        fmt.Println("Error: client registry has been corrupted, aborting")
        os.Exit(1)
    }
    
    // Make a buffer to hold incoming data and read it
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }
    
    // Do something specific for the data and close the connection
    routeRequest(string(buf[:n]), conn)
    conn.Close()
}

// Kind of action controller, seems quite okay, nice to test also
func routeRequest(request string, conn core.Handler) {
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

func handleClientJoinRequest(body string, conn core.Handler) {
    // Lazy, I know, but does the trick
    id := uint64(len(clients.Id))
    // Debugging code, which I am too tired of removing/adding all the time.
    fmt.Printf("Body: %s", body)
    
    addressParts := strings.SplitN(body, ":", 2)
    if (len(addressParts) == 1) {
        conn.Write([]byte ("Port missing!"))
    } else {
        // This works, but it could be better (more safe)
        clients.Id = append(clients.Id, id)
        clients.Address = append(clients.Address, core.Address{IP: addressParts[0], Port: addressParts[1]})
        conn.Write([]byte ("Welcome! Your id is: " + strconv.Itoa(int(id)) + ", you address is: " + clients.Address[id].IP + ":" + clients.Address[id].Port))
    }
}

func handlePeopleRequest(body string, conn core.Handler) {
    requestId, err := strconv.ParseUint(body, 10, 64)
    fmt.Printf("Request id: %q, err %s", requestId, err)
    
    // Seems a bit clumsy, but will do for now
    var resultIds []string
    for _, id := range clients.Id {
        fmt.Printf("Request Id: %s, iter id: %s, body: %s", strconv.FormatUint(requestId, 10), strconv.FormatUint(id, 10), body)
        if (id != requestId) {
            resultIds = append(resultIds, strconv.FormatUint(id, 10))
        }
    }
    conn.Write([]byte (strings.Join(resultIds,",")))
}


func handleMessageRequest(body string, conn core.Handler) {
    bodySplit := strings.SplitN(body, ":", 2)
    fmt.Printf("Recipients: %s", bodySplit[0])
    recipients := strings.Split(bodySplit[0], ",")
    message := bodySplit[1]
    
    // 1024 kilobyte limit
    if len(message) > 1048576 {
        conn.Write([]byte ("Error: Message too long!"))
        return
    }
    // Max 255 recipients
    if len(recipients) > 255 {
        conn.Write([]byte ("Error: Too many recipients!"))
        return
    }
    
    for _, recipient := range recipients {
        recipientId, _ := strconv.ParseUint(recipient, 10, 64)
        if clients.MessageQueue.InsertNewMessage(recipientId, message) == false {
            fmt.Println("Error: MessageQueue full")
            // Might add debug message, which tells recipients, that didn't get
            // delivered before exiting
            return
        }
    }
    conn.Write([]byte ("Sent: \"" + message + "\" to users " + strings.Join(recipients, ",")))
}
