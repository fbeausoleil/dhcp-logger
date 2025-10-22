#!/bin/bash

echo "Starting DHCPv4/v6 server test..."

# Build binaries
go build -o ./bin/dhcp-server ./cmd/dhcp-server
go build -o ./bin/test-client ./cmd/test-client

# Start the server in background on high ports
DHCP_PORT=6767 DHCPV6_PORT=6667 ./bin/dhcp-server &
SERVER_PID=$!

# Wait a moment for server to start
sleep 1

# Send a test packet
echo "Sending test DHCP packet..."
./bin/test-client

# Wait a moment for packet to be processed
sleep 2

# Stop the server
kill $SERVER_PID 2>/dev/null

echo "Test completed."
