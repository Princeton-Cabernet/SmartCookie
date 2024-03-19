# SmartCookie Artifact 

This repository contains the prototype source code and instructions for artifact evaluation for our USENIX Security'24 paper [SmartCookie: Blocking Large-Scale SYN Floods with a Split-Proxy Defense on Programmable Data Planes](#).

## Contents
The artifact consists of two major pieces: 1) source code for the switch agent and server agent of SmartCookie's split-proxy SYN-flooding defense, related benchmark code, and measurement scripts (showing availability), and 2) a hardware testbed for running and evaluating SmartCookie under key attack scenarios (showing functionality and reproducibility). 
* `p4src/` includes the Switch Agent program that calculates SYN cookies using HalfSipHash.
	* `p4src/benchmark/` contains variants of the Switch Agent, for benchmarking max hashing rate using a different hash function (AES).
* `ebpf/` includes the Server Agent programs that process cookie-verified new connection handshake and false positive packets.
	* `ebpf/benchmark/` contains a XDP-based SYN cookie generator, for benchmarking max hashing rate of a server-only solution.
* `experiments/` includes the relevant scripts for running key experiments.
	* `experiments/measurements/` contains scripts for collecting client-side latency and server-side CPU measurements.

## Description & Requirements
For the purposes of this artifact evaluation, our testbed consists of five servers and an Intel Tofino Wedge32X-BF programmable switch.
Three machines act as adversaries, each with a \textit{XX}-core \textit{Intel Xeon Silver 4114 CPU and a Mellanox ConnectX-5 2x100Gbps NIC}, generating attack traffic using DPDK 19.12.0 and pktgen-DPDK.
Two other machines act as server and client, each with 8-core Intel Xeon D-1541 CPUs and Intel X552 2x10Gbps NICs. 
**For simplicity of artifact evaluation, we are providing evaluators with access to our preconfigured testbed (access instructions below). Instructions for installations and dependencies are briefly included for completeness, but all installations and dependencies are already in place for the evaluation testbed.** 
Next, we describe how to access the testbed, what hardware and software dependencies are required (these are preconfigured for the testbed), and what additional benchmarks can be run. 

### Security, privacy, and ethical concerns
There are no security, privacy, or ethical concerns or risks to evaluators or their machines. All experiments can be run on the authors’ testbed, which is provisioned for the planned attack rates. For testbed access, please do not share or distribute the private key (discussed further below).

### Accessing the testbed 
Save the SSH private access key (shared with you
directly on the submission site) to your local ma-
chine under ~/.ssh/usenixsec24ae.priv.id_rsa.
Update the permissions with chmod 600
~/.ssh/usenixsec24ae.priv.id_rsa. Start the ssh-
agent and load the key: eval $(ssh-agent -s) and
ssh-add ~/.ssh/usenixsec24ae.priv.id_rsa.
• Put the following text into your local machine’s
~/.ssh/config, such that you can ssh into the machines
by hostname using the public-facing proxy port. Your
public keys are already in place.


## Usage

### Loading the Server Agent

Prerequisite: please use kernel 5.10 or newer and install the entire `bcc` toolkit.
(For Ubuntu, you may run `sudo apt-get install bpfcc-tools python3-bpfcc linux-headers-$(uname -r)`)

Use the provided python scripts to compile and load the eBPF programs to the interface connected to the programmable switch:

1. `sudo python xdp_load.py *if_name*` for ingress path
2. `sudo python tc_load.py *if_name*` for egress path

### Loading the Switch Agent

Prerequisite: please use `bf-sde` version 9.7.1 or newer to compile the P4 program. Then, the program can be deployed via `$SDE/run_switchd.sh -p SmartCookie-HalfSipHash`.

### Benchmarking hash rate

All variants of the Switch Agent will respond to any incoming SYN packet from all non-server ports. To measure maximum hash rate, simply direct your packet generator to generate any TCP packet with TCP flags set to `0x02`, and increase sending rate (observe response packet rate) until loss is observed.

Note: for AES variant, please first run the controller script to load an arbitrary key; this is required to set up recirculation rounds correctly. 


## Citing
If you find this implementation or our paper useful, please consider citing:

    @inproceedings{yoo2023smartcookie,
        title={SmartCookie: Blocking Large-Scale SYN Floods with a Split-Proxy Defense on Programmable Data Planes},
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
