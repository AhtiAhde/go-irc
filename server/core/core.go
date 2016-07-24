package core


type Handler interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
}

// These are for testability (duck typing?)
type Dialer func (Address) Handler

type Connections struct {
    Id []uint64
    Address []Address
    MessageQueue MessageQueue
}

type Address struct {
    IP string
    Port string
}

type Message struct {
    Recipient uint64
    Payload string
}

type MessageQueue struct {
    MessageQueue [1024]Message //This could be refactored to plain MessageQueue
}

func (messageQueue MessageQueue) InsertNewMessage(recipient uint64, payload string) bool {
    for index, slot := range messageQueue.MessageQueue {
        empty := Message{}
        if slot == empty {
            messageQueue.MessageQueue[index] = Message{recipient, payload}
            return true
        }
    }
    return false
}

// Passing contact dialer feels stupid; there is a better way?
func (messageQueue MessageQueue) HandleMessage(contact Dialer, clients Connections) {
    for index, message := range messageQueue.MessageQueue {
        emptyMessage := Message{}
        if message != emptyMessage {
            conn :=establishConnection(message.Recipient, contact, clients)
            conn.Write([]byte (message.Payload))
            messageQueue.MessageQueue[index] = emptyMessage
        }
    }
}

// A bit lame trick for testability; hopefully refactoring for channels will fix
func establishConnection(id uint64, contact Dialer, clients Connections) Handler {
    var ret Handler
    for i, _ := range clients.Id {
        if clients.Id[i] == id {
            ret = contact(clients.Address[i])
        }
    }
    return ret
}