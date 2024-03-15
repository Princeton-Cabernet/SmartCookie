package main

import (
    "io"
    "fmt"
/*    "bufio"
    "os"
//    "log"
    "runtime"
    "strings"
    "strconv"
    "sync"
    */
    "net"
    "net/http"
    "time"
    "syscall"
    "golang.org/x/sys/unix"
    "math/rand"
)
var local_port string ="3043"
var remote_port string = "8090"
var server_dpid string = "128"

func check(e error) {
    if e != nil {
        panic(e)
    }
}
func getAddr() (string) {
	return fmt.Sprintf("130.0.0.8:"+local_port)
}
func rndAddr() (string) {
     _ = rand.Intn(1)
	//without randomized src port 
	return fmt.Sprintf("130.0.%d.%d:"+local_port,rand.Intn(255),rand.Intn(255))
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
func http_client_use_Addr() (*http.Client) {
/*	var (
 		clientsock int
		err        error
	)

	if clientsock, err = Socket(AF_INET, SOCK_STREAM, IPPROTO_IP); err != nil {
		fmt.Println("Client Socket() called error:", err.Error())
		return
	}
	if errReuse := SetsockoptInt(clientsock, SOL_SOCKET, SO_REUSEADDR, 1); errReuse != nil {
		fmt.Printf("reuse addr error: %v\n", errReuse)
		return
	}    

*/
    OKAddr := getAddr() // local IP address to use
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

func http_get(url string) (int,error){
    client:=http_client_use_Addr()
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
	//randomized 
	//return fmt.Sprintf("http://"+server_dpid+".0.%d.%d:8090/%d",rand.Intn(60), rand.Intn(255), size)
	return fmt.Sprintf("http://"+server_dpid+".0.0.6:8090/%d",size)
}

func http_get_retry_timing(size int) (bool, int, time.Duration){
    start := time.Now()
    succ:=false

    retries:=1
    for i:=0;i<retries;i++{
        retsize,err := http_get(get_url_from_size(size))
	if err != nil{
	    //fmt.Errorf("%v\n",err);
	    fmt.Println(err);
            time.Sleep(time.Millisecond)
        }
        if retsize>=size{
            succ=true
            retries=i
            break
        }
    }
    end := time.Now()
    return succ, retries, end.Sub(start)
}

func main(){
    succ,retr,dur:=http_get_retry_timing(14500)
    fmt.Println(succ,retr,dur)
    succ2,retr2,dur2:=http_get_retry_timing(14500)
    fmt.Println(succ2,retr2,dur2)
}

