package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sync"
)

func parse(input string) ([]string, error) {
	ips := make([]string, 0)
	r, err := regexp.Compile("^([0-9]{1,3}\\.){3}[0-9]{1,3}$")
	if err != nil {
		return nil, err
	}

	if r.MatchString(input) {
		ips = append(ips, r.FindString(input))
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

func fetchIPCN(ip string) {
	_url := "http://ip.cn/" + ip
	cmd := exec.Command("curl", "-s", "-L", _url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	fmt.Print(string(out))
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fetchIPCN("")
		os.Exit(0)
	}

	if len(args) > 2 {
		fmt.Print("no or one arg only")
		os.Exit(1)
	}

	ips, err := parse(args[1])
	if err != nil {
		fmt.Print("no IP found.")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	for _, v := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			fetchIPCN(ip)
		}(v)
	}

	wg.Wait()
}
