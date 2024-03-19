#!/bin/bash

time="1" 	# one second
int="enp134s0f1np1"   # network interface

while true
    do
   	 txpkts_old="`ethtool -S enp134s0f1np1 | grep tx_packets | cut -d \: -f 2`"
   		 sleep $time
   	 txpkts_new="`ethtool -S enp134s0f1np1 | grep tx_packets | cut -d \: -f 2`"
   	 txpkts="`expr $txpkts_new - $txpkts_old`"   	  	# evaluate expressions for recv packets
   	 txpkts="`expr $txpkts / $time`"
   		 echo "\"tx $txpkts pkts/ on interface $int\","
    done

