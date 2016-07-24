package core

import (
    "fmt"
    "strings"
    "strconv"
)

type Handler interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
}

type Connections struct {
    AddressBook []AddressEntry
    MessageQueue MessageQueue
}

type AddressEntry struct {
    Id uint64
    IP string
    Port string
}

func (this *Connections) InsertAddress(addressString string) bool {
    addressParts := strings.Split(addressString, ":")
    if (len(addressParts) != 2) {
        return false
    }
    newEntry := AddressEntry{Id: uint64(len(this.AddressBook)), IP: addressParts[0], Port: addressParts[1]}
    this.AddressBook = append(this.AddressBook, newEntry)
    fmt.Println(len(this.AddressBook))
    return true
}

////// Router //////

type Router struct {}

// Kind of action controller, seems quite okay, nice to test also
func (this Router) RouteRequest(request string, conn Handler, clients *Connections) {
    requestSplit := strings.SplitN(request, ":", 2)
    action := requestSplit[0]
    body := ""
    if (len(requestSplit) > 1) {
        body = requestSplit[1]
        body = strings.TrimSuffix(body, "\n")
        fmt.Printf("ACTION IS: %s", action)
        switch {
            case action == "JOIN":
                this.handleClientJoinRequest(body, conn, clients)
            case action == "PEOPLE":
                this.handlePeopleRequest(body, conn, clients)
            case action == "MESSAGE":
                this.handleMessageRequest(body, conn, clients)
        }
    }
}

func (this Router) handleClientJoinRequest(body string, conn Handler, clients *Connections) {
    noErrors := clients.InsertAddress(body)
    if noErrors == false {
        conn.Write([]byte ("Invalid parameters, unable to add address!"))
    } else {
        newEntry := clients.AddressBook[len(clients.AddressBook) - 1]
        conn.Write([]byte ("Welcome! Your id is: " + strconv.Itoa(int(newEntry.Id)) + ", you address is: " + newEntry.IP + ":" + newEntry.Port))
    }
}

func (this Router) handlePeopleRequest(body string, conn Handler, clients *Connections) {
    requestId, err := strconv.ParseUint(body, 10, 64)
    fmt.Printf("Request id: %q, err %s", requestId, err)
    
    // Seems a bit clumsy, but will do for now
    var resultIds []string
    for _, address := range clients.AddressBook {
        fmt.Printf("Request Id: %s, iter id: %s, body: %s", strconv.FormatUint(requestId, 10), strconv.FormatUint(address.Id, 10), body)
        if (address.Id != requestId) {
            resultIds = append(resultIds, strconv.FormatUint(address.Id, 10))
        }
    }
    conn.Write([]byte (strings.Join(resultIds,",")))
}

func (this Router) handleMessageRequest(body string, conn Handler, clients *Connections) {
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


////// Message Queue //////

// These are for testability (duck typing?)
type Dialer func (AddressEntry) Handler

type Message struct {
    Recipient uint64
    Payload string
}

type MessageQueue struct {
    MessageQueue [1024]Message //This could be refactored to plain MessageQueue
}

func (this *MessageQueue) InsertNewMessage(recipient uint64, payload string) bool {
    for index, slot := range this.MessageQueue {
        empty := Message{}
        if slot == empty {
            this.MessageQueue[index] = Message{recipient, payload}
            return true
        }
    }
    return false
}

// Passing contact dialer feels stupid; there is a better way?
func (this *MessageQueue) HandleMessage(contact Dialer, clients *Connections) {
    for index, message := range this.MessageQueue {
        emptyMessage := Message{}
        if message != emptyMessage {
            conn :=establishConnection(message.Recipient, contact, clients)
            conn.Write([]byte (message.Payload))
            this.MessageQueue[index] = emptyMessage
        }
    }
}

// A bit lame trick for testability; hopefully refactoring for channels will fix
func establishConnection(id uint64, contact Dialer, clients *Connections) Handler {
    var ret Handler
    for _, address := range clients.AddressBook {
        if address.Id == id {
            ret = contact(address)
        }
    }
    return ret
}