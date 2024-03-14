#!/bin/bash

# Network interface configuration
INTERFACE="enp3s0f1"
IP_ADDRESS="128.0.0.6"
HW_ADDRESS="00:00:00:00:00:80"
GATEWAY="1.0.0.1"
ARP_ADDRESS="00:00:00:00:00:01"

# Configure IP address and MAC address for the interface
sudo ifconfig $INTERFACE $IP_ADDRESS
sudo ifconfig $INTERFACE hw ether $HW_ADDRESS

# Add static route and ARP entry
sudo route add $GATEWAY dev $INTERFACE
sudo arp -s $GATEWAY $ARP_ADDRESS

# Add route entries
sudo route add -net 128.0.0.0/8 gw $GATEWAY dev $INTERFACE
sudo route add -net 136.0.0.0/8 gw $GATEWAY dev $INTERFACE
sudo route add -net 130.0.0.0/8 gw $GATEWAY dev $INTERFACE

# Turn off TX and RX checksum offloading
sudo ethtool -K $INTERFACE tx off rx off

echo "Network configuration completed successfully."

