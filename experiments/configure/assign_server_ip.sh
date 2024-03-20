#!/bin/bash 
if="enp3s0f1"

server=131.0 #jc6 machine 
for I in `seq 0 60`; do	
for J in `seq 0 255`; do
	var_ip=$I.$J
	ip addr add $server.$var_ip dev $if;
	if [ $J == 255 ]
	then
		echo "ip addr add $server.$var_ip dev $if";
		timestamp="$(date +"%M.%S,%N")"
		echo "$timestamp";
	fi
done
done
