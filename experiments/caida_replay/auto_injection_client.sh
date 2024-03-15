#!/bin/bash
set -e 
set -x 

src_port=${1?Error: specify the desired starting local src port.}


iter=8 


for (( i=0; i<$iter; i++))
    do
	 src_port_i=$(( $src_port + $i )) 
	 go run injection_client.go $src_port_i
    done
echo "\"Successfully injected 2550*$i flows.\""
