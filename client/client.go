package main

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
    status := this.sendRequest(joinRequest)
    if status[:21] == "Welcome! Your id is: " {
        this.id = strings.Split(status[21:], ",")[0]
    }
    this.userList = this.sendRequest("PEOPLE:" + this.id)
    
    fmt.Println("Connected to server successfully, your id is: " + this.id + ", other users are: " + this.userList)
}

func (this *Client) WaitForCommands() {
    for {
        input, _ := this.reader.ReadString('\n')
        this.handleCommands(input)
    }
}

func (this *Client) sendRequest(request string) string {
    server, _ := net.Dial("tcp", this.serverIp + ":50500")
    fmt.Fprintf(server, request)
    
    ret, _ := bufio.NewReader(server).ReadString('\n')
    return ret
}

func main() {
    client := Client{ip.GetIP(), "50501", "", os.Args[1], "", bufio.NewReader(os.Stdin)}
    client.Init()

    go client.StartListeningForMessages()
    go client.WaitForCommands()
    for {}
}

func (this *Client) handleCommands(input string) {
    switch {
        case len(input) > 5 && input[:6] == "/WHOIS": 
            this.userList = this.sendRequest("PEOPLE:" + this.id)
            fmt.Println("Users online: " + this.userList)
            break;
        case len(input) > 6 && input[:7] == "/WHOAMI": 
            response := this.sendRequest("WHOAMI:" + this.ip + ":" + this.port)
            fmt.Println(response)
            break;
        case len(input) > 7 && input[:8] == "/PRIVATE":
            // input[9:] 9th char contains the first space
            payload := strings.SplitN(input[9:], " ", 2)
            this.sendRequest("MESSAGE:" + payload[0] + ":" + payload[1])
            fmt.Printf("To %s: %s\n", payload[0], payload[1])
            break;
        case len(input) > 4 && input[:5] == "/QUIT":
            os.Exit(1)
            //break; Maybe not needed?
        default:
            this.sendRequest("MESSAGE:" + this.userList + ":" + input)
            fmt.Println("All: " + input)
    }
}
