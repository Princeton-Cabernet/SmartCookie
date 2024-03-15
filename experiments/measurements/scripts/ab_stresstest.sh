#!/bin/bash

time="0.5" 	# time in seconds
int="enp134s0f1np1"   # network interface

num_requests=${1?Error: enter the number of requests to send}
concurrency=1
ip="12.0.0.3"
port="8080"


	 request_rate="`ab -n $num_requests -c $concurrency http://$ip:$port/ | grep "Requests per second"`"

   	 echo "\"$request_rate    on interface $int\","


