package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func validURL(s string) bool {
	r := regexp.MustCompile("(?i)^http(?:s)?://")

	return r.MatchString(s)
}

func jsonGET(url string, wrapper interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	return dec.Decode(wrapper)
}

func httpGet(u string) ([]byte, error) {
	res, err := http.Get(u)
	if err != nil {
		return []byte{}, err
	}

	raw, err := ioutil.ReadAll(res.Body)

	res.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	return raw, nil
}

func cleanURL(u string, domain string) string {
	ind := strings.Index(u, domain)
	u = u[:ind+len(domain)]

	
	prefixes := []string{"http://", "https://", "www."}

	for _, p := range prefixes {
		u = strings.TrimPrefix(u, p)
	}

	return u
}
