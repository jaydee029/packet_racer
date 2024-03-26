package packet

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Rawsocket(c *Commands) int {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Fatalf("Failed to create raw socket: %v", err)
		return 0
	}
	defer syscall.Close(fd)

	// Set options: here, we enable IP_HDRINCL to manually include the IP header
	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		log.Fatalf("Failed to set IP_HDRINCL: %v", err)
		return 0
	}

	config, err := NewPacketConfig(
		WithIpLayer(net.IP{127, 0, 0, 1}, net.IP{127, 0, 0, 1}),
		WithUdpLayer(int(8081), int(c.Port)),
		WithPayloadSize(1500),
	)
	if err != nil {
		fmt.Println("error configuring packet: %v", err)
		return 0
	}
	packet, err := BuildPacket(config, c.Size)
	if err != nil {
		fmt.Println("failed to build packet: %w", err)
		return 0
	}

	dstaddr := &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}

	total := 0

	timerCh := time.After(time.Duration(c.Time) * time.Second)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-timerCh:
			// If 1 second elapsed, print the final counter value and shutdown the server
			fmt.Println("Server shutting down due to timeout...")
			fmt.Printf("Total Packets Sent: %d\n", total)
			syscall.Close(fd)
			return total

		case sig := <-signalCh:
			// If termination signal received, print the final counter value and shutdown the server
			fmt.Printf("Received signal %s. Shutting down...\n", sig)
			fmt.Printf("Total Packets sent: %d\n", total)
			syscall.Close(fd)
			return total

		default:
			err := syscall.Sendto(fd, packet, 0, dstaddr)
			if err != nil {
				fmt.Println("Error sending package", err)
				continue
			}
			total++
		}
	}
}
