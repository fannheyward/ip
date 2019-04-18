package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

func parse(input string) ([]string, error) {
	ips := make([]string, 0)

	ip := net.ParseIP(input)
	if ip != nil {
		ips = append(ips, ip.String())
		return ips, nil
	}

	nips, err := net.LookupIP(input)
	if err != nil {
		return nil, err
	}

	for _, v := range nips {
		ips = append(ips, v.String())
	}

	return ips, nil
}

func fetchIP(ip string) {
	local := false
	_url := "http://freeapi.ipip.net/" + ip
	if ip == "" {
		_url = "http://ip.cn"
		local = true
	}
	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		log.Fatalln("new req error:", err.Error())
	}

	req.Header.Add("User-Agent", "curl/7.54.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("HTTP error:", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("read body error:", err.Error())
	}

	if local {
		log.Println(string(body))
		os.Exit(0)
	}

	m := make([]string, 0)
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatalln("json.Unmarshal error:", err.Error())
	}

	s := ip
	for _, e := range m {
		if e == "" {
			continue
		}
		if s == ip {
			s = s + ": " + e
		} else {
			s = s + "-" + e
		}
	}
	log.Println(s)
}

func main() {
	log.SetFlags(0)
	args := os.Args
	if len(args) == 1 {
		fetchIP("")
		os.Exit(0)
	}

	if len(args) > 2 {
		log.Fatalln("too many args")
	}

	ips, err := parse(args[1])
	if err != nil {
		log.Fatalln("no IP found.")
	}

	wg := sync.WaitGroup{}
	for _, v := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			fetchIP(ip)
		}(v)
	}

	wg.Wait()
}
