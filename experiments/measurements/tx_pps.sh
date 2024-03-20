#!/bin/bash

time="1" 	# one second
int="enp3s0f1"   # network interface

while true
    do
   	 txpkts_old="`ethtool -S enp3s0f1 | grep tx_packets | cut -d \: -f 2`"
   		 sleep $time
   	 txpkts_new="`ethtool -S enp3s0f1 | grep tx_packets | cut -d \: -f 2`"
   	 txpkts="`expr $txpkts_new - $txpkts_old`"   	  	# evaluate expressions for recv packets
   	 txpkts="`expr $txpkts / $time`"
   		 echo "\"tx $txpkts pkts/ on interface $int\","
    done

