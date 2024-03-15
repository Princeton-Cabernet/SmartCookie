#!/bin/bash

time="1" 	# one second
int="enp3s0f1"   # network interface

while true
    do
   	 rxpkts_old="`ethtool -S enp3s0f1 | grep rx_packets | cut -d \: -f 2`"
   		 sleep $time
   	 rxpkts_new="`ethtool -S enp3s0f1 | grep rx_packets | cut -d \: -f 2`"
   	 rxpkts="`expr $rxpkts_new - $rxpkts_old`"   	  	# evaluate expressions for recv packets
   	 rxpkts="`expr $rxpkts / $time`"
   		 echo "\"rx $rxpkts pkts/ on interface $int\","
    done

