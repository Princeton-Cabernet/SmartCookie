#!/bin/bash

time="1" 	# one second
int="enp3s0f1"   # network interface
dst_addr="128.0.0.6" 
dport="8090"

filename=rps_${1?Error: enter attack rps to specify file name to save data in.}
dirpath=../data/latency/lucid/jaqen2_synflood
mkdir $dirpath

src_port=${2?Error: specify the desired starting local src port.}
payload_size=14500
iter=31 

# first delete any old file there might have been 
rm $dirpath/$filename.pcap

tcpdump -evvvnX -i $int -w $dirpath/$filename.pcap & 

# for some reason, the first curl request doesn't trigger a RST
# so collect 31 measurements and in Jupyter Notebook script disregard first measurement
for (( i=0; i<$iter; i++))
    do
	 src_port_i=$(( $src_port + $i )) 
	 latency="`curl  --retry 2 --retry-connrefused --retry-delay 0 --retry-all-errors --local-port $src_port_i -s $dst_addr:$dport/$payload_size -w %{time_connect}:%{time_starttransfer}:%{time_total} | tail -1 | rev | cut -c -26 | rev `"
   	 echo "\"curl latency: $latency on interface $int with local src port $src_port_i\","

    done
echo "\"Wrote output to $dirpath/$filename.pcap.\""
