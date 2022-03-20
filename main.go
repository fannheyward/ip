package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type hostIP struct {
	domain string
	ip     string
}

func parse(input string, c chan hostIP) {
	if ip := net.ParseIP(input); ip != nil {
		c <- hostIP{domain: "", ip: ip.String()}
		return
	}

	url, err := url.Parse(input)
	if err != nil {
		log.Println("parse url error:", err.Error())
		return
	}

	if len(url.Hostname()) > 0 {
		input = url.Hostname()
	}

	nips, err := net.LookupIP(input)
	if err != nil {
		log.Println("unable to lookup:", input)
		return
	}

	for _, v := range nips {
		c <- hostIP{domain: input, ip: v.String()}
	}
}

func fetchIP(h hostIP) {
	_url := "https://ip.fm/" + h.ip
	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		log.Println("new req error:", err.Error())
		return
	}

	req.Header.Add("User-Agent", "curl/7.54.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("HTTP error:", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read body error:", err.Error())
		return
	}
	if len(body) == 0 {
		log.Println("response body none:", resp.Status)
		return
	}
	log.Println(h.domain + " " + string(body[:len(body)-1]))
}

func main() {
	log.SetFlags(0)
	args := os.Args
	if len(args) == 1 {
		fetchIP(hostIP{domain: "Local", ip: ""})
		os.Exit(0)
	}

	ch := make(chan hostIP, 10)
	go func() {
		for _, arg := range args[1:] {
			parse(arg, ch)
		}
		close(ch)
	}()

	wg := sync.WaitGroup{}
	for {
		v, ok := <-ch
		if !ok {
			break
		}

		wg.Add(1)
		go func(hi hostIP) {
			defer wg.Done()
			fetchIP(hi)
		}(v)
	}

	wg.Wait()
}
