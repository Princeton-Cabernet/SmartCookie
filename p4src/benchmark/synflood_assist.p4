
#include <core.p4>
#include <tna.p4>

typedef bit<48> mac_addr_t;
typedef bit<32> ipv4_addr_t;
typedef bit<16> ether_type_t;
const ether_type_t ETHERTYPE_IPV4 = 16w0x0800;
const ether_type_t ETHERTYPE_VLAN = 16w0x0810;

typedef bit<8> ip_protocol_t;
const ip_protocol_t IP_PROTOCOLS_ICMP = 1;
const ip_protocol_t IP_PROTOCOLS_TCP = 6;
const ip_protocol_t IP_PROTOCOLS_UDP = 17;


header ethernet_h {
    mac_addr_t dst_addr;
    mac_addr_t src_addr;
    bit<16> ether_type;
}

header ipv4_h {
    bit<4> version;
    bit<4> ihl;
    bit<8> diffserv;
    bit<16> total_len;
    bit<16> identification;
    bit<3> flags;
    bit<13> frag_offset;
    bit<8> ttl;
    bit<8> protocol;
    bit<16> hdr_checksum;
    ipv4_addr_t src_addr;
    ipv4_addr_t dst_addr;
}

header tcp_h {
    bit<16> src_port;
    bit<16> dst_port;

    bit<32> seq_no;
    bit<32> ack_no;
    bit<4> data_offset;
    bit<4> res;
    bit<8> flags;
    bit<16> window;
    bit<16> checksum;
    bit<16> urgent_ptr;
}

header udp_h {
    bit<16> src_port;
    bit<16> dst_port;
    bit<16> udp_total_len;
    bit<16> checksum;
}
header udp_padding_h{
   bit<32> a; 
   bit<32> b; 
   bit<32> c; 
   bit<32> d; 
  
}
struct header_t {
    ethernet_h ethernet;
    ipv4_h ipv4;
    tcp_h tcp;
    udp_h udp;
    udp_padding_h udp_padding; 
}

struct ig_metadata_t {
    bit<16> tcp_total_len;
    bit<1> redo_cksum; 	
}
struct eg_metadata_t {
}


parser TofinoIngressParser(
        packet_in pkt,
        inout ig_metadata_t ig_md,
        out ingress_intrinsic_metadata_t ig_intr_md) {
    state start {
        pkt.extract(ig_intr_md);
        transition select(ig_intr_md.resubmit_flag) {
            1 : parse_resubmit;
            0 : parse_port_metadata;
        }
    }

    state parse_resubmit {
        // Parse resubmitted packet here.
        pkt.advance(64);
        transition accept;
    }

    state parse_port_metadata {
        pkt.advance(64);  //tofino 1 port metadata size
        transition accept;
    }
}
parser SwitchIngressParser(
        packet_in pkt,
        out header_t hdr,
        out ig_metadata_t ig_md,
        out ingress_intrinsic_metadata_t ig_intr_md) {

    TofinoIngressParser() tofino_parser;

    state start {
        tofino_parser.apply(pkt, ig_md, ig_intr_md);
        transition parse_ethernet;
    }

    state parse_ethernet {
        pkt.extract(hdr.ethernet);
        transition select (hdr.ethernet.ether_type) {
            ETHERTYPE_IPV4 : parse_ipv4;
            default : reject;
        }
    }

    state parse_ipv4 {
        pkt.extract(hdr.ipv4);
        transition select(hdr.ipv4.protocol) {
            IP_PROTOCOLS_TCP : parse_tcp;
            IP_PROTOCOLS_UDP : parse_udp;
            default : accept;
        }
    }

    state parse_tcp {
        pkt.extract(hdr.tcp);
        transition select(hdr.ipv4.total_len) {
            default : accept;
        }
    }

    state parse_udp {
        pkt.extract(hdr.udp);
	pkt.extract(hdr.udp_padding);
        transition select(hdr.udp.dst_port) {
            default: accept;
        }
    }
}

// ---------------------------------------------------------------------------
// Ingress Deparser
// ---------------------------------------------------------------------------
control SwitchIngressDeparser(
        packet_out pkt,
        inout header_t hdr,
        in ig_metadata_t ig_md,
        in ingress_intrinsic_metadata_for_deparser_t ig_intr_dprsr_md) {


      Checksum() ipv4_checksum;
      Checksum() tcp_checksum;

    apply {
	if(ig_md.redo_cksum == 0x1){

	        hdr.ipv4.hdr_checksum = ipv4_checksum.update({
	            hdr.ipv4.version,
	            hdr.ipv4.ihl,
	            hdr.ipv4.diffserv,
	            hdr.ipv4.total_len,
	            hdr.ipv4.identification,
	            hdr.ipv4.flags,
	            hdr.ipv4.frag_offset,
	            hdr.ipv4.ttl,
	            hdr.ipv4.protocol,
	            hdr.ipv4.src_addr,
	            hdr.ipv4.dst_addr
	        });
	
	        hdr.tcp.checksum = tcp_checksum.update({
	            //==pseudo header
                    hdr.ipv4.src_addr,
                    hdr.ipv4.dst_addr,
                    8w0,
                    hdr.ipv4.protocol,
                    ig_md.tcp_total_len,
	            //==actual header
                    hdr.tcp.src_port,
                    hdr.tcp.dst_port,
                    hdr.tcp.seq_no,
                    hdr.tcp.ack_no,
                    hdr.tcp.data_offset,
                    hdr.tcp.res,
                    hdr.tcp.flags,
                    hdr.tcp.window,
                    hdr.tcp.urgent_ptr
                });
	}

        pkt.emit(hdr.ethernet);
        pkt.emit(hdr.ipv4);
        pkt.emit(hdr.tcp);
        pkt.emit(hdr.udp);
    }
}

// ---------------------------------------------------------------------------
// Egress parser
// ---------------------------------------------------------------------------
parser SwitchEgressParser(
        packet_in pkt,
        out header_t hdr,
        out eg_metadata_t eg_md,
        out egress_intrinsic_metadata_t eg_intr_md) {
    state start {
        pkt.extract(eg_intr_md);
        transition accept;
    }
}

// ---------------------------------------------------------------------------
// Egress Deparser
// ---------------------------------------------------------------------------
control SwitchEgressDeparser(
        packet_out pkt,
        inout header_t hdr,
        in eg_metadata_t eg_md,
        in egress_intrinsic_metadata_for_deparser_t eg_intr_md_for_dprsr) {
    apply {
    }
}


// ---------------------------------------------------------------------------
// Ingress Control
// ---------------------------------------------------------------------------
control SwitchIngress(
        inout header_t hdr,
        inout ig_metadata_t ig_md,
        in ingress_intrinsic_metadata_t ig_intr_md,
        in ingress_intrinsic_metadata_from_parser_t ig_intr_prsr_md,
        inout ingress_intrinsic_metadata_for_deparser_t ig_intr_dprsr_md,
        inout ingress_intrinsic_metadata_for_tm_t ig_intr_tm_md) {

        action drop() {
            ig_intr_dprsr_md.drop_ctl = 0x1; // Drop packet.
        }
        action nop() {
        }
	action route_to(bit<9> dst){
		ig_intr_tm_md.ucast_egress_port=dst;
	}
        action reflect(){
            //send you back to where you're from
            ig_intr_tm_md.ucast_egress_port=ig_intr_md.ingress_port;
        }

        action fill_ip_header(){
            hdr.ipv4.version=4;
            hdr.ipv4.ihl=5;
            hdr.ipv4.diffserv=0;
            //len: assigned by TCP-related stuff
            hdr.ipv4.identification=11234;
            hdr.ipv4.flags=2;
            hdr.ipv4.frag_offset=0;
            hdr.ipv4.ttl=64;
            hdr.ipv4.protocol=6;
        }
        action fill_ip_header_tcp_no_payload(){
            //assume a 64-bit/8-byte payload
            hdr.ipv4.total_len=5*4+5*4;
            ig_md.tcp_total_len=5*4;
        }

        action fill_tcp_common(){
            hdr.tcp.data_offset=5;
            hdr.tcp.res=0;
            hdr.tcp.window=16384;
            hdr.tcp.urgent_ptr=0;
        }



        action fill_ether_header(){
	    //jc4
            hdr.ethernet.src_addr=0x3cfdfeccb1c0;
            //jc5
	    hdr.ethernet.dst_addr=0xd05099e828cb;
        }
        action fill_tcp_syn_packet(bit<32> seq_no){
		//this may not be necessary
                //fill_ether_header();
                fill_ip_header();
                fill_ip_header_tcp_no_payload();
                fill_tcp_common();

            hdr.tcp.seq_no=seq_no;
            hdr.tcp.ack_no=0;
            hdr.tcp.flags=2;//TCP_FLAGS_S;

            hdr.tcp.src_port=hdr.udp.src_port;
            hdr.tcp.dst_port=hdr.udp.dst_port;
        }



        apply {

            	//for testing with jc4 (dpid 160 and 168), jc5 (dpid 129), jc6 (dpid 131), opti1 (dpid 144), and opti2 (dpid 152)  
		//dst ip based routing 

		//dst jc4 port 0, ip 168.0.0.4
		if(hdr.ipv4.dst_addr == 0xA8000004){
			route_to(168); 
	    		hdr.ethernet.dst_addr=0x0000000000A8;
		}
		//dst jc4 port 1, ip 160.0.0.4
		else if(hdr.ipv4.dst_addr == 0xA0000004){
			route_to(160); 
	    		hdr.ethernet.dst_addr=0x0000000000A0;
		}
		//dst jc5, ip 129.0.0.5
		else if(hdr.ipv4.dst_addr == 0x81000005){
			route_to(129); 
	    		hdr.ethernet.dst_addr=0x000000000081;
		}
		//dst jc6, ip 131.0.0.6
		else if(hdr.ipv4.dst_addr == 0x83000006){
			route_to(131); 
	    		hdr.ethernet.dst_addr=0x000000000083;
		}
		//dst opti1, ip 144.0.0.7
		else if(hdr.ipv4.dst_addr == 0x90000007){
			route_to(144); 
	    		hdr.ethernet.dst_addr=0x000000000090;
		}
		else{ //dst opti2
			route_to(152);
            		hdr.ethernet.dst_addr=0x000000000098;
		}

		if(hdr.udp.isValid()){
			hdr.tcp.setValid();
			fill_tcp_syn_packet(hdr.ipv4.src_addr);
			ig_md.redo_cksum = 0x1; 
			hdr.udp.setInvalid();
		}

        }
}

// ---------------------------------------------------------------------------
// Egress Control
// ---------------------------------------------------------------------------
control SwitchEgress(
        inout header_t hdr,
        inout eg_metadata_t eg_md,
        in egress_intrinsic_metadata_t eg_intr_md,
        in egress_intrinsic_metadata_from_parser_t eg_intr_md_from_prsr,
        inout egress_intrinsic_metadata_for_deparser_t ig_intr_dprs_md,
        inout egress_intrinsic_metadata_for_output_port_t eg_intr_oport_md) {
    apply {
    }
}



Pipeline(SwitchIngressParser(),
         SwitchIngress(),
         SwitchIngressDeparser(),
         SwitchEgressParser(),
         SwitchEgress(),
         SwitchEgressDeparser()
         ) pipe;

Switch(pipe) main;
