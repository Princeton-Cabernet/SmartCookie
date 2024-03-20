package main

import (
    "io"
    "fmt"
    "os"
/*    "bufio"
//    "log"
    "strings"
    "math/rand"
    */
    "strconv"
    "net"
    "net/http"
    "runtime"
    "time"
    "sync"
    "syscall"
    "golang.org/x/sys/unix"
)
var remote_port string = "8090"
var server_dpid string = "131"

func check(e error) {
    if e != nil {
        panic(e)
    }
}
func getAddr(addr_a string, addr_b string, local_port string) (string) {
     //_ = rand.Intn(1)
	//without randomized src port 
	return fmt.Sprintf("129.0."+addr_a+"."+addr_b+":"+local_port)

     //return fmt.Sprintf("130.0.%d.%d:%d",rand.Intn(255),rand.Intn(255),1024+rand.Intn(65535-1024))
    //return fmt.Sprintf("130.%d.%d.%d:%d",rand.Intn(254),rand.Intn(254),rand.Intn(254),10000+rand.Intn(10000))
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
func http_client_use_Addr(addr_a string, addr_b string, local_port string) (*http.Client) {
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

    OKAddr := getAddr(addr_a, addr_b, local_port) // local IP address to use
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

func http_get(addr_a string, addr_b string, local_port string, url string) (int,error){
    client:=http_client_use_Addr(addr_a, addr_b, local_port)
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
	return fmt.Sprintf("http://"+server_dpid+".0.0.6:8090/%d",size)
}

func http_get_retry_timing(addr_a string, addr_b string, local_port string, size int) (bool, int, time.Duration){
    start := time.Now()
    succ:=false
    
    retries:=2
    for i:=0;i<retries;i++{
        retsize,err := http_get(addr_a, addr_b, local_port, get_url_from_size(size))
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

func start_run(wg *sync.WaitGroup, fid int, i int, j int, local_port string) () {
    defer wg.Done() 
    if(fid%200==0){fmt.Errorf("flow %d started\n",fid)}
    var addr_a string = strconv.Itoa(i)
    var addr_b string = strconv.Itoa(j)
    succ,retr,dur:=http_get_retry_timing(addr_a, addr_b, local_port, 1)
    fmt.Println(succ,retr,dur)

}

func main() {
    var local_port string = os.Args[1]
    var N = 10
    var M = 255 
    numCPUs := runtime.NumCPU()
    runtime.GOMAXPROCS(numCPUs+1)
    var wg sync.WaitGroup
    for i:=0;i<N;i++{
	for j:=0;j<M;j++{
    	    wg.Add(1)
     	    go start_run(&wg, i, i, j, local_port) // waitgroup, fid, addr_a, addr_b, local_port 
    	}
    }
    wg.Wait()
}

