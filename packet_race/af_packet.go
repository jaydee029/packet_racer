package packet

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Afpacket(c *Commands) int {

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(tonet(syscall.ETH_P_IP)))
	if err != nil {
		fmt.Errorf("failed to create socket: %w", err)
		return 0
	}
	defer syscall.Close(fd)

	ifi, err := net.InterfaceByName("eth0")
	if err != nil {
		fmt.Errorf("failed to get interface %s: %w", "eth0", err)
		return 0
	}

	dstaddr := &syscall.SockaddrLinklayer{
		Protocol: tonet(syscall.ETH_P_IP),
		Ifindex:  ifi.Index,
	}

	// Construct the packet once outside the loop
	// create a packet configuration
	config, err := NewPacketConfig(
		WithEthernetLayer(net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe}, net.HardwareAddr{0xc0, 0xff, 0xee, 0x00, 0x00, 0x00}),
		WithIpLayer(net.IP{127, 0, 0, 1}, net.IP{127, 0, 0, 1}),
		WithUdpLayer(8081, c.Port),
		//WithPayloadSize(1490),
	)
	if err != nil {
		fmt.Errorf("error configuring packet: %v", err)
		return 0
	}
	// build the packet
	packet, err := BuildPacket(config, c.Size)
	if err != nil {
		fmt.Errorf("failed to build packet: %w", err)
		return 0
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
