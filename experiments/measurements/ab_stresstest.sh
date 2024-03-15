#!/bin/bash

time="0.5" 	# time in seconds
int="enp5s0f0"   # network interface

num_requests=${1?Error: enter the number of requests to send}
concurrency=1
ip="128.0.0.5"
port="8080"


	 request_rate="`ab -n $num_requests -c $concurrency http://$ip:$port/ | grep "Requests per second"`"

   	 echo "\"$request_rate    on interface $int\","


