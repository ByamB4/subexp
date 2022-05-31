package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

type fetchFn func() ([]string, error)

func init() {
	flag.Usage = func() {
		h := "Sub domain explorer tool.\n\n"

		h += "Examples: \n"
		h += "	subexp unitel.mn\n"

		fmt.Fprintf(os.Stderr, h)
	}
}

func main() {
	flag.Parse()
	domain := strings.ToLower(flag.Arg(0))

	if domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	f := fetch{domain}

	sources := []fetchFn{
		f.urlScan,
		f.bufferOverrun,
		f.crtSh,
		f.hackerTarget,
		f.certSpotter,
		f.wayArchive,
	}

	out := make(chan string)
	var wg sync.WaitGroup

	for _, source := range sources {
		wg.Add(1)
		go func(source fetchFn) {
			defer wg.Done()
			subdomains, err := source()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				return
			}
			for _, subdomain := range subdomains {
				out <- cleanDomain(subdomain)
			}
		}(source)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	result := []string{}

	printed := make(map[string]bool)

	for n := range out {
		if _, ok := printed[n]; ok {
			continue
		}
		printed[n] = true
		result = append(result, n)
	}

	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) < len(result[j])
	})

	for _, subdomain := range result {
		fmt.Println(subdomain)
	}
}

func cleanDomain(d string) string {
	d = strings.ToLower(d)

	if len(d) < 2 {
		return d
	}

	if d[0] == '*' || d[0] == '%' {
		d = d[1:]
	}

	if d[0] == '.' {
		d = d[1:]
	}

	if d[:4] == "www1" {
		return d
	}

	if d[:3] == "www" {
		d = d[4:]
	}

	return d

}
