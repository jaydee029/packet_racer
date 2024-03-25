package main

import (
	"flag"
	"fmt"
	"os/exec"
	"sort"
	"syscall"
	"time"

	packet "github.com/jaydee029/packet_racer/packet_race"
)

type Commands struct {
	time  int
	size  int
	port  int
	total map[string]int
}

var cmd *exec.Cmd

type count struct {
	packets int
	Mbs     int
}

func Mb(i int) int {
	return (i * 1485) / (1024 * 1024)
}

func stats(c *Commands) map[string]count {
	keys := make([]string, 0, 4)

	for k, _ := range c.total {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i int, j int) bool {
		return c.total[keys[i]] > c.total[keys[j]]
	})

	for _, k := range keys {
		println(k)
	}

	out := make(map[string]count, 4)

	for _, k := range keys {
		out[k] = count{
			packets: c.total[k],
			Mbs:     Mb(c.total[k]),
		}

	}

	return out

}

func stopNetcat() {
	if cmd == nil || cmd.Process == nil {
		fmt.Println("Server is not running.")
		return
	}

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		fmt.Println("Error stopping Server:", err)
		return
	}

	fmt.Println("Server stopped.")
}

func (c *Commands) server() {
	cmd = exec.Command("nc", "-ul", "8080")

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting Server:", err)
		return
	}

	fmt.Printf("Listening at %v.", c.port)

	duration := 10 * time.Second
	time.AfterFunc(duration, func() {
		stopNetcat()
	})

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for Server:", err)
		return
	}
}

func (c *Commands) print(out map[string]count) {
	fmt.Println("*--------------------------------------------------------------*")
	fmt.Println("|UDP  METHOD |Packets Transfered|Mega Bytes Transfered|")
	for i, k := range out {
		fmt.Printf("|%s|%v|%v|\n", i, k.packets, k.Mbs)
	}
	fmt.Println("*--------------------------------------------------------------*")
}

func cli() {
	inptime := flag.Int("duration", 1, "duration for each method")
	inpsize := flag.Int("size", 1485, "size of each packet in bytes")
	inpport := flag.Int("port", 8080, "port number")

	flag.Parse()

	c := &Commands{
		time:  *inptime,
		size:  *inpsize,
		port:  *inpport,
		total: map[string]int{},
	}

	go c.server()

	time.Sleep(2 * time.Second)

	c.total["udp"] = packet.Udpserver()
	c.total["afinet"] = packet.NewAFinet()
	c.total["afpacket"] = packet.Afpacket()
	c.total["socket"] = packet.Rawsocket()

	output := stats(c)

	c.print(output)

}
