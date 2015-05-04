// The MIT License (MIT)

// Copyright (c) 2015 Chris Olstrom

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"flag"
	"github.com/fatih/color"
	"net"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"syscall"
	"time"
)

func sameLength(a, b []string) bool {
	return len(a) == len(b)
}

func sameContents(a, b []string) bool {
	return reflect.DeepEqual(a, b)
}

func equivalent(a, b []string) bool {
	return sameLength(a, b) && sameContents(a, b)
}

func stringify(addresses []net.IP) []string {
	var transformed []string
	for _, ip := range addresses {
		transformed = append(transformed, net.IP.String(ip))
	}
	return transformed
}

func lookup(hostname string) []string {
	addresses, _ := net.LookupIP(hostname)
	sorted := stringify(addresses)
	sort.Strings(sorted)
	return sorted
}

func react(hostname, execute string, changes <-chan string) {
	changed := <-changes
	output, error := exec.Command(execute).Output()
	if error != nil {
		color.Red("[ERROR] %s", error)
	} else {
		color.Green("[CHANGE] %s %s", changed, output)
	}
}

func monitor(hostname string, interval int, changes chan<- string) {
	var knownAddresses []string
	for tick := range time.NewTicker(time.Duration(interval) * time.Second).C {
		addresses := lookup(hostname)
		color.White("[LOOKUP] %s %s %s", hostname, addresses, tick)
		if !equivalent(knownAddresses, addresses) {
			changes <- hostname
			knownAddresses = addresses
		}
	}
}

func noHostsProvided() bool {
	return flag.NArg() == 0
}

func stayAlive() {
	signals := make(chan os.Signal, 1)
	for {
		signal := <-signals
		if signal == syscall.SIGTERM {
			return
		}
	}
}

func monitorHosts(interval int, execute string) {
	color.Yellow("Monitoring %d hosts every %d seconds. Will invoke '%s' if any change.", flag.NArg(), interval, execute)
	changes := make(chan string)
	for _, hostname := range flag.Args() {
		go monitor(hostname, interval, changes)
		go react(hostname, execute, changes)
	}
}

func main() {
	interval := flag.Int("interval", 5, "interval for DNS queries")
	execute := flag.String("execute", "", "command to execute when a change is detected")
	flag.Parse()
	if noHostsProvided() {
		color.Red("[ERROR] You must provide one or more hostnames to monitor.")
		return
	}
	monitorHosts(*interval, *execute)
	stayAlive()
}
