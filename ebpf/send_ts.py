import time 
import socket 
import os 
#from ctypes import c_int32 

#set ip to client ip to direct to switch agent, which will process and drop pkt 
UDP_IP = "129.0.0.5"
UDP_PORT = 5555; 
period = 2

print("Sending one timesync UDP pkt every", period, "seconds, hit CTRL+C to stop")

#f = open('/proc/uptime', 'r')

#def int32(val):
#    return c_int32(val).value

while 1:
    try:
        #get total number of seconds since system boot 
        #ts = os.popen('uptime -p').read()[:-1]
        #ts = float(f.readline().split()[0])
        print("Sending UDP pkt")
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) 
        sock.sendto(bytes(4), (UDP_IP, UDP_PORT))
        time.sleep(period)
    except KeyboardInterrupt:
        print("\nEnding timesync pkt send")
        break 
