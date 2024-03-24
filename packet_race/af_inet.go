package packet

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewAFinet() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)

	if err != nil {
		fmt.Println("1", err)
		return
	}

	addr := &syscall.SockaddrInet4{
		Addr: [4]byte{127, 0, 0, 1},
		Port: 8080,
	}
	buf := getTestMsg()
	total := 0

	timerCh := time.After(1 * time.Second)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-timerCh:
			// If 8 seconds elapsed, print the final counter value and shutdown the server
			fmt.Println("Server shutting down due to timeout...")
			fmt.Printf("Total Packets Sent: %d\n", total)
			syscall.Close(fd)
			return

		case sig := <-signalCh:
			// If termination signal received, print the final counter value and shutdown the server
			fmt.Printf("Received signal %s. Shutting down...\n", sig)
			fmt.Printf("Total Packets sent: %d\n", total)
			syscall.Close(fd)
			return

		default:
			err := syscall.Sendto(fd, buf, 0, addr)
			if err != nil {
				fmt.Println("Error sending package", err)
				continue
			}
			total++
		}
	}
}
