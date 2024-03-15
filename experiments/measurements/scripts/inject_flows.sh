#!/bin/bash

time="1" 	# one second
int="enp134s0f1"   # network interface
dst_addr="128.0.0.6" 
dport="8090"

num_requests=${1?Error: enter the number of curl requests to make/flows to inject.}
src_port=${2?Error: specify the desired starting local src port.}

for (( i=0; i<$num_requests; i++))
    do
	 src_port_i=$(( $src_port + $i )) 
	 latency="`curl  --retry 2 --retry-connrefused --retry-delay 0 --retry-all-errors --local-port $src_port_i -s $dst_addr:$dport`"
   	 echo "\"Made curl request on interface $int with local src port $src_port_i\","
    done
echo "\"Completed $i curl requests.\""
