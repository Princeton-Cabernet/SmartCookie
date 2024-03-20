# SmartCookie Artifact 

This repository contains the prototype source code and instructions for artifact evaluation for our USENIX Security'24 paper [SMARTCOOKIE: Blocking Large-Scale SYN Floods with a Split-Proxy Defense on Programmable Data Planes](#).

## 1. Contents
The artifact consists of two major pieces: 1) source code for the switch agent and server agent of SMARTCOOKIE's split-proxy SYN-flooding defense, related benchmark code, and measurement scripts (showing availability), and 2) a hardware testbed for running and evaluating SMARTCOOKIE under key attack scenarios (showing functionality and reproducibility). 
* `p4src/` includes the Switch Agent program that calculates SYN cookies using HalfSipHash.
	* `p4src/benchmark/` contains variants of the Switch Agent, for benchmarking max hashing rate using a different hash function (AES).
* `ebpf/` includes the Server Agent programs that process cookie-verified new connection handshake and false positive packets.
	* `ebpf/benchmark/` contains a XDP-based SYN cookie generator, for benchmarking max hashing rate of a server-only solution.
* `experiments/` includes the relevant scripts for running key experiments.
	* `experiments/measurements/` contains scripts for collecting client-side latency and server-side CPU measurements.

## 2. Description & Requirements
For the purposes of this artifact evaluation, our testbed consists of five servers and an Intel Tofino Wedge32X-BF programmable switch.
Three machines act as adversaries, each with a XX-core Intel Xeon Silver 4114 CPU and a Mellanox ConnectX-5 2x100Gbps NIC, generating attack traffic using DPDK 19.12.0 and pktgen-DPDK.
Two other machines act as server and client, each with 8-core Intel Xeon D-1541 CPUs and Intel X552 2x10Gbps NICs. 
**For simplicity of artifact evaluation, we are providing evaluators with access to our preconfigured testbed (access instructions below). Instructions for installations and dependencies are briefly included for completeness, but all installations and dependencies are already in place for the evaluation testbed.** 
Next, we describe how to access the testbed, what hardware and software dependencies are required (these are preconfigured for the testbed), and what additional benchmarks can be run. 

### 2.1 Security, privacy, and ethical concerns
There are no security, privacy, or ethical concerns or risks to evaluators or their machines. All experiments can be run on the authors’ testbed, which is provisioned for the planned attack rates. For testbed access, please do not share or distribute the private key (discussed further below).

### 2.2 Accessing the testbed 
* Save the SSH private access key (shared with you directly on the submission site) to your local machine under `~/.ssh/usenixsec24ae.priv.id_rsa`. Note your `sudo` password was also shared with you on the submission site. 
* Update the permissions with `chmod 600 ~/.ssh/usenixsec24ae.priv.id_rsa`. 
* Start the ssh-agent and load the key: `eval $(ssh-agent -s)` and `ssh-add ~/.ssh/usenixsec24ae.priv.id_rsa`.
* Put the following text into your local machine’s `~/.ssh/config`, such that you can ssh into the machines by hostname using the public-facing proxy port. Your public keys are already in place.
```
    Host jc-gateway # Gateway server (jumphost) 
        HostName 24.229.186.134
        User usenixsec24ae
        IdentityFile ~/.ssh/usenixsec24ae.priv.id_rsa
    Host jc5        # Client server
        HostName 10.0.0.5
        User usenixsec24ae
        IdentityFile ~/.ssh/usenixsec24ae.priv.id_rsa
        ProxyJump jc-gateway
    Host jc6        # Server agent
        HostName 10.0.0.6
        User usenixsec24ae
        IdentityFile ~/.ssh/usenixsec24ae.priv.id_rsa
        ProxyJump jc-gateway
    Host opti1        # Attack server 1 
        HostName 10.0.0.7             
        User usenixsec24ae
        IdentityFile ~/.ssh/usenixsec24ae.priv.id_rsa
        ProxyJump jc-gateway
    Host opti2        # Attack server 2 
        HostName 10.0.0.8             
        User usenixsec24ae
        IdentityFile ~/.ssh/usenixsec24ae.priv.id_rsa
        ProxyJump jc-gateway
    Host jc4        # Attack server 3 
        HostName 10.0.0.4
        User usenixsec24ae
        IdentityFile ~/.ssh/usenixsec24ae.priv.id_rsa
        ProxyJump jc-gateway
    Host jc-tofino  # Switch agent 
        HostName 10.0.0.100
        User jsonch
        IdentityFile ~/.ssh/usenixsec24ae.priv.id_rsa
        ProxyJump jc-gateway
```
### 2.3 Hardware dependencies 
The switch agent requires an Intel Tofino Wedge32X-BF programmable switch. In order to stress test the switch agent and observe the full capacity of the defense, the adversarial machines must be capable of generating at least 150 Mpps of combined adversarial traffic. This can be accomplished with either two attack machines with 20 cores and 2x100Gbps links, or with three or more machines with fewer cores. 

### 2.4 Software dependencies (For evaluation simplicity, all software dependencies are pre-installed and configured on the artifact testbed.)
* Switch Agent Prerequisite: please use `bf-sde` version 9.7.1 or newer to compile the P4 program. 
* Server Agent Prerequisite: please use kernel `5.10` or newer and the latest version of the `bcc` toolkit. (For Ubuntu, you may run `sudo apt-get install bpfcc-tools python3-bpfcc linux-headers-$(uname -r)`)
* Adversary Machine Prerequisite: please use `DPDK` `19.12.0` or newer and a matching `pktgen-DPDK` version.

### 2.5 Benchmarks 
We compare the cookie hashing performance of SMARTCOOKIE’s switch-based HalfSipHash to that of AES (on the switch) and XDP (on the server). Source code and setup instructions are under `/p4src/benchmark` and `/ebpf/benchmark/` respectively.

## 3. Usage and a Basic Test (Estimate: 15 human-minutes)
We next describe the setup and configuration steps to launch SMARTCOOKIE and prepare the testbed environment for evaluation. We also walk through a simple functionality test of the switch agent and server agent, with an end-to-end connection test between a client and server. 

### 3.1 Compiling and launching the Switch Agent (Terminal 1) 
* First, open a new terminal window and SSH into the switch `ssh jc-tofino`.
* Clone the SMARTCOOKIE artifact repo and `cd SmartCookie-Artifact/p4src`.
* Run the `./switchagent_compile.sh` script to compile the program. This may take a few seconds, and you will see some warnings, but these can safely be ignored. Note this step only needs to be done once, unless there are changes made to the program. 
* Once the compilation is complete, run `./switchagent_load.sh` to load the `SMARTCOOKIE-HalfSipHash.p4` program onto the switch. A successful load should output `bfruntime gRPC server started` as the last log line and land on the switch driver shell starting with `bfshell>`.
* Keep this terminal open while running experiments, and open other terminals for other operations.
* If, for any reason you need to restart the switch driver, run `sudo killall bf_switchd` first, then run `./switchagent_load.sh` to reload the program again.
* Next, to configure the switch interfaces, copy and paste the below in the `bfshell>` to manually bring up ports:
  ```
    ucli
    pm
    port-add 1/1 10G NONE
    an-set 1/1 2
    port-enb 1/1
    an-set 1/1 1   
    port-add 1/3 10G NONE
    an-set 1/3 2
    port-enb 1/3
    an-set 1/3 1    
    port-add 3/0 100G RS     
    port-enb 3/0
    port-add 4/0 100G RS 
    port-enb 4/0
    port-add 5/0 40G NONE
    an-set 5/0 2
    port-enb 5/0
    an-set 5/0 1   
    port-add 6/0 40G NONE
    an-set 6/0 2
    port-enb 6/0
    an-set 6/0 1  
    show
    rate-period 1
    rate-show 
  ```
The last three `bfshell>` commands will list packet counts and throughput rates for each of the interfaces linked to the servers. 
This is the mapping between the servers and switch ports.
* CLIENT (10G link): Port 1/1 with DPID 129 (hex 0x81) is linked to jc5, with assigned client IP address `129.0.0.5`
* SERVER (10G link): Port 1/3 with DPID 131 (hex 0x83) is linked to jc6, with assigned server IP address `131.0.0.6` 
* ATTACK SERVER 1 (100G link): Port 3/0 with DPID 144 (hex 0x90) is linked to `opti1`, with assigned IP address `144.0.0.7`. 
* ATTACK SERVER 2 (100G link): Port 4/0 with DPID 152 (hex 0x98) is linked to `opti2`, with assigned IP address `152.0.0.8`. 
* ATTACK SERVER 3 (two 40G links): Port 5/0 with DPID 160 (hex 0xA0) is linked to `jc4` port 1, with assigned IP address `160.0.0.4`, and Port 6/0 with DPID 168 (hex 0xA8) is linked to `jc4` port 0, with assigned IP address `168.0.0.4`. 

### 3.2 Launching the Server Agent (Terminal 2, 3, & 4)
* Open three other terminal windows and access the server agent with `ssh jc6` in each window.
* Clone the artifact repo on `jc6` if you haven't already, and `cd SmartCookie-Artifact/ebpf`.
* Run `./configure/configure_jc6.sh` once to configure static IP addresses and ARP entries.
* Next, use the provided python scripts in the separate terminals to compile and load the eBPF programs to the interface connected to the switch:
	* 1) `sudo python3 xdp_load.py enp3s0f1` for ingress 
	* 2) `sudo python3 tc_load.py enp3s0f1` for egress 
* You should see output that the programs have been loaded.
* Finally, run the following python script to sync timestamps between the server agent and switch agent, which is necessary for cookie checks: `sudo python3 send_ts.py`. 

### 3.3 A Quick Functionality Test (Terminal 5 & 6) 
* To test a simple end-to-end connection between the `jc5` client and `jc6` server (protected by the intermediate switch agent and server agent), open two more terminals.
* SSH into the client with `ssh jc5` and SSH once more into the server with `ssh jc6`.
* On the server `jc6`, start up a `netcat` server with `nc -l -p 2020`.
* On the client `jc5`, connect to the `netcat` server with `nc 131.0.0.6 2020`.
* The client will seamlessly connect to the server after verification at the switch agent, and you can send messages between the client and server, with the messages popping up on the receiving side.
* If you are curious, you can use `tcpdump -evvvnX -i enp3s0f1` on both client and server to view the full packet sequence during connection setup, and map it to that of Figure 4 in the paper.
* Note that tcpdump is positioned after XDP on the _ingress_ pipeline, and after TC on the _egress_ pipeline (XDP-->tcpdump--> network stack on ingress, and network stack-->TC->tcpdump on egress).

# 4. Evaluation Workflow 
There are three main experiments that showcase the key results and major claims of our work. These are described next. 

## 4.0 Major Claims 

* C1: SMARTCOOKIE defends against attacks _without packet loss_ until high rates (up to 136 Mpps), significantly outperforming the benchmarks of other defenses, which become exhausted at attack rates starting at only ~1.3 Mpps up to ~52 Mpps. This is proven by experiment (E1), and described in Section 8.2 of the paper.
* C2: During attacks, SMARTCOOKIE protects benign clients from performance penalties and protects servers from additional CPU usage. It adds little to no latency overhead to benign connections during attacks, and any latency is comparable to the baseline latency with no ongoing attack. Additionally, it protects the server's CPU during attacks, fully keeping the CPU resources for other applications. This is proven by experiments (E2) and (E3), and shown in Section 8.3 and 8.4 of the paper.

## 4.1 Experiment 1 - Hashing Throughput (Estimate: 1 human-hour)
**Description:** Compare the hashing throughput SMARTCOOKIE-HalfSipHash (SC-HSH) can achieve _without packet loss_ to the maximum hashing throughput of the three benchmarks: 
* SMARTCOOKIE-AES (SC-AES)
* XDP-HalfSipHash (XDP-HSH)
* Kernel-SipHash (K-SH)
  
Use DPDK to send spoofed attack packets to the server while increasing sending rates, and observe the response packet rates untl loss is observed on the switch (for SC-AES and SC-AES) or server (for K-SH and XDP-HSH). (As noted in the paper, since our benchmarks perform one hash calculation per SYN packet, we effectively measure maximum hashing throughput.) The Tx (response) rates should exactly match Rx (received) rates for as long as SMARTCOOKIE or the benchmark is handling the attack without any packet loss. Once a defense begins to reach its capacity, the Tx rate will begin to dip below Rx rates. 

Each of the defense benchmarks have slightly different setup and attack steps, which are described next. 
**Note that the workflow for each experiment benchmark MUST BE RUN SEPARATELY, but instructions are grouped together below for some experiments, since many steps overlap.**

### Experiment 1A and 1B: SC-HSH and SC-AES 

**Initial Preparation for Experiment 1A:** 
* For SC-HSH, launch the switch agent in `jc-tofino`, as described in `3.1` above.

**Initial Preparation for Experiment 1B:** 
* For SC-AES, launch the AES variant of the switch agent by following the workflow described in `3.1`, with the exception of these different steps:
	* `cd /p4src/benchmark` (_instead_ of `cd /p4src`).
   	* Run the `./aes_switchagent_compile.sh` script to compile the program (_instead_ of `./switchagent_compile.sh`). Note this step only needs to be done once, unless there are changes made to the program. 
   	* Once the compilation is complete, run `./aes_switchagent_load.sh` to load the `SMARTCOOKIE-AES.p4` program onto the switch (_instead_ of `./switchagent_load.sh`).
   	* Follow the remainder of the steps in `3.1` (e.g., initializing ports). 
* Then, run the controller script to load an arbitrary encryption key (this is required to set up recirculation rounds correctly): `python3 SmartCookie-AES-controller/install_key.py 0x000102030405060708090a0b0c0d0e0f`. The script may take a few seconds to a minute to install the key. 

**Attack Preparation for both Experiment 1A and 1B:**
* In three additional terminals, SSH into the attack machines: `ssh opti1`, `ssh opti2`, and `ssh jc4`. `DPDK` and `pktgen-DPDK` are already configured for you.
* For each attack terminal, `cd /home/shared/pktgen-dpdk` and launch pktgen with `sudo -E tools/run.py testbed`.
* If the server has been rebooted recently, reconfigure the huge pages: `cd /home/shared/dpdk/usertools` and run `./dpdk-setup.sh`. 
  	* On `opti1` and `opti2`, we have non-NUMA systems, so choose the option to setup hugepage mappings for non-NUMA systems [5]. Meanwhile, for `jc4`, choose the option for NUMA systems [52].
  	* Enter 8192 pages per node.
  	* Exit the script and return to the above steps to launch pktgen. 

**Attack Execution for both Experiment 1A and 1B:** 
* From within the `Pktgen:/>` console of each of the attack machines, launch the SYN flood against the `jc6` server, using the following commands (which set the SYN flag 0x02 with a random mask, and spoof source IPs).
* On attack server `opti1`, copy-paste the following comands:
   	```
	set 0 type ipv4
	set 0 count 0
	set 0 burst 10000
	set 0 rate 1
	enable 0 random 
	set 0 rnd 0 46 ........00000010................
	set 0 proto tcp
	set 0 size 40 
	set 0 src mac 00:00:00:00:00:90
	set 0 dst mac 00:00:00:00:00:83 
	set 0 src ip 144.0.0.7/32 
	set 0 dst ip 131.0.0.6
	set 0 dport 8090 
	start 0
	```
* On attack server `opti2`, copy-paste the following comands:
   	```
	set 0 type ipv4
	set 0 count 0
	set 0 burst 10000
	set 0 rate 1
	enable 0 random 
	set 0 rnd 0 46 ........00000010................
	set 0 proto tcp
	set 0 size 40 
	set 0 src mac 00:00:00:00:00:90
	set 0 dst mac 00:00:00:00:00:83 
	set 0 src ip 144.0.0.7/32 
	set 0 dst ip 131.0.0.6
	set 0 dport 8090 
	start 0
	```
* On attack server `jc4`, copy-paste the following comands: 
   	```
	set 0 type ipv4
	set 0 count 0
	set 0 burst 10000
	set 0 rate 1
	enable 0 random 
	set 0 rnd 0 46 ........00000010................
	set 0 proto tcp
	set 0 size 40 
	set 0 src mac 00:00:00:00:00:A8
	set 0 dst mac 00:00:00:00:00:83
	set 0 src ip 168.0.0.4/32 
	set 0 dst ip 131.0.0.6
	set 0 dport 8090 
	
	set 1 type ipv4
	set 1 count 0
	set 1 burst 10000
	set 1 rate 1
	enable 1 random 
	set 1 rnd 0 46 ........00000010................
	set 1 proto tcp
	set 1 size 40 
	set 1 src mac 00:00:00:00:00:A0
	set 1 dst mac 00:00:00:00:00:83
	set 1 src ip 160.0.0.4/32 
	set 1 dst ip 131.0.0.6
	set 1 dport 8090 
	
	start 0
	start 1   
	```
* The commands `start 0` and `start 1` begin the attack, and you should see `pktgen`'s continuous Rx/Tx rates in the `Pktgen:/>` consoles. (Note: If the console displays ever get messy, `page main` will reset the display.)
* In the switch agent's `bf-sde.pm>` console, the command `rate-show` will also show Rx/Tx rates of the attack on the switch (ports `3/0`, `4/0`, `5/0`, and `6/0`).

**Results for both Experiment 1A and 1B:**
* To observe the maximum attack rate that SC-HSH and SC-AES can handle before any packet loss, you can play around with increasing the sending attack rate with commands `set 0 rate X` and `set 1 rate X`, with a maximum `X` of 100. As long as the Rx/Tx rates observed with `rate-show` in the switch agent match each other, the switch agent is successfully defending against the SYN flood attack packets without any packet loss.
* **SC-HSH can accomplish this up until rates of ~136 Mpps, while SC-AES can only achieve ~52 Mpps.**
* To verify these maximum rates directly, use the following attack rates: 
	* Turn off attack traffic from both `opti1` and `opti2` with `stop 0`.
   	* Max out the sending rates on both ports on `jc4` with `set 0 rate 100` and `set 1 rate 100`. Using `rate-show` in the switch console, you should see the combined Tx (response) rate from both ports matches the combined Rx (received) rate at ~37.8 Mpps. 
   	* For SC-AES: Turn on the attack from `opti1` with `set 0 rate 1` and `start 0`. Refresh the switch counters (`rate-show`), and confirm the Rx/Tx rates still match (it should be ~39 Mpps). Try inching the attack rate up on `opti1` with `set 0 rate 5`, and observe that at ~45 Mpps total attack rates the loss on the switch remains 0 (or very close to it). Increase the attack rate to ~52 Mpps with `set 0 rate 10` on `opti1`, and observe some now-consistent loss on the Tx rate of the switch (although it should still be a relatively low loss rate). Finally, increase the attack rate to ~59 Mpps with `set 0 rate 15` on `opti1`, and observe that the Tx rate on the switch has now dropped consistently and significantly below the Rx rate (the Tx should be about ~53 Mpps), showing that SC-AES has reached its maximum defense capacity. To see this effect magnified, increase the attack rate to 80 Mpps with `set 0 rate 100` on `opti1`, and observe that the response rate drops to ~41 Mpps, compared to the ~80 Mpps attack rate. At this point, the defense is effectively dropping 50% of the traffic it receives without being able to process it.  
  	* For SC-HSH: Directly max out the sending rate from `opti1` with `set 0 rate 100` (with `jc4` attack continuing as well). You should see that even under an ~80 Mpps attack, the SMARTCOOKIE-HalfSipHash switch agent defends against the attack _without any packet loss_ (the Tx rate closely matching the Rx rate on the switch). SC-HSH can achieve this performance up to ~135 Mpps, but unfortunately due to the hardware limitations in our artifact testbed (which is separately operated from the testbed hardware we used in our original experiments), we can only demonstrate attack rates up to XX Mpps. ### FIX ME! 


### Experiment 1C and 1D: XDP-HSH and K-SH 

**Initial Preparation for both Experiment 1C and 1D:** 
* For XDP-HSH and K-SH, there is no defense running on the switch. Instead, we need to launch a supporting wire program (`synflood_assist.p4`) on the switch, which ensures attack traffic from DPDK are properly formed SYN packets before they are forwarded onto the server (where the XDP and kernel defenses operate). Follow the workflow described in `3.1`, with the exception of these different steps:
	* `cd /p4src/benchmark` (_instead_ of `cd /p4src`).
   	* Run the `./switch_assist_compile.sh` script to compile the program (_instead_ of `./switchagent_compile.sh`). Note this step only needs to be done once, unless there are changes made to the program. 
   	* Once the compilation is complete, run `./switch_assist_load.sh` to load the `synflood_assist.p4` program onto the switch (_instead_ of `./switchagent_load.sh`).
   	* Follow the remainder of the steps in `3.1` (e.g., initializing ports).
* For just XDP-HSH, we also need to bring up benchmark defense code on the server: 
	* SSH into the server with `ssh jc6` and `cd SmartCookie-Artifact/ebpf/benchmark`.
   	* Run `sudo python3 xdp_cookie_load.py enp3s0f1 3` to launch the XDP cookie defense. (The final argument specifies the IFINDEX which is associated with the interface, and can be found with `ip a`.)
* For both XDP-HSH and K-SH, launch an `http` server on the `jc6` server that will be ready to accept connection requests from incoming SYN packets:
	* SSH into the server with `ssh jc6` and `cd SmartCookie-Artifact/experiments`.
   	* Run `go run http_server.go`, which will launch an HTTP server that listens for connections on port 8090. 
* Finally, for both XDP-HSH and K-SH, load measurement scripts on the `jc6` server that will track its Rx and Tx rates:
	* Open two new terminals. In each terminal, SSH into the server with `ssh jc6` and `cd SmartCookie-Artifact/experiments/measurements`.
   	  	* 1) `./rx_pps.sh` to capture a continous Rx packet count  
	 	* 2) `./tx_pps.sh` to capture a continous Tx packet count
  
**Attack Preparation for both Experiment 1C and 1D:**
* In one additional terminal, SSH into JUST the `opti1` attack machine: `ssh opti1` (_we only need one attack machine for these benchmarks, as they are easily overwhelmed by lower attack rates_). `DPDK` and `pktgen-DPDK` are already configured for you.
* In the `opti1` terminal, `cd /home/shared/pktgen-dpdk` and launch pktgen with `sudo -E tools/run.py testbed`.

**Attack Execution for both Experiment 1C and 1D:** 
* From within the `Pktgen:/>` console, launch the attack against the `jc6` server, using the following commands (**NOTE:** we are generating spoofed UDP packets here, which the `synflood_assist.p4` program on the switch converts to properly formed SYN packets before sending to the server, as generating properly formed SYN packets was less convenient to do directly using `pktgen`).
* On attack server `opti1`, copy-paste the following comands: 
```
	set 0 type ipv4
	set 0 count 0
	set 0 burst 10000
	set 0 rate 0.01
	enable 0 range 
	range 0 proto udp
	range 0 size 64 64 64 0 
	range 0 src mac 00:00:00:00:00:90 00:00:00:00:00:90 00:00:00:00:00:90 00:00:00:00:00:00
	range 0 dst mac 00:00:00:00:00:83 00:00:00:00:00:83 00:00:00:00:00:83 00:00:00:00:00:00
	range 0 src ip 144.0.0.7 144.0.0.7 144.0.255.255 0.0.0.1
	range 0 dst ip 131.0.0.6 131.0.0.6 131.0.0.6 0.0.0.0 
	range 0 dst port 8090 8090 8090 0
	start 0
```
* The commands `start 0` and `start 1` begin the attack, and you should see `pktgen`'s continuous Rx/Tx rates in the `Pktgen:/>` consoles. (Note: If the console displays ever get messy, `page main` will reset the display. Also, `page main` displays the sending counters, but does not reflect the configurations from the commands above. Go to `page range` to see these.)
* In the switch's `bf-sde.pm>` console, the command `rate-show` will also show Rx/Tx rates of the attack on the switch (port `3/0` for `opti1`). **Look at Tx rates on `1/3` to see the packets sent to the `jc6` server, and compare this to the Rx rates on `1/3`, which shows the packets received back from the server (the server's response rate).**

**Results for both Experiment 1C and 1D:**
* To observe the maximum attack rate that XDP-HSH and K-SH can handle before any packet loss, you can play around with increasing the sending attack rate with commands `set 0 rate X` with a maximum `X` of 100. As long as the Rx/Tx rates observed with `rate-show` for port `1/3` in the switch match each other, and the `rx_pps.sh` and `tx_pps.sh` measurement script outputs on the server also generally match each other, the benchmark is defending against the SYN flood attack packets without packet loss.
* **XDP-HSH can only accomplish this until around ~7.3 Mpps, while K-SH can only reach ~1.3 Mpps before the server's CPUs are exhausted.**
  * To verify these maximum rates directly, use the following attack rates: 
   	* Turn on the attack from just `opti1` with `set 0 rate 0.5` and `start 0`. Using `rate-show` in the switch console, you should see the Tx (send) rate on port `1/3` match the Rx (received) rate at ~0.74 Mpps. The measurement scripts on the server should also show a similar count of Rx and Tx packets.  
   	* For XDP-HSH: Increase the attack from `opti1` with `set 0 rate 3`. Refresh the switch counters (`rate-show`), and confirm the Rx/Tx rates still match (it should be ~4.45 Mpps). Try inching the attack rate up with `set 0 rate 4.5`, and observe that at ~6.6 Mpps attack rates we are beginning to see some loss in the server's response (Rx counters on port `1/3` on the switch). Increase the attack rate to ~7.3 Mpps with `set 0 rate 5` on `opti1`, and observe some now-consistent loss on the Rx rate of port `1/3` on the switch. Finally, increase the attack rate to ~14.7 Mpps with `set 0 rate 10` on `opti1`, and observe that the Rx rate on `1/3` of the switch has now dropped consistently and significantly below the Tx rate, showing that XDP-HSH has reached its maximum defense capacity. 
  	* For K-SH: Slowly increase the attack from `opti1` with `set 0 rate 0.75`. Refresh the switch counters (`rate-show`), and confirm the Rx/Tx rates still match (it should be ~1.11 Mpps). Try inching the attack rate up with `set 0 rate 0.9`, and observe that at ~1.3 Mpps attack rates we are beginning to see some loss in the server's response (Rx counters on port `1/3` on the switch). Increase the attack rate to ~2.2 Mpps with `set 0 rate 1.5` on `opti1`, and observe some now-consistent loss on the Rx rate of port `1/3` on the switch. Finally, increase the attack rate to ~14.7 Mpps with `set 0 rate 10` on `opti1`, and observe that the Rx rate on `1/3` of the switch has now dropped consistently and significantly below the Tx rate, showing that K-SH has reached its maximum defense capacity. 



## 4.2 Experiment 2 - Latency Overhead (Estimate: 30 human-minutes)
**Description:** Measure the end-to-end setup latency for benign client connections during an attack, to show that SMARTCOOKIE adds _little to no latency overhead_ to the baseline without any attack. 

## 4.3 Experiment 3 - Server CPU Usage (Estimate: 30 human-minutes)
**Description:** Measure the CPU usage on the server during an attack, to show that SMARTCOOKIE fully protects server CPU usage during attack. 





## Citing
If you find this implementation or our paper useful, please consider citing:

    @inproceedings{yoo2024smartcookie,
        title={SMARTCOOKIE: Blocking Large-Scale SYN Floods with a Split-Proxy Defense on Programmable Data Planes},
        author={Yoo, Sophia and Chen, Xiaoqi and Rexford, Jennifer},
        booktitle={33rd USENIX Security Symposium (USENIX Security 24)},
        year={2024},
        publisher={USENIX Association}
    }

## License

Copyright 2023 Sophia Yoo & Xiaoqi Chen, Princeton University.

The project's source code are released here under the [GNU Affero General Public License v3](https://www.gnu.org/licenses/agpl-3.0.html). In particular,
- You are entitled to redistribute the program or its modified version, however you must also make available the full source code and a copy of the license to the recipient. Any modified version or derivative work must also be licensed under the same licensing terms.
- You also must make available a copy of the modified program's source code available, under the same licensing terms, to all users interacting with the modified program remotely through a computer network.

(TL;DR: you should also open-source your derivative work's source code under AGPLv3.)
