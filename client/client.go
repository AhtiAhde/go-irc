package main

import (
    "fmt"
    "os"
    "net"
    "github.com/ThatGuyFromFinland/utils"
)

func main() {
    conn, _ := net.Dial("tcp", "172.17.20.181:50500")
    fmt.Fprintf(conn, "JOIN:" + ip.GetIP() + ":" + os.Args[1])
}