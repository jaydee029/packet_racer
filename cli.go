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

func Mb(i int) int {
	return (i * 1485) / (1024 * 1024)
}

func stats(c *Commands) {
	keys := make([]string, 0, 4)

	for k, _ := range c.total {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i int, j int) bool {
		return c.total[keys[i]] > c.total[keys[j]]
	})

	fmt.Println("*--------------------------------------------------------------*")
	fmt.Println("|UDP  METHOD |Packets Transfered|Mega Bytes Transfered|")

	for _, k := range keys {
		fmt.Printf("|%s|%v|%v|\n", k, c.total[k], Mb(c.total[k]))
	}

	fmt.Println("*--------------------------------------------------------------*")

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

	c.total["net.dial  "] = packet.Udpserver()
	c.total["AF inet   "] = packet.NewAFinet()
	c.total["AF packet "] = packet.Afpacket()
	c.total["Raw socket"] = packet.Rawsocket()

	stats(c)

}
