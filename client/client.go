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
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Printf("Your IP address is " + clientIP + ", which port would you prefer using?\n")
    fmt.Println("Press enter to use default [50501]:")
    port := "50501"
    var input string
    input, _ = reader.ReadString('\n')
    if (input != "\n") {
        port = strings.Split(input, "\n")[0]
    }
    joinRequest := "JOIN:" + clientIP + ":" + port
    fmt.Printf("Attempting: %s\n", joinRequest)
    status := sendRequest(joinRequest, serverIP)
    var clientId string
    if status[:21] == "Welcome! Your id is: " {
        clientId = strings.Split(status[21:], ",")[0]
    }
    fmt.Println("Connected to server successfully, your id is: " + clientId)
    
    go startListeningForMessages(clientIP, port)
    
    for {
        input, _ = reader.ReadString('\n')
        if input[:6] == "/WHOIS" {
            people := sendRequest("PEOPLE:" + clientId, serverIP)
            fmt.Println("Who is there? There are users: " + people)
        } else if input[:5] == "/QUIT" {
            os.Exit(1)
        } else {
            fmt.Println("You: " + input)
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