package ip

import (
    "os"
    "net"
)

func GetIP() string {
    var ret string
    host, _ := os.Hostname()
        addrs, _ := net.LookupIP(host)
        for _, addr := range addrs {
            if ipv4 := addr.To4(); ipv4 != nil {
                ret = ipv4.String()
            }   
        }
    return ret
}