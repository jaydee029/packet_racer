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
	return (i * 1470) / (1024 * 1024)
}

func spaces(packets, length int) {
	for i := 0; i <= length-packets; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("|")
}

func stats(c *packet.Commands) {
	keys := make([]string, 0, 4)

	for k, _ := range c.Total {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i int, j int) bool {
		return c.Total[keys[i]] > c.Total[keys[j]]
	})

	fmt.Println("*--------------------------------------------------------------*")
	fmt.Println()
	fmt.Println("|UDP  METHOD |Packets Transfered|Mega Bytes Transfered|")
	fmt.Println()

	for _, k := range keys {
		p := fmt.Sprintf("%s", c.Total[k])
		m := fmt.Sprintf("%s", Mb(c.Total[k]))
		fmt.Printf("|%s|%v", k, c.Total[k])
		spaces(len(p), 26)
		fmt.Printf("%v", Mb(c.Total[k]))
		spaces(len(m), 29)
		fmt.Println()
	}
	fmt.Println()

	fmt.Printf("Port used for server:%v\nSize of each Packet(in bytes):%v\nDuration of each method:%v\n", c.Port, c.Size, c.Time)

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

func server(c *packet.Commands) {
	cmd = exec.Command("nc", "-ul", "8080")

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting Server:", err)
		return
	}

	fmt.Printf("Listening at %v.", c.Port)

	total_duration := (c.Time*4 + 6)
	duration := time.Duration(total_duration) * time.Second
	time.AfterFunc(duration, func() {
		stopNetcat()
	})

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for Server:", err)
		return
	}
}

func Cli() {
	inptime := flag.Int("duration", 1, "duration for each method")
	inpsize := flag.Int("size", 1470, "size of each packet in bytes")
	inpport := flag.Int("port", 8080, "port number")

	flag.Parse()

	c := &packet.Commands{
		Time:  *inptime,
		Size:  *inpsize,
		Port:  *inpport,
		Total: map[string]int{},
	}

	go server(c)

	time.Sleep(2 * time.Second)

	c.Total["net.dial    "] = packet.Udpserver(c)
	c.Total["AF inet     "] = packet.NewAFinet(c)
	c.Total["AF packet   "] = packet.Afpacket(c)
	c.Total["Raw socket  "] = packet.Rawsocket(c)

	stats(c)

}
