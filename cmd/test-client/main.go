package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
)

func main() {
	// Defaults for high ports for local testing
	v4Target := "127.0.0.1:6767"
	if p := os.Getenv("DHCP_PORT"); p != "" {
		v4Target = "127.0.0.1:" + p
	}
	v6Target := "[::1]:6667" // use high port for tests
	if p := os.Getenv("DHCPV6_PORT"); p != "" {
		v6Target = "[::1]:" + p
	}

	// Send DHCPv4 DISCOVER
	hwAddr, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	v4, err := dhcpv4.NewDiscovery(hwAddr)
	if err == nil {
		if conn, err := net.Dial("udp4", v4Target); err == nil {
			_, _ = conn.Write(v4.ToBytes())
			_ = conn.Close()
			fmt.Println("Sent DHCPv4 DISCOVER to", v4Target)
		}
	}

	// Send DHCPv6 SOLICIT
	v6, err := dhcpv6.NewSolicit(hwAddr)
	if err == nil {
		if conn, err := net.Dial("udp6", v6Target); err == nil {
			_, _ = conn.Write(v6.ToBytes())
			_ = conn.Close()
			fmt.Println("Sent DHCPv6 SOLICIT to", v6Target)
		}
	}

	time.Sleep(100 * time.Millisecond)
}
