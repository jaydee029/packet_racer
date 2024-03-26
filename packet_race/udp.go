package packet

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Udpserver(c *Commands) int {
	addra := &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: c.Port,
	}

	fd, err := net.DialUDP("udp", nil, addra)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("the sever is live")

	total := 0
	buf := getMsg(c.Size)

	// Create a channel to listen for termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	timerCh := time.After(time.Duration(c.Time) * time.Second)

	for {
		select {
		case <-timerCh:
			// If 1 second elapsed, print the final counter value and shutdown the server
			fmt.Println("Server shutting down due to timeout...")
			fmt.Printf("Total Packets Sent: %d\n", total)
			fd.Close()
			return total
		case sig := <-signalCh:
			// If termination signal received, print the final counter value and shutdown the server
			fmt.Printf("Received signal %s. Shutting down...\n", sig)
			fmt.Printf("Total Packets Sent : %d\n", total)
			fd.Close()
			return total
		default:

			_, err := fd.Write(buf)
			if err != nil {
				fmt.Println("Error sending the package", err)
				continue
			}
			total++
		}
	}

}
