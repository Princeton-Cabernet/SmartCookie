#!/bin/bash

time="1" 	# one second
int="enp3s0f0"   # network interface

while true
    do
	 latency="`curl -s 128.0.0.5:8080 -w %{time_connect}:%{time_starttransfer}:%{time_total} | tail -1`"
   	 echo "\"curl latency: $latency on interface $int\","
    done

