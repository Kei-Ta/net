package parse

import (
	"encoding/binary"
	"fmt"
	"net"
)

type EthernetFrame struct {
	Destination net.HardwareAddr
	Source      net.HardwareAddr
	EtherType   uint16
	Payload     []byte
}

type IpFrame struct {
	Version   uint8
	HeaderLen uint8
	ToS       uint8
	PacketLen uint16
	TTL       uint8
	Protocol  uint8
	SrcIp     net.IP
	DesIp     net.IP
	Payload   []byte
}

type IcmpFrame struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	Payload  []byte
}

type UdpFrame struct {
	SrcPort  uint16
	DesPort  uint16
	Datagram uint16
	Checksum uint16
	Payload  []byte
}

type TcpFrame struct {
	SrcPort    uint16
	DesPort    uint16
	Sequence   uint32
	AckNumner  uint32
	ControlBit uint8
	Payload    []byte
}

// ParseEthernetFrameはEthernetフレームを解析する関数

// Parameters:
// - data: バイナリデータ
// - output: Ethernetフレームを出力するかのフラグ
//
// Returns:
// - *EthernetFrame : EthernetFrame構造体のポインタ
// - error : エラー
func ParseEthernetFrame(data []byte, output bool) (*EthernetFrame, error) {
	if len(data) < 14 {
		return nil, fmt.Errorf("data too short for Ethernet frame")
	}

	destMAC := data[0:6]
	srcMAC := data[6:12]
	etherType := binary.BigEndian.Uint16(data[12:14])
	payload := data[14:] // ペイロードはヘッダの後に続く

	ethFrame := &EthernetFrame{
		Destination: destMAC,
		Source:      srcMAC,
		EtherType:   etherType,
		Payload:     payload,
	}

	if output {
		fmt.Println("------------------Ethernet Frame Start---------------------")
		fmt.Printf("Received Ethernet Frame:\n")
		fmt.Printf("  Destination MAC: %s\n", ethFrame.Destination)
		fmt.Printf("  Source MAC: %s\n", ethFrame.Source)
		fmt.Printf("  EtherType: %d\n", ethFrame.EtherType)
		fmt.Println("------------------Ethernet Frame End---------------------")
	}
	return ethFrame, nil
}

// IPフレームを解析する関数
// Parameters
// Returns
func ParseIpFrame(data []byte, output bool) (*IpFrame, error) {
	//Frameの大きさを検証する処理

	firstByte := data[0]
	firstFourBits := firstByte >> 4 // 最初の4ビットを取得
	ipFrame := &IpFrame{
		Version:   firstFourBits & 0x0F,
		HeaderLen: firstByte & 0x0F,
		ToS:       data[1],
		PacketLen: binary.BigEndian.Uint16(data[2:4]),
		TTL:       data[8],
		Protocol:  data[9],
		SrcIp:     net.IP(data[12:16]),
		DesIp:     net.IP(data[16:20]),
		Payload:   data[20:],
	}
	if output {
		fmt.Println("------------------IP Frame Start---------------------")
		fmt.Printf("Version: %04b\n", ipFrame.Version)
		fmt.Printf("HeaderLen: %04b\n", ipFrame.HeaderLen)
		fmt.Printf("ToS: %08b\n", ipFrame.ToS)
		fmt.Printf("PacketLen: %d\n", ipFrame.PacketLen)
		fmt.Printf("TTL: %d\n", ipFrame.TTL)
		fmt.Printf("Protocol: %d\n", ipFrame.Protocol)
		fmt.Printf("SrcIp: %d\n", ipFrame.SrcIp)
		fmt.Printf("DesIp: %d\n", ipFrame.DesIp)
		fmt.Println("------------------IP Frame End---------------------")
	}
	return ipFrame, nil
}

// ICMPフレームを解析する関数
func ParseIcmpFrame(payload []byte, output bool) (*IcmpFrame, error) {
	//Frameの大きさを検証する処理
	icmpFrame := &IcmpFrame{
		Type:     payload[0],
		Code:     payload[1],
		Checksum: binary.BigEndian.Uint16(payload[2:4]),
		Payload:  payload[4:],
	}
	if output {
		fmt.Println("------------------ICMP Frame Start---------------------")
		fmt.Printf("Type: %d\n", icmpFrame.Type)
		fmt.Printf("Code: %d\n", icmpFrame.Code)
		fmt.Println("------------------ICMP Frame End---------------------")
	}
	return icmpFrame, nil
}

// UDPフレームを解析する関数
func ParseUdpFrame(payload []byte, output bool) (*UdpFrame, error) {
	//Frameの大きさを検証する処理
	udpFrame := &UdpFrame{
		SrcPort:  binary.BigEndian.Uint16(payload[0:2]),
		DesPort:  binary.BigEndian.Uint16(payload[2:4]),
		Datagram: binary.BigEndian.Uint16(payload[4:6]),
		Checksum: binary.BigEndian.Uint16(payload[6:8]),
		Payload:  payload[8:],
	}
	if output {
		fmt.Println("------------------UDP Frame Start---------------------")
		fmt.Printf("SrcPort: %d\n", udpFrame.SrcPort)
		fmt.Printf("DesPort: %d\n", udpFrame.DesPort)
		fmt.Printf("Datagram: %d\n", udpFrame.Datagram)
		fmt.Printf("Checksum: %d\n", udpFrame.Checksum)
		fmt.Println("------------------UDP Frame End---------------------")
	}
	return udpFrame, nil
}

// TCPフレームを解析する関数
func ParseTcpFrame(payload []byte, output bool) (*TcpFrame, error) {
	//Frameの大きさを検証する処理
	tcpFrame := &TcpFrame{
		SrcPort:    binary.BigEndian.Uint16(payload[0:2]),
		DesPort:    binary.BigEndian.Uint16(payload[2:4]),
		Sequence:   binary.BigEndian.Uint32(payload[4:8]),
		AckNumner:  binary.BigEndian.Uint32(payload[8:12]),
		ControlBit: payload[33],
		Payload:    payload[8:],
	}
	if output {
		fmt.Println("------------------TCP Frame Start---------------------")
		fmt.Printf("SrcPort: %d\n", tcpFrame.SrcPort)
		fmt.Printf("DesPort: %d\n", tcpFrame.DesPort)
		fmt.Printf("Sequence: %d\n", tcpFrame.Sequence)
		fmt.Printf("AckNumner: %d\n", tcpFrame.AckNumner)
		fmt.Printf("ControlBit: %08b\n", tcpFrame.ControlBit)
		fmt.Println("------------------TCP Frame End---------------------")
	}
	return tcpFrame, nil
}
