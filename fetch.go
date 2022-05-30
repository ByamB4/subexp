package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type fetch struct {
	domain string
}

func (f fetch) hackerTarget() ([]string, error) {
	ret := make([]string, 0)
	raw, err := httpGet(
		fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", f.domain),
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

type CrtShRet struct {
	Name string `json:"name_value"`
}

func (f fetch) crtSh() ([]string, error) {
	var ret []CrtShRet

	resp, err := http.Get(
		fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", f.domain),
	)

	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()

	output := make([]string, 0)

	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &ret); err != nil {
		return []string{}, err
	}

	for _, res := range ret {
		output = append(output, strings.Split(res.Name, "\n")...)
	}

	return output, nil
}

func (f fetch) urlScan() ([]string, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", f.domain),
	)

	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()

	output := make([]string, 0)

	dec := json.NewDecoder(resp.Body)

	wrapper := struct {
		Results []struct {
			Task struct {
				URL string `json:"url"`
			} `json:"task"`

			Page struct {
				URL string `json:"url"`
			} `json:"page"`
		} `json:"results"`
	}{}

	err = dec.Decode(&wrapper)

	if err != nil {
		return []string{}, err
	}

	for _, r := range wrapper.Results {
		u, err := url.Parse(r.Task.URL)
		if err != nil {
			continue
		}
		output = append(output, u.Hostname())
	}

	for _, r := range wrapper.Results {
		u, err := url.Parse(r.Page.URL)
		if err != nil {
			continue
		}
		output = append(output, u.Hostname())
	}
	return output, nil
}

func (f fetch) bufferOverrun() ([]string, error) {
	ret := make([]string, 0)

	wrapper := struct {
		Records []string `json:"FDNS_A"`
	}{}
	err := jsonGET(
		fmt.Sprintf("https://dns.bufferover.run/dns?q=.%s", f.domain),
		&wrapper,
	)

	if err != nil {
		return ret, err
	}

	for _, r := range wrapper.Records {
		parts := strings.SplitN(r, ",", 2)
		if len(parts) != 2 {
			continue
		}
		ret = append(ret, parts[1])
	}

	return ret, nil
}

func (f fetch) certSpotter() ([]string, error) {
	ret := make([]string, 0)

	wrapper := []struct {
		ID           string   `json:"id"`
		TbsSha256    string   `json:"tbs_sha256"`
		DNSNames     []string `json:"dns_names"`
		PubkeySha256 string   `json:"pubkey_sha256"`
		NotBefore    string   `json:"not_before"`
		NotAfter     string   `json:"not_after"`
	}{}

	err := jsonGET(
		fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names", f.domain),
		&wrapper,
	)

	if err != nil {
		return ret, err
	}

	for _, r := range wrapper {
		ret = append(ret, r.DNSNames...)
	}

	return ret, nil
}
