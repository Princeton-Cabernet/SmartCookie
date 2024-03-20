package main

import (
    "fmt"
    "net/http"
    "runtime"
    "strings"
    "strconv"
)

const N=4096
var block = []byte(strings.Repeat("0",N-1)+"\n")

func hello(w http.ResponseWriter, r *http.Request) {
    //fmt.Println(r.RemoteAddr)

    i, err:=strconv.Atoi(r.URL.Path[1:])
    if err!=nil || i<=0 {
        fmt.Fprintf(w, "Error: cannot parse a number from Path %s", r.URL.Path)
        return
    }
    for i>N{
        w.Write(block)
        i-=N
    }
    w.Write(block[:i])
}

func main() {
    fmt.Println("Launching server and listening for connection requests...")
    runtime.GOMAXPROCS(38)
    http.HandleFunc("/", hello)
    http.ListenAndServe(":8090", nil)
}
