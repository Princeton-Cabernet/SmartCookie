#!/bin/bash
# Calls a go client script 30 times to send 30 requests with a 30 different src ip and fixed src port. Src port netem should be throttled to 1 Mbit prior.  
# Look at the go script for dst ip and dst port. 

time="1" 	# one second
int="enp134s0f1"   # network interface

src_ip_top="130.0.30"
src_ip_bottom="30"
#src_ip=${3?Error: specify the desired src ip.}
src_port=${1?Error: specify the desired starting local src port.}

filename=rps_${2?Error: enter attack rps to specify file name to save data in.}
dirpath=../notebook_eval/data/latency/lucid/jaqen_synflood_nothrottle
mkdir $dirpath


iter=3 

# first delete any old file there might have been 
rm $dirpath/$filename.txt

for (( i=0; i<$iter; i++))
    do
	 src_ip_i=$src_ip_top.$(( $src_ip_bottom + $i ))
	 latency="`go run one_http_req.go $src_ip_i $src_port | tail -1 | cut -c 8- | tee -a $dirpath/$filename.txt`"
   	 echo "\"go client latency: $latency on interface $int with src ip $src_ip_i and src port $src_port\","

    done
echo "\"Wrote output to $dirpath/$filename.txt.\""

echo -ne '\007' #bell
