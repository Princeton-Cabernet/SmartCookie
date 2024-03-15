#!/bin/bash

time="1"        # one second
#cabernet 
#int="enp134s0f1np1"   # network interface
#lucid 
int="enp3s0f1"
while true
    do
         rxpkts_old="`ethtool -S $int | grep rx_packets: | cut -d \: -f 2`"
                 sleep $time
         rxpkts_new="`ethtool -S $int | grep rx_packets: | cut -d \: -f 2`"
         rxpkts="`expr $rxpkts_new - $rxpkts_old`"              # evaluate expressions for recv packets
         rxpkts="`expr $rxpkts / $time`"
                 echo "\"rx $rxpkts pkts/ on interface $int\","
    done
