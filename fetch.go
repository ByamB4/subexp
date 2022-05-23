package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func fetchHackerTarget(domain string) ([]string, error) {
	ret := make([]string, 0)
	raw, err := httpGet(
		fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain),
	)

	if err != nil {
		return ret, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), ",", 2)
		if len(parts) != 2 {
			continue
		}
		ret = append(ret, parts[0])
	}
	return ret, sc.Err()
}
