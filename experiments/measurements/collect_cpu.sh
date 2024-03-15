#!/bin/bash

time="0.5" 	# time in seconds
int="enp5s0f0"   # network interface
dirname=b_rps_${1?Error: enter directory name to save files in}
iter=20

mkdir base_overhead/$dirname
for (( i=0; i<$iter; i++ ))
do
	collect_cpu="`sudo perf stat -a --no-aggr sleep 1 2>&1  | grep cycles | grep -v insn | tr -s ' ' |cut -d' ' -f2 | tee base_overhead/$dirname/output$i.txt`"

   	 echo "\"Wrote output to base_overhead/$dirname/output$i.txt\","
 done

