# SmartCookie Artifact 

This repository contains the prototype source code and instructions for artifact evaluation for our USENIX Security'24 paper [SMARTCOOKIE: Blocking Large-Scale SYN Floods with a Split-Proxy Defense on Programmable Data Planes](#).

## Contents
The artifact consists of two major pieces: 1) source code for the switch agent and server agent of SMARTCOOKIE's split-proxy SYN-flooding defense, related benchmark code, and measurement scripts (showing availability), and 2) a hardware testbed for running and evaluating SMARTCOOKIE under key attack scenarios (showing functionality and reproducibility). 
* `p4src/` includes the Switch Agent program that calculates SYN cookies using HalfSipHash.
	* `p4src/benchmark/` contains variants of the Switch Agent, for benchmarking max hashing rate using a different hash function (AES).
* `ebpf/` includes the Server Agent programs that process cookie-verified new connection handshake and false positive packets.
	* `ebpf/benchmark/` contains a XDP-based SYN cookie generator, for benchmarking max hashing rate of a server-only solution.
* `experiments/` includes the relevant scripts for running key experiments.
	* `experiments/measurements/` contains scripts for collecting client-side latency and server-side CPU measurements.

## Description & Requirements
For the purposes of this artifact evaluation, our testbed consists of five servers and an Intel Tofino Wedge32X-BF programmable switch.
Three machines act as adversaries, each with a XX-core Intel Xeon Silver 4114 CPU and a Mellanox ConnectX-5 2x100Gbps NIC, generating attack traffic using DPDK 19.12.0 and pktgen-DPDK.
Two other machines act as server and client, each with 8-core Intel Xeon D-1541 CPUs and Intel X552 2x10Gbps NICs. 
**For simplicity of artifact evaluation, we are providing evaluators with access to our preconfigured testbed (access instructions below). Instructions for installations and dependencies are briefly included for completeness, but all installations and dependencies are already in place for the evaluation testbed.** 
Next, we describe how to access the testbed, what hardware and software dependencies are required (these are preconfigured for the testbed), and what additional benchmarks can be run. 

### Security, privacy, and ethical concerns
There are no security, privacy, or ethical concerns or risks to evaluators or their machines. All experiments can be run on the authors’ testbed, which is provisioned for the planned attack rates. For testbed access, please do not share or distribute the private key (discussed further below).

### Accessing the testbed 
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
        HostName 10.0.0.7             
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
### Hardware dependencies 
The switch agent requires an Intel Tofino Wedge32X-BF programmable switch. In order to stress test the switch agent and observe the full capacity of the defense, the adversarial machines must be capable of generating at least 150 Mpps of combined adversarial traffic. This can be accomplished with either two attack machines with 20 cores and 2x100Gbps links, or with three or more machines with fewer cores. 

### Software dependencies (For evaluation simplicity, all software dependencies are pre-installed and configured on the artifact testbed.)
* Switch Agent Prerequisite: please use `bf-sde` version 9.7.1 or newer to compile the P4 program. 
* Server Agent Prerequisite: please use kernel `5.10` or newer and the latest version of the `bcc` toolkit. (For Ubuntu, you may run `sudo apt-get install bpfcc-tools python3-bpfcc linux-headers-$(uname -r)`)
* Adversary Machine Prerequisite: please use `DPDK` `19.12.0` or newer and a matching `pktgen-DPDK` version.

### Benchmarks 
We compare the cookie hashing performance of SMARTCOOKIE’s switch-based HalfSipHash to that of AES (on the switch) and XDP (on the server). Source code and setup instructions are under `/p4src/benchmark` and `/ebpf/benchmark/` respectively.

## Usage and a Basic Test (Estimate: 15 human-minutes)
We next describe the setup and configuration steps to launch SMARTCOOKIE and prepare the testbed environment for evaluation. We also walk through a simple functionality test of the switch agent and server agent, with an end-to-end connection test between a client and server. 

### Compiling and launching the Switch Agent (Terminal 1) 
* First, open a new terminal window and SSH into the switch `ssh jc-tofino`.
* Clone the SMARTCOOKIE artifact repo and `cd SmartCookie-Artifact/p4src`.
* Run the `./switchagent_compile.sh` script to compile the program. This may take a few seconds, and you will see some warnings, but these can safely be ignored.
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
* ATTACK SERVER 1 (100G link): Port 3/0 with DPID 144 (hex 0x90) is linked to `opti1`. 
* ATTACK SERVER 2 (): 
* ATTACK SERVER 3 (two 40G links): Port 5/0 with DPID 160 (hex 0xA0) is linked to `jc4` port 1 and Port 6/0 with DPID 168 (hex 0xA8) is linked to `jc4` port 0. 

### Launching the Server Agent (Terminal 2, 3, & 4)
* Open three other terminal windows and access the server agent with `ssh jc6` in each window.
* Clone the artifact repo on `jc6` if you haven't already, and `cd SmartCookie-Artifact/ebpf`.
* Run `./configure/configure_jc6.sh` once to configure static IP addresses and ARP entries.
* Next, use the provided python scripts in the separate terminals to compile and load the eBPF programs to the interface connected to the switch:
	* 1) `sudo python3 xdp_load.py enp3s0f1` for ingress 
	* 2) `sudo python3 tc_load.py enp3s0f1` for egress 
* You should see output that the programs have been loaded.
* Finally, run the following python script to sync timestamps between the server agent and switch agent, which is necessary for cookie checks: `sudo python3 send_ts.py`. 

### A Quick Functionality Test (Terminal 5 & 6) 
* To test a simple end-to-end connection between the `jc5` client and `jc6` server (protected by the intermediate switch agent and server agent), open two more terminals.
* SSH into the client with `ssh jc5` and SSH once more into the server with `ssh jc6`.
* On the server `jc6`, start up a `netcat` server with `nc -l -p 2020`.
* On the client `jc5`, connect to the `netcat` server with `nc 131.0.0.6 2020`.
* The client will seamlessly connect to the server after verification at the switch agent, and you can send messages between the client and server, with the messages popping up on the receiving side.
* If you are curious, you can use `tcpdump -evvvnX -i enp3s0f1` on both client and server to view the full packet sequence during connection setup, and map it to that of Figure 4 in the paper.
* Note that tcpdump is positioned after XDP on the _ingress_ pipeline, and after TC on the _egress_ pipeline (XDP-->tcpdump--> network stack on ingress, and network stack-->TC->tcpdump on egress).

## Evaluation Workflow 
There are three main experiments that showcase the key results and major claims of our work. These are described next. 

### Major Claims 

* C1: SMARTCOOKIE defends against attacks _without packet loss_ until rates of 136.9Mpps (which is 2.6x more than the next fastest defense). This is proven by experiment (E1), and described in Section 8.2 of the paper.
* C2: During attacks, SMARTCOOKIE protects benign clients from performance penalties and protects servers from additional CPU usage. It adds little to no latency overhead to benign connections during attacks, and any latency is comparable to the baseline latency with no ongoing attack. Additionally, it protects the server's CPU during attacks, fully keeping the CPU resources for other applications. This is proven by experiments (E2) and (E3), and shown in Section 8.3 and 8.4 of the paper.

## Experiment 1 - Hashing Throughput (Estimate: 45 human-minutes)
**Description:** Compare the maximum hashing throughput SMARTCOOKIE-HalfSipHash (SC-HSH) can achieve _without packet loss_ to the maximum hashing throughput of the three benchmarks: Kernel-SipHash (K-SH), XDP-HalfSipHash (XDP-HSH), and SMARTCOOKIE-AES (SC-AES). Use DPDK to send spoofed attack packets to the server and observe Rx to Tx packet rates to measure loss on the switch (for SC-AES and SC-AES) or server (for K-SH and XDP-HSH). (As noted in the paper, since our benchmarks performs one hash calculation per SYN packet, we effectively measure maximum hashing throughput.) The Tx rates should exactly match Rx rates for as long as SMARTCOOKIE or the benchmark is handling the attack without any packet loss. Once the defense begins to reach its capacity, the Tx (response) rate will begin to dip below Rx (received) rates.

**Preparation:** 
* Launch the switch agent, as described above. In three additional terminals, SSH into the attack machines: `ssh opti1`, `ssh opti2`, and `ssh jc4`. `DPDK` and `pktgen-DPDK` are already configured for you.
* For each attack terminal, `cd /home/shared/pktgen-dpdk` and launch pktgen with `sudo -E tools/run.py testbed`.
* If the server has been rebooted recently, reconfigure the huge pages: `cd /home/shared/dpdk/usertools` and run `./dpdk-setup.sh`. 
  	* On `opti1` and `opti2`, we have non-NUMA systems, so choose the option to setup hugepage mappings for non-NUMA systems [5]. Meanwhile, for `jc4`, choose the option for NUMA systems [52].
  	* Enter 8192 pages per node.
  	* Exit the script and return to the above steps to launch pktgen. 

**Execution:** 
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
* In the switch agent's `bf-sde.pm>` console, the command `rate-show` will also show Rx/Tx rates of the attack on the switch (ports `3/0`, `5/0`, and `6/0`). #### FIX ME!  

**Results:**
* To verify the maximum attack rate that SMARTCOOKIE can handle before any packet loss, increase the sending attack rate with `set 0 rate X` and `set 1 rate X`, with a maximum `X` of 100. As long as the Rx/Tx rates match in the switch agent, the switch agent is successfully defending against the SYN flood attack packets without any packet loss, up to 135 Mpps. 

## Experiment 2 - Latency (Estimate: 20 human-minutes)

## Experiment 3 - CPU (Estimate: 20 human-minutes)



### Benchmarking hash rate

The AES variant of the Switch Agent will respond to any incoming SYN packet from all non-server ports. To measure maximum hash rate, simply direct your packet generators to generate any TCP packet with TCP flags set to `0x02`, and increase sending rate (observe response packet rate) until loss is observed.

Note: for AES variant, please first run the controller script to load an arbitrary key; this is required to set up recirculation rounds correctly. 


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
