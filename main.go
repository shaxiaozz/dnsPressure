package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var concurrency int

func init() {
	flag.IntVar(&concurrency, "c", 1, "requests per second")
}

func resolveDomain(domain string) (time.Duration, error) {
	resolver := &net.Resolver{
		PreferGo: true,
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := resolver.LookupHost(ctx, domain)
	duration := time.Since(start)
	return duration, err
}

func logResult(domain string, duration time.Duration, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05,000")
	if err != nil {
		fmt.Printf("%s %s fail %v %s\n", timestamp, domain, duration.Truncate(time.Millisecond), parseDNSError(err))
	} else {
		fmt.Printf("%s %s success %v\n", timestamp, domain, duration.Truncate(time.Millisecond))
	}
}

func parseDNSError(err error) string {
	msg := err.Error()
	if strings.Contains(msg, "i/o timeout") {
		return "timeout"
	}
	return msg
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: dnsPressure -c <concurrency> <domain>")
		os.Exit(1)
	}
	domain := flag.Arg(0)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup

	for {
		<-ticker.C
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				duration, err := resolveDomain(domain)
				logResult(domain, duration, err)
			}()
		}

		wg.Wait()
	}
}
