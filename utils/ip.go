package ip

import "net"

func Get_IP() {
    ifaces, err := net.Interfaces()
    fmt.Printf("Errors: %d", err)
    for _, i := range ifaces {
        if i.Name == "eth0" {
            fmt.Printf("Yay!")
            addrs, err := i.Addrs()
            fmt.Printf("Errors: %d", err)
            for _, addr := range addrs {
                
                switch v := addr.(type) {
                // IP address
                case *net.IPNet:
                        fmt.Printf("\nAdresses: addr(%d) i(%s) ip(%s) \n", addr, i.Name, v.IP.String())
                        //fmt.Printf("IP Address: %s", v.IP)
                /* This would be for IPv6
                case *net.IPAddr:
                        return v.IP
                */
                }
            }
        }
    }
}