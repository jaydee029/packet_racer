package packet

import (
	"crypto/rand"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type PacketConfig struct {
	SrcIP, DstIP     net.IP
	SrcPort, DstPort layers.UDPPort
	SrcMAC, DstMAC   net.HardwareAddr
	PayloadSize      int
}
type PacketOption func(*PacketConfig) error

func NewPacketConfig(opts ...PacketOption) (*PacketConfig, error) {
	config := &PacketConfig{}

	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func getTestMsg() []byte {
	// Generate a 1500 byte random message
	buf := make([]byte, 1470)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}

func WithIpLayer(srcIp, dstIp net.IP) PacketOption {
	return func(c *PacketConfig) error {
		c.SrcIP = srcIp
		c.DstIP = dstIp
		return nil
	}
}

// WithUdpLayer enables the UDP layer in the packet.
func WithUdpLayer(srcPort, dstPort int) PacketOption {
	return func(c *PacketConfig) error {
		c.SrcPort = layers.UDPPort(srcPort)
		c.DstPort = layers.UDPPort(dstPort)
		return nil
	}
}

// WithPayloadSize sets the payload size for the packet.
func WithPayloadSize(size int) PacketOption {
	return func(c *PacketConfig) error {
		c.PayloadSize = size
		return nil
	}
}
func WithEthernetLayer(srcMAC, dstMAC net.HardwareAddr) PacketOption {
	return func(c *PacketConfig) error {
		c.SrcMAC = srcMAC
		c.DstMAC = dstMAC
		return nil
	}
}
func BuildPacket(c *PacketConfig) ([]byte, error) {
	buf := gopacket.NewSerializeBuffer()
	var layersToSerialize []gopacket.SerializableLayer

	// Automatically include the Ethernet layer if MAC addresses are provided
	if c.SrcMAC != nil && c.DstMAC != nil {
		ethLayer := &layers.Ethernet{
			SrcMAC:       c.SrcMAC,
			DstMAC:       c.DstMAC,
			EthernetType: layers.EthernetTypeIPv4,
		}
		layersToSerialize = append(layersToSerialize, ethLayer)
	}

	// Set IP layer
	ipLayer := &layers.IPv4{
		Version:  4,
		TTL:      64,
		SrcIP:    c.SrcIP,
		DstIP:    c.DstIP,
		Protocol: layers.IPProtocolUDP,
	}
	layersToSerialize = append(layersToSerialize, ipLayer)

	udpLayer := &layers.UDP{
		SrcPort: c.SrcPort,
		DstPort: c.DstPort,
	}
	udpLayer.SetNetworkLayerForChecksum(ipLayer) // Important for checksum calculation
	layersToSerialize = append(layersToSerialize, udpLayer)

	payload := getTestMsg() //make([]byte, c.PayloadSize)

	// Optionally, fill the payload with data
	layersToSerialize = append(layersToSerialize, gopacket.Payload(payload))

	// Serialize the packet layers into the buffer
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}, layersToSerialize...); err != nil {
		return nil, fmt.Errorf("error serializing packet: %w", err)
	}

	return buf.Bytes(), nil
}

// converts a uint16 from host- to network byte order.
func tonet(u uint16) uint16 {
	return (u<<8)&0xff00 | u>>8
}
