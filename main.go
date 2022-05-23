package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func init() {
	flag.Usage = func() {
		h := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n\n"

		h += "Examples: \n"
		h += "	subexp http://example.com\n"

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

	domains, err := fetchHackerTarget(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("%s\n", strings.Join(domains, "\n"))
	}
}
