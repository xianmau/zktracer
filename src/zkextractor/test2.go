package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

var cnt int = 1

func main() {
	fmt.Printf("Start.\n")
	timer := time.Tick(100 * time.Millisecond)
	for _ = range timer {
		doing([]string{"172.19.32.16"})
	}
}

func doing(host []string) {
	conn, _, err := zk.Connect(host, time.Second)
	if err != nil {
		panic(err)
	}
	_, _, err = conn.Get("/ymb")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cnt)
	cnt++
	fmt.Printf("done. %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
