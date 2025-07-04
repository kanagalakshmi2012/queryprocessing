package main
import (
	"fmt"
	"net"
	"time"
	"sync"
)
const (
	targetDNS  = "127.0.0.1:53"
	domainName = "example.com."
	workers    = 50
	duration   = 1 * time.Second
)
var (
	queryCount int64
	wg         sync.WaitGroup
)
func dnsQuery() {
	defer wg.Done()
	msg := []byte{
		0xaa, 0xaa, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x07, 'e', 'x', 'a', 'm',
		'p', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, 0x00,
		0x01, 0x00, 0x01,
	}
	conn, err := net.Dial("udp", targetDNS)
	if err != nil {
		return
	}
	defer conn.Close()
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		_, err := conn.Write(msg)
		if err != nil {
			continue
		}
		buffer := make([]byte, 512)
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, err = conn.Read(buffer)
		if err == nil {
			queryCount++
		}
	}
}
func main() {
	start := time.Now()
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go dnsQuery()
	}
	wg.Wait()
	elapsed := time.Since(start).Seconds()
	fmt.Printf("Total queries: %d\n", queryCount)
	fmt.Printf("QPS: %.2f\n", float64(queryCount)/elapsed)
}
