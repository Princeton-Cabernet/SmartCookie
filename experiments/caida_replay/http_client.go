package main

import (
    "io"
    "bufio"
    "os"
    "fmt"
//    "log"
    "net"
    "net/http"
    "runtime"
    "time"
    "strings"
    "strconv"
    "sync"
    "syscall"
    "golang.org/x/sys/unix"
    "math/rand"
)
var local_addr1 string =  "129.0.0.1"
var local_addr2 string =  "129.0.0.2"
var local_addr3 string =  "129.0.0.3"
var local_addr4 string =  "129.0.0.4"
var local_addr5 string =  "129.0.0.5"
var local_addr6 string =  "129.0.0.6"
var local_addr7 string =  "129.0.0.7"
var local_port string ="2020"
var remote_addr1 string = "36.0.0.1"
var remote_addr2 string = "36.0.0.2"
var remote_port string = "8090"
var server_dpid string = "131"

func check(e error) {
    if e != nil {
        panic(e)
    }
}
func rndAddr() (string) {
     _ = rand.Intn(1)
     //return fmt.Sprintf(local_addr1+":"+local_port)
/*     var rnd_ctr int = rand.Intn(140)
     if(rnd_ctr > 120){
     	return fmt.Sprintf(local_addr1+":%d",1024+rand.Intn(65535-1024))
     } else if (rnd_ctr > 100){
     	return fmt.Sprintf(local_addr2+":%d",1024+rand.Intn(65535-1024))
     } else if (rnd_ctr > 80){
     	return fmt.Sprintf(local_addr3+":%d",1024+rand.Intn(65535-1024))
     } else if (rnd_ctr >60){
     	return fmt.Sprintf(local_addr4+":%d",1024+rand.Intn(65535-1024))
     } else if (rnd_ctr > 40){
     	return fmt.Sprintf(local_addr5+":%d",1024+rand.Intn(65535-1024))
     } else if (rnd_ctr > 20){
     	return fmt.Sprintf(local_addr6+":%d",1024+rand.Intn(65535-1024))
     } else {
     	return fmt.Sprintf(local_addr7+":%d",1024+rand.Intn(65535-1024))
     }
*/   
	//without randomized src port 
	return fmt.Sprintf("129.0.%d.%d:local_port",rand.Intn(255),rand.Intn(255))

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
func http_client_use_rndAddr() (*http.Client) {
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

    OKAddr := rndAddr() // local IP address to use
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
    client:=http_client_use_rndAddr()
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
        /*var rnd_ctr int = rand.Intn(140)
	if (rnd_ctr > 70) {
	   return fmt.Sprintf("http://"+remote_addr1+":"+remote_port+"/%d",size)
        } else {
	   return fmt.Sprintf("http://"+remote_addr2+":"+remote_port+"/%d",size)
    	}*/
	return fmt.Sprintf("http://"+server_dpid+".0.%d.%d:8090/%d",rand.Intn(60), rand.Intn(255), size)
}

func http_get_retry_timing(size int) (bool, int, time.Duration){
    start := time.Now()
    succ:=false

    retries:=20
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

/* 
func main() {
    succ,retr,dur:=http_get_retry_timing(12345)
    fmt.Println(succ,retr,dur)
}*/


func parse_flow_schedule(fn string) ([]int64, []int){
    readFile, err := os.Open(fn)
    defer readFile.Close()
    if err != nil {
        panic(err)
    }

    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)

    //skip first line
    fileScanner.Scan()
    _=fileScanner.Text()

    var allsizes []int
    var allbeginNS []int64  

    for fileScanner.Scan() {
        line:=fileScanner.Text()
        arr:=strings.Split(line, "\t")
        if len(line)<2 || len(arr)!=5{
            continue
        }
        //fmt.Println(arr)
        begin,err:=strconv.ParseFloat(arr[1], 64)
        if err!=nil{panic(err)}
        size,err:=strconv.Atoi(arr[4])
        if err!=nil{panic(err)}

        beginNs:=int64(begin*(1e9))

        allbeginNS=append(allbeginNS,beginNs)
        allsizes=append(allsizes,size)
        //fmt.Println(begin,size)
    }
    return allbeginNS,allsizes
}

type Result struct {
    fid int
    succ bool
    retries int
    elapsed time.Duration

    ActualStart time.Time
    ActualStartDelayed time.Duration
}

func wait_and_run(wg *sync.WaitGroup, fid int, results []Result, beginNS int64, oldbeginNS int64, size int){
    defer wg.Done()
    time.Sleep((time.Duration(beginNS)-time.Duration(oldbeginNS))*time.Nanosecond)
    if(fid%200==0){fmt.Errorf("flow %d started\n",fid)}
    start:=time.Now()
    succ,retr,dur:=http_get_retry_timing(size)

    results[fid]=Result{fid: fid, succ: succ, retries: retr, elapsed: dur, ActualStart:start}
}
func directly_run(wg *sync.WaitGroup, fid int, results []Result, beginNS int64, oldbeginNS int64, size int){
    defer wg.Done()
    //fmt.Println("Goroutine executing for fid ", fid)
    if(fid%200==0){fmt.Errorf("flow %d started\n",fid)}
    start:=time.Now()
    succ,retr,dur:=http_get_retry_timing(size)

    results[fid]=Result{fid: fid, succ: succ, retries: retr, elapsed: dur, ActualStart:start}
}
func main(){
    //fn:="traces/caida_15s_flows.csv" // full trace 
    //fn:="traces/1percent_sample.csv" // percentage of 1 min trace 
    fn:="traces/caida10.csv" // n flows
    allbeginNS,allsizes:=parse_flow_schedule(fn)
    N:=len(allbeginNS)
    results:=make([]Result, N)
    numCPUs := runtime.NumCPU()
    runtime.GOMAXPROCS(numCPUs+1)
    globalStart:=time.Now()
    t0 := globalStart.UnixNano()
    var start_time int64 = int64(t0)
    var wg sync.WaitGroup
    
    //dummy line 
    allsizes[0] = 0 
    for i:=0;i<N;i++{
        //fmt.Println("Sleep for ", start_time + allbeginNS[i] - int64(time.Now().UnixNano()))
        wg.Add(1)
        //hard code the size for fast injection 
	go directly_run(&wg, i, results, allbeginNS[i], 0, 1)
        //go directly_run(&wg, i, results, allbeginNS[i], 0, allsizes[i])
        time.Sleep(time.Duration( start_time + allbeginNS[i] - int64(time.Now().UnixNano()) )*time.Nanosecond)
   }
    wg.Wait()
    
    fmt.Println("fid\t succ \t num_retries \t latency(us) \t schedBegin(us) \t actualBegin(us)")
    for i:=0;i<N;i++{
        results[i].ActualStartDelayed = results[i].ActualStart.Sub(globalStart)        
        r:=results[i]
        fmt.Printf("%d \t %v \t %v \t %v \t %v \t %v \n", r.fid, r.succ, r.retries, r.elapsed.Microseconds(),
            allbeginNS[i]/1000, r.ActualStartDelayed.Microseconds()    );
    }
}
