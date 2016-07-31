package main

import (
    "bufio"
    "net"
    "fmt"
    "os"
    "github.com/ThatGuyFromFinland/client/core"
)

func sendServerRequest(request string, serverIp string) string {
    server, _ := net.Dial("tcp", serverIp + ":50500")
    fmt.Fprintf(server, request)
    
    ret, _ := bufio.NewReader(server).ReadString('\n')
    return ret
}

func main() {
    client := core.NewClient("", os.Args[1], sendServerRequest)
    client.Init()

    go client.StartListeningForMessages()
    go client.WaitForCommands()
    for {}
}
