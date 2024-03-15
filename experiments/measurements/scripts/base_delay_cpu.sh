#!/bin/bash


dirname=bd_rps_${1?Error: enter directory name to save files in}
iter=20

mkdir ../data/base_delay_overhead/$dirname
for (( i=0; i<$iter; i++ ))
do
	collect_cpu="`sudo perf stat -a --no-aggr sleep 1 2>&1  | grep cycles | grep -v insn | tr -s ' ' |cut -d' ' -f2 | tee base_delay_overhead/$dirname/output$i.txt`"

   	 echo "\"Wrote output to base_delay_overhead/$dirname/output$i.txt\","
 done

