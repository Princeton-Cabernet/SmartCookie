#!/bin/bash

time="0.5" 	# time in seconds
#cabernet 
#int="enp134s0f1np1"   # network interface
#lucid 
int="enp3s0f1" 
dirpath=../data/cpu/lucid/jaqen2_synflood
subdirname=rps_${1?Error: enter attack rps to specify directory to save files in.}
iter=5

mkdir $dirpath

# first delete any old subdirectory there might have been
rm -rf $dirpath/$subdirname 
mkdir $dirpath/$subdirname

for (( i=0; i<$iter; i++ ))
do
	collect_cpu="`sudo perf stat -a --no-aggr sleep 1 2>&1  | grep cycles | grep -v insn | tr -s ' ' |cut -d' ' -f2 | tee $dirpath/$subdirname/output$i.txt`"

   	 echo "\"Wrote output to $dirpath/$subdirname/output$i.txt\","
 done

