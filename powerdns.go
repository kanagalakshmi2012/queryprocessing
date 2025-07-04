package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	serverAddr = "127.0.0.1:53"
	domain     = "example.com."
	concurrency = 50
	testDuration = 1 * time.Second
)

var (
	totalQueries int64
	wg sync.WaitGroup
)

func buildQuery() []byte {
	return []byte{
		0xaa, 0xaa, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x07, 'e', 'x', 'a', 'm',
		'p', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, 0x00,
		0x01, 0x00, 0x01,
	}
}

func sendQuery() {
	defer wg.Done()
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		return
	}
	defer conn.Close()
	query := buildQuery()
	deadline := time.Now().Add(testDuration)
	for time.Now().Before(deadline) {
		_, err := conn.Write(query)
		if err != nil {
			continue
		}
		buf := make([]byte, 512)
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, err = conn.Read(buf)
		if err == nil {
			totalQueries++
		}
	}
}

func main() {
	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go sendQuery()
	}
	wg.Wait()
	elapsed := time.Since(start).Seconds()
	fmt.Printf("Total Queries: %d\n", totalQueries)
	fmt.Printf("QPS: %.2f\n", float64(totalQueries)/elapsed)
}
