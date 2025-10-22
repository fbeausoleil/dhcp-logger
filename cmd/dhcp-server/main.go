package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
)

func main() {
	// Allow custom ports via environment variables for testing without root
	v4Port := ":67"
	if p := os.Getenv("DHCP_PORT"); p != "" {
		v4Port = ":" + p
		fmt.Printf("Using custom IPv4 port: %s\n", v4Port)
	}
	v6Port := ":547"
	if p := os.Getenv("DHCPV6_PORT"); p != "" {
		v6Port = ":" + p
		fmt.Printf("Using custom IPv6 port: %s\n", v6Port)
	}

	// Start both listeners concurrently
	errCh := make(chan error, 2)
	go func() { errCh <- listenDHCPv4(v4Port) }()
	go func() { errCh <- listenDHCPv6(v6Port) }()

	fmt.Println("DHCPv4/v6 packet logger started")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("==========================================")

	// Block until any listener returns an error
	err := <-errCh
	if err != nil {
		log.Fatalf("Listener error: %v", err)
	}
}

func listenDHCPv4(port string) error {
	conn, err := net.ListenPacket("udp4", port)
	if err != nil {
		return fmt.Errorf("DHCPv4 listen on %s failed: %w", port, err)
	}
	defer conn.Close()

	fmt.Printf("DHCPv4 Server listening on %s\n", port)

	buffer := make([]byte, 1500)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			log.Printf("DHCPv4 read error: %v", err)
			continue
		}
		packet, err := dhcpv4.FromBytes(buffer[:n])
		if err != nil {
			fmt.Printf("Error parsing DHCPv4 packet from %s: %v\n", clientAddr, err)
			continue
		}
		printDHCPv4Packet(packet, clientAddr)
	}
}

func listenDHCPv6(port string) error {
	conn, err := net.ListenPacket("udp6", port)
	if err != nil {
		return fmt.Errorf("DHCPv6 listen on %s failed: %w", port, err)
	}
	defer conn.Close()

	fmt.Printf("DHCPv6 Server listening on %s\n", port)

	buffer := make([]byte, 4096)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			log.Printf("DHCPv6 read error: %v", err)
			continue
		}
		d, err := dhcpv6.FromBytes(buffer[:n])
		if err != nil {
			fmt.Printf("Error parsing DHCPv6 packet from %s: %v\n", clientAddr, err)
			continue
		}
		printDHCPv6Packet(d, clientAddr)
	}
}

func printDHCPv4Packet(packet *dhcpv4.DHCPv4, clientAddr net.Addr) {
	fmt.Printf("\n=== DHCPv4 Packet Received ===\n")
	fmt.Printf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("From: %s\n", clientAddr)
	fmt.Printf("Message Type: %s\n", packet.MessageType().String())
	fmt.Printf("Operation: %s\n", packet.OpCode.String())
	fmt.Printf("Hardware Type: %s\n", packet.HWType.String())
	fmt.Printf("Hardware Address Length: %d\n", len(packet.ClientHWAddr))
	fmt.Printf("Hops: %d\n", packet.HopCount)
	fmt.Printf("Transaction ID: 0x%x\n", packet.TransactionID)
	fmt.Printf("Seconds: %d\n", packet.NumSeconds)
	fmt.Printf("Flags: 0x%x\n", packet.Flags)
	fmt.Printf("Client IP: %s\n", packet.ClientIPAddr)
	fmt.Printf("Your IP: %s\n", packet.YourIPAddr)
	fmt.Printf("Server IP: %s\n", packet.ServerIPAddr)
	fmt.Printf("Gateway IP: %s\n", packet.GatewayIPAddr)
	fmt.Printf("Client MAC: %s\n", packet.ClientHWAddr)
	fmt.Printf("Server Name: %s\n", packet.ServerHostName)
	fmt.Printf("Boot File: %s\n", packet.BootFileName)
	fmt.Printf("Options:\n")
	for code, option := range packet.Options {
		fmt.Printf("  %d: %v\n", code, option)
	}
	fmt.Printf("================================\n")
}

func printDHCPv6Packet(d dhcpv6.DHCPv6, clientAddr net.Addr) {
	fmt.Printf("\n=== DHCPv6 Packet Received ===\n")
	fmt.Printf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("From: %s\n", clientAddr)

	inner := d
	if decap, err := dhcpv6.DecapsulateRelay(d); err == nil && decap != nil {
		fmt.Printf("Encapsulated in Relay-Message\n")
		inner = decap
	}

	if msg, ok := inner.(*dhcpv6.Message); ok && msg != nil {
		fmt.Printf("Message Type: %s\n", msg.Type().String())
		fmt.Printf("Transaction ID: 0x%02x%02x%02x\n", msg.TransactionID[0], msg.TransactionID[1], msg.TransactionID[2])
		fmt.Printf("Options:\n")
		for _, opt := range msg.Options.Options {
			fmt.Printf("  %s\n", opt.String())
		}
	} else {
		fmt.Printf("%s\n", d.Summary())
	}
	fmt.Printf("================================\n")
}
