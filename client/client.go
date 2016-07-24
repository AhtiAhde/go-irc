package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "net"
    "github.com/ThatGuyFromFinland/utils"
)

func main() {
    clientIP := ip.GetIP()
    serverIP := os.Args[1]
    var otherUsers string
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Printf("Your IP address is " + clientIP + ", which port would you prefer using?\n")
    fmt.Println("Press enter to use default [50501]:")
    port := "50501"
    var input string
    
    // Read the port
    input, _ = reader.ReadString('\n')
    if (input != "\n") {
        port = strings.Split(input, "\n")[0]
    }
    
    // Join the server
    joinRequest := "JOIN:" + clientIP + ":" + port
    fmt.Printf("Attempting: %s\n", joinRequest)
    status := sendRequest(joinRequest, serverIP)
    var clientId string
    if status[:21] == "Welcome! Your id is: " {
        clientId = strings.Split(status[21:], ",")[0]
    }
    fmt.Println("Connected to server successfully, your id is: " + clientId)
    otherUsers = sendRequest("PEOPLE:" + clientId, serverIP)
    
    go startListeningForMessages(clientIP, port)
    
    for {
        input, _ = reader.ReadString('\n')
        switch {
            case len(input) > 5 && input[:6] == "/WHOIS": 
                otherUsers = sendRequest("PEOPLE:" + clientId, serverIP)
                fmt.Println("Users online: " + otherUsers)
                break;
            case len(input) > 6 && input[:7] == "/WHOAMI": 
                response := sendRequest("WHOAMI:" + clientIP + ":" + port, serverIP)
                fmt.Println(response)
                break;
            case len(input) > 7 && input[:8] == "/PRIVATE":
                // input[9:] 9th char contains the first space
                payload := strings.SplitN(input[9:], " ", 2)
                sendRequest("MESSAGE:" + payload[0] + ":" + payload[1], serverIP)
                fmt.Printf("To %s: %s\n", payload[0], payload[1])
                break;
            case len(input) > 4 && input[:5] == "/QUIT":
                os.Exit(1)
                //break; Maybe not needed?
            default:
                sendRequest("MESSAGE:" + otherUsers + ":" + input, serverIP)
                fmt.Println("All: " + input)
        }
    }
}

func startListeningForMessages(clientIP string, port string) {
    l, err := net.Listen("tcp", clientIP + ":" + port)
    if err != nil {
        fmt.Println("Error failed starting to listen for server:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Ready for messages")
    
    for {
        // Listen for an incoming connection.
        client, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        buf := make([]byte, 1024)
        // Read the incoming connection into the buffer.
        n, err := client.Read(buf)
        if err != nil {
            fmt.Println("Error reading:", err.Error())
        }
        fmt.Println(string(buf[:n]))
        client.Close()
    }
}

func sendRequest(request string, serverIP string) string {
    server, _ := net.Dial("tcp", serverIP + ":50500")
    fmt.Fprintf(server, request)
    
    ret, _ := bufio.NewReader(server).ReadString('\n')
    return ret
}