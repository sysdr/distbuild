// kvbench hammers a kvserver with concurrent clients and reports throughput.
//
//	go run ./cmd/kvbench -addr localhost:7070 -clients 50 -ops 20000
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	addr := flag.String("addr", "localhost:7070", "server address")
	clients := flag.Int("clients", 50, "concurrent connections")
	ops := flag.Int("ops", 20000, "SET+GET pairs per client")
	flag.Parse()

	var wg sync.WaitGroup
	start := time.Now()
	for c := 0; c < *clients; c++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", *addr)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			r := bufio.NewReader(conn)
			w := bufio.NewWriter(conn)
			for i := 0; i < *ops; i++ {
				fmt.Fprintf(w, "SET k%d:%d v%d\nGET k%d:%d\n", c, i, i, c, i)
				w.Flush()
				r.ReadString('\n') // OK
				r.ReadString('\n') // value
			}
		}(c)
	}
	wg.Wait()
	elapsed := time.Since(start)
	total := *clients * *ops * 2
	fmt.Printf("%d requests from %d clients in %v -> %.0f req/s\n",
		total, *clients, elapsed.Round(time.Millisecond), float64(total)/elapsed.Seconds())
}
