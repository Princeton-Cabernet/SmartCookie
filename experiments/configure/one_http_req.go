package main

import (
    "io"
    "fmt"
    "os"
/*    "bufio"
//    "log"
    "strings"
    "math/rand"
    "runtime"
    "sync"
    "strconv"
    */
    "net"
    "net/http"
    "time"
    "syscall"
    "golang.org/x/sys/unix"
)
var remote_port string = "8090"
var server_dpid string = "128"
func check(e error) {
    if e != nil {
        panic(e)
    }
}
func getAddr(srcport string, srcip string) (string) {
	return fmt.Sprintf(srcip+":"+srcport)
}
func SocketControl(network, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(fd uintptr) {
		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
        check(err)

		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
        check(err)
        
        var l unix.Linger
        
        l.Linger=0
        l.Onoff=0
        err = unix.SetsockoptLinger(int(fd), unix.SOL_SOCKET, unix.SO_LINGER, &l)
        check(err)
	})
	return err
}
func http_client_use_Addr(srcport string, srcip string) (*http.Client) {

    OKAddr := getAddr(srcport, srcip) // local IP address to use
    OKAddress, _ := net.ResolveTCPAddr("tcp", OKAddr)
    transport := &http.Transport{
                 Dial: (&net.Dialer{
        		 Control: SocketControl,
			 LocalAddr: OKAddress,
                     }).Dial }

    client := &http.Client{
         Transport: transport,
    }
    return client

}

func http_get(srcport string, srcip string, url string) (int,error){
    client:=http_client_use_Addr(srcport, srcip)
    resp, err := client.Get(url)
    if err != nil {
        return 0,err
    }
    defer resp.Body.Close()

    data := make([]byte, 4096)
    total:=0
    for {
        n, err := resp.Body.Read(data)
        total+=n
        if err != nil {
            if err == io.EOF {
                return total,nil
            }
            return -1,err
        }
    }
}

func get_url_from_size(size int) string{
	return fmt.Sprintf("http://"+server_dpid+".0.0.6:8090/%d", size)
}

func http_get_retry_timing(srcport string, srcip string, size int) (bool, int, time.Duration){
    start := time.Now()
    succ:=false
    //force second retry in one ms 
    retries:=8
    for i:=0;i<retries;i++{
        retsize,err := http_get(srcport, srcip, get_url_from_size(size))
	if err != nil{
	    //fmt.Errorf("%v\n",err);
	    fmt.Println(err);
            time.Sleep(time.Millisecond)
        } else if retsize>=size{
            succ=true
            retries=i
            break
        }
    }
    end := time.Now()
    return succ, retries, end.Sub(start)
}


// make a single http request with provided srcport, srcip, and payload size (that automatically retries once upon RST) 
func main() {
    var srcip string = os.Args[1] 
    var srcport string = os.Args[2] 
    var size int = 14500  
    succ,retr,dur:=http_get_retry_timing(srcport, srcip, size)
    fmt.Println(succ,retr,dur)
    
/*    if retr == 1 { //we were successful, the timing measurement is accurate 
    	fmt.Println(succ,retr,dur)
    } 
    if retr == 2{ //we weren't successful after one retry (i.e., packets were dropped, server overloaded). print a large time to signify timeout. 
    	fmt.Println(succ,retr,9999.99999)
    }
*/
}

