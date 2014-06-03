package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"zk"
)

var (
	zones = map[string][]string{
		"zone1": []string{"192.168.56.101"},
		//"zone1": []string{"172.19.32.16"},
		//"zone2": []string{"172.19.32.153"},
		//"zone3": []string{"172.19.32.46"},
	}
)

func main() {

	// xmzk
	conn := zk.New()
	conn.Connect(zones["zone1"])
	t1 := time.Now()
	pdo(conn, "/ymb/loggers")
	//traverse(conn, "/ymb/loggers")
	t2 := time.Now()
	log.Printf("Extract data from [%s] using %v\n", "zone1", t2.Sub(t1))
}

func pdo(conn *zk.ZkCli, path string) {
	var wg sync.WaitGroup
	children, _ := conn.Children(path)
	for _, znode := range children {
		wg.Add(1)
		go func(path string) {
			_, err := conn.Get(path)
			if err == nil {
				//fmt.Println(string(data))
			} else {
				fmt.Println(err)
			}
			wg.Done()
		}(path + "/" + znode)
	}
	wg.Wait()
}

// traverse all znodes under the specified path
func traverse(conn *zk.ZkCli, path string) {
	children, err := conn.Children(path)
	if err != nil {
		panic(err)
	}

	if len(children) <= 0 {
		_, err := conn.Get(path)
		if err == nil {
			//fmt.Println("#Leaf ZNode Found:")
			//fmt.Println("#PATH: ", path)
			//fmt.Println("#DATA: ", string(data))
		}
	}
	for _, znode := range children {
		if path == "/" {
			//fmt.Printf("Searching ZNode: /%s\n", znode)
			traverse(conn, "/"+znode)
		} else {
			//fmt.Printf("Searching ZNode: %s/%s\n", path, znode)
			traverse(conn, path+"/"+znode)
		}
	}
}
