#!/bin/bash 

if="enp3s0f1"
client=129.0 #jc5 machine 
for J in `seq 26 40`; do	
for K in `seq 0 254`; do
	var_ip=$J.$K
	ip addr add $client.$var_ip dev $if;
	if [ $K == 254 ] 
	then
		echo "ip addr add $client.$var_ip dev $if";
		timestamp="$(date +"%M.%S.%N")"
		echo "$timestamp";
	fi
done
done
