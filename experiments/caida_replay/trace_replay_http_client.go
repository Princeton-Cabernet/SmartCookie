package main

import (
    "io"
    "bufio"
    "os"
    "fmt"
    "net"
    "net/http"
    "time"
    "strings"
    "strconv"
    "math/rand"
)

func rndAddr() (string) {
    //return fmt.Sprintf("169.254.0.1:%d",1025+rand.Intn(65533-1024))
    return fmt.Sprintf("130.0.%d.%d:%d",rand.Intn(254),rand.Intn(254),1024+rand.Intn(60000))
}

func http_client_use_addr(addr string) (*http.Client) {
    OKAddress, _ := net.ResolveTCPAddr("tcp", addr)

    transport := &http.Transport{
                 Dial: (&net.Dialer{
                         LocalAddr: OKAddress,
                     }).Dial }

    client := &http.Client{
         Transport: transport,
    }
    return client
}

func http_get(addr string, url string) (int,error){
    client:=http_client_use_addr(addr)
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
    return fmt.Sprintf("http://128.0.0.6:8090/%d",size)
}


type Request struct {
    fid int
    size int
    ScheduledStart time.Time
    ScheduledStartDelay time.Duration
}
type Result struct {
    fid int
    succ bool
    resLen int
    retries int
    elapsed time.Duration

    ActualStart time.Time
    ActualStartDelayed time.Duration
}

func http_get_with_retry(fid int, size int) Result {
    retries:=0
    succ:=false
    start:=time.Now()
    url:=get_url_from_size(size)
    addr:=rndAddr()//local ip:port address
    resLen:=-1
    var err error
    for i:=0;i<10;i++{
        resLen, err=http_get(addr,url)
        if err==nil{
            succ=true
            break
        }
        time.Sleep(time.Millisecond)
    }
    if !succ{
        fmt.Println("Failed after 10 retries!",err)
    }else if resLen!=size{
        fmt.Println("Anomaly response size!",resLen)
    }

    return Result{
        fid:fid,
        succ:succ,
        resLen:resLen,
        retries:retries,
        ActualStart:start,
    }
}


func worker_consume(myid int, req_channel <-chan Request, results_channel chan<- Result){
    fmt.Println("//Worker start,",myid)
    for {//loop forever
        req,ok := <-req_channel
        if !ok{
            //channel closed, we're done
            fmt.Println("//Worker done,",myid)
        }
	//if too early, wait a bit
	//tdiff:=req.ScheduledStart.Sub( time.Now() )
	//if tdiff>0 {time.Sleep(tdiff)}

        res:=http_get_with_retry(req.fid,req.size)
        results_channel<-res
    }
}

func read_trace(fn string) []Request {
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

    var AllRequest []Request
    i:=0
    for fileScanner.Scan() {
        line:=fileScanner.Text()
        arr:=strings.Split(line, "\t")
        if len(line)<2 || len(arr)!=5{
            continue
        }
        //fmt.Println(arr)
	i+=1
        begin,err:=strconv.ParseFloat(arr[1], 64)
        if err!=nil{panic(err)}
        size,err:=strconv.Atoi(arr[4])
        if err!=nil{panic(err)}

        beginNs:=int64(begin*(1e9))
	AllRequest=append(AllRequest,Request{
		fid:i,
		size:size,
		ScheduledStartDelay: time.Nanosecond * time.Duration(beginNs),
	})
        //fmt.Println(begin,size)
    }
    return AllRequest
}


func main(){
    reqs:=read_trace("caida_15sec_flows.csv")
    fmt.Println("Loaded trace, #flows=",len(reqs))

    N_workers:=1000
    N_jobs:=len(reqs)

    req_channel:=make(chan Request,N_jobs)
    results_channel:=make(chan Result,N_jobs)
    for i:=0;i<N_workers;i++{
        go worker_consume(i,req_channel,results_channel)
    }
    
    fmt.Println("//[main] start injection")
    begin:=time.Now()

    for _,r:= range reqs {
	r.ScheduledStart = begin.Add(r.ScheduledStartDelay)
        req_channel <- r
    }
    close(req_channel)

    //gather results
    ret:=0
    succ:=0
    fail:=0
    for j:=0;j<N_jobs;j++{
        res:=<-results_channel
        ret+=res.resLen
    	if res.succ {succ+=1} else {fail+=1}
    } 
    fmt.Println("Got all results",ret)
    end:=time.Now()
    fmt.Println("//[main] Time diff",end.Sub(begin)," for #conn=",N_jobs, "total succ/fail=", succ,fail)
}

