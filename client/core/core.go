package core

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "net"
    "github.com/ThatGuyFromFinland/utils"
)

type Client struct {
    ip string
    port string
    id string
    serverIp string
    userList string
    reader *bufio.Reader
    sendRequest Dialer
}

type Dialer func (string, string) string

func NewClient (clientId string, serverIp string, handler Dialer) Client {
    return Client{ip.GetIP(), "50501", clientId, serverIp, "", bufio.NewReader(os.Stdin), handler}
}

func (this *Client) StartListeningForMessages() {
    l, err := net.Listen("tcp", this.ip + ":" + this.port)
    if err != nil {
        fmt.Println("Error failed starting to listen for server:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    
    fmt.Println("Ready for messages")
    
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        
        buf := make([]byte, 1024)
        
        // Read the incoming connection into the buffer.
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading:", err.Error())
        }
        fmt.Println(string(buf[:n]))
        conn.Close()
    }
}

func (this *Client) Init() {
    fmt.Printf("Your IP address is " + this.ip + ", which port would you prefer using?\n")
    fmt.Println("Press enter to use default [50501]:")
    
    // Read the port
    input, _ := this.reader.ReadString('\n')
    if (input != "\n") {
        this.port = strings.Split(input, "\n")[0]
    }
    
    // Join the server
    joinRequest := "JOIN:" + this.ip + ":" + this.port
    fmt.Printf("Attempting: %s\n", joinRequest)
    status := this.sendRequest(joinRequest, this.serverIp)
    if status[:21] == "Welcome! Your id is: " {
        this.id = strings.Split(status[21:], ",")[0]
    }
    this.userList = this.sendRequest("PEOPLE:" + this.id, this.serverIp)
    
    fmt.Println("Connected to server successfully, your id is: " + this.id + ", other users are: " + this.userList)
}

func (this *Client) WaitForCommands() {
    for {
        input, _ := this.reader.ReadString('\n')
        this.HandleCommands(input)
    }
}

func (this *Client) HandleCommands(input string) {
    switch {
        case len(input) > 5 && input[:6] == "/WHOIS": 
            this.userList = this.sendRequest("PEOPLE:" + this.id, this.serverIp)
            fmt.Println("Users online: " + this.userList)
            break;
        case len(input) > 6 && input[:7] == "/WHOAMI": 
            response := this.sendRequest("WHOAMI:" + this.ip + ":" + this.port, this.serverIp)
            fmt.Println(response)
            break;
        case len(input) > 7 && input[:8] == "/PRIVATE":
            // input[9:] 9th char contains the first space
            payload := strings.SplitN(input[9:], " ", 2)
            this.sendRequest("MESSAGE:" + payload[0] + ":" + payload[1], this.serverIp)
            fmt.Printf("To %s: %s\n", payload[0], payload[1])
            break;
        case len(input) > 4 && input[:5] == "/QUIT":
            os.Exit(1)
            //break; Maybe not needed?
        default:
            this.sendRequest("MESSAGE:" + this.userList + ":" + input, this.serverIp)
            fmt.Println("All: " + input)
    }
}
