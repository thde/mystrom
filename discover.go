package mystrom

import (
	"context"
	"fmt"
	"net"
)

type DeviceType byte

// https://api.mystrom.ch -> Discovery devices -> Device types
const (
	DeviceTypeBulb                    DeviceType = 102
	DeviceTypeButtonPlus1stGeneration DeviceType = 103
	DeviceTypeButtonSmall             DeviceType = 104
	DeviceTypeLEDStrip                DeviceType = 105
	DeviceTypeSwitchCH                DeviceType = 106
	DeviceTypeSwitchEU                DeviceType = 107
	DeviceTypeMotionSensor            DeviceType = 110
	DeviceTypeGateway                 DeviceType = 112
	DeviceTypeSTECCO                  DeviceType = 113
	DeviceTypeButtonPlus2ndGeneration DeviceType = 118
	DeviceTypeSwitchZero              DeviceType = 120
)

func (d DeviceType) String() string {
	switch d {
	case DeviceTypeBulb:
		return "Bulb"
	case DeviceTypeButtonPlus1stGeneration:
		return "Button plus 1st generation"
	case DeviceTypeButtonSmall:
		return "Button small/simple"
	case DeviceTypeLEDStrip:
		return "LED Strip"
	case DeviceTypeSwitchCH:
		return "Switch CH"
	case DeviceTypeSwitchEU:
		return "Switch EU"
	case DeviceTypeMotionSensor:
		return "Motion Sensor"
	case DeviceTypeGateway:
		return "Gateway"
	case DeviceTypeSTECCO:
		return "STECCO/CUBO"
	case DeviceTypeButtonPlus2ndGeneration:
		return "Button Plus 2nd generation"
	case DeviceTypeSwitchZero:
		return "Switch Zero"
	default:
		return fmt.Sprintf("%d", d)
	}
}

type Device struct {
	Address net.Addr
	MAC     net.HardwareAddr
	Type    DeviceType

	Cloud      bool
	Registered bool
	MeshChild  bool
}

type Discover struct {
	// See func net.Dial for a description of the Address parameter.
	Address string
	// Network must be "udp", "udp4", "udp6", "unixgram", or an IP transport.
	// See net.ListenPacket for more details on the Network parameter.
	Network      string
	ListenConfig net.ListenConfig
}

// Device blocks until a MyStrom device has been discovered.
// Each device cyclically (every 5 seconds) sends a broadcast packet
// using the UDP protocol to the address 255.255.255.255 and port 7979.
func (d *Discover) Device(ctx context.Context) (Device, error) {
	address := defaultString(d.Address, ":7979")
	network := defaultString(d.Network, "udp")

	pc, err := d.ListenConfig.ListenPacket(ctx, network, address)
	if err != nil {
		return Device{}, fmt.Errorf("listen packet error for address %s, network %s: %w", address, network, err)
	}
	defer pc.Close()

	buf := make([]byte, 8)
	_, addr, err := pc.ReadFrom(buf)
	if err != nil {
		return Device{}, fmt.Errorf("error reading packet: %w", err)
	}

	return parseDevicePayload(buf, addr)
}

func defaultString(s, def string) string {
	if s != "" {
		return s
	}

	return def
}

// https://api.mystrom.ch -> Discovery devices
func parseDevicePayload(buf []byte, addr net.Addr) (d Device, err error) {
	if len(buf) < 8 { // broadcast packet is 8 bytes
		return d, fmt.Errorf("payload too small: %d", len(buf))
	}

	d.Address = addr
	d.MAC = net.HardwareAddr(buf[:6])
	d.Type = DeviceType(buf[6])

	status := buf[7]
	d.MeshChild = status&(1<<0) != 0
	d.Registered = status&(1<<1) != 0
	d.Cloud = status&(1<<2) != 0

	return d, nil
}
