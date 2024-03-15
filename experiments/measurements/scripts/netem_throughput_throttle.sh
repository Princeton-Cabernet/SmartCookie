#cabernet testbed
#DEV='enp134s0f1'

#lucid testbed
DEV='enp3s0f1'
sudo tc qdisc del dev $DEV root #reset
sudo tc qdisc add dev $DEV root handle 1: htb

for I in `seq 100`; do

	PORT=`expr 2000 + ${I}`;

	echo "Setting throughput for TCP srcport ${PORT}"

	sudo -E tc class add dev $DEV parent 1: classid 1:${PORT} htb rate 100gbit
	sudo -E tc filter add dev $DEV parent 1: protocol ip prio 1 u32 flowid 1:${PORT} match ip sport ${PORT} 0xffff
	sudo -E tc qdisc add dev $DEV parent 1:${PORT} handle ${PORT}:1 netem rate 5Mbit

done
