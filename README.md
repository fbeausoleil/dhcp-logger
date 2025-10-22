# DHCP Packet Logger

A simple DHCP server that listens for DHCP packets and prints them to stdout without responding to clients.

## Features

- Listens on UDP port 67 (standard DHCP port)
- Prints detailed information about received DHCP packets
- Does not respond to clients (packet logger only)
- Shows packet headers, options, and metadata

## Usage

### Build and Run

```bash
# Build the server
go build -o dhcp-server main.go

# Run the server (requires root privileges for port 67)
sudo ./dhcp-server
```

### What it displays

For each DHCP packet received, the server prints:
- Timestamp
- Source address
- Message type (DISCOVER, REQUEST, etc.)
- Operation code
- Hardware type and MAC address
- IP addresses (client, server, gateway)
- Transaction ID
- DHCP options

### Example Output

```
DHCP Server started - listening on port 67
Press Ctrl+C to stop
==========================================

=== DHCP Packet Received ===
Timestamp: 2024-01-15 14:30:25
From: 192.168.1.100:68
Message Type: DHCPDISCOVER
Operation: BOOTREQUEST
Hardware Type: Ethernet
Hardware Address Length: 6
Hops: 0
Transaction ID: 0x12345678
Seconds: 0
Flags: 0x8000
Client IP: 0.0.0.0
Your IP: 0.0.0.0
Server IP: 0.0.0.0
Gateway IP: 0.0.0.0
Client MAC: aa:bb:cc:dd:ee:ff
Server Name: 
Boot File: 
Options:
  53: [1]
  55: [1 3 6 15 31 33 43 44 46 47 119 121 249 252]
  60: [80 120 105 110 101 116 45 85 66 73 45 69 70 73]
================================
```

## Requirements

- Go 1.19 or later
- Root privileges (to bind to port 67)
- Network interface with DHCP traffic

## Notes

- This server only logs packets and does not respond to DHCP requests
- Requires root privileges to bind to port 67
- Useful for debugging DHCP traffic and understanding packet structure
