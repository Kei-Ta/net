package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/kei-ta/net/internal/parse"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type App struct {
	Cli        *cli.App
	Connection *raw.Conn
}

func NewApp(version string) *App {
	app := App{}

	app.Cli = &cli.App{
		Name:  "net",
		Usage: "The cli tool to handle network",
		Commands: []*cli.Command{
			{
				Name:    "capture",
				Aliases: []string{"c"},
				Usage:   "capture paclet",
				Action: func(cCtx *cli.Context) error {
					app.Capture()
					return nil
				},
			}, {
				Name:    "arp",
				Aliases: []string{"a"},
				Usage:   "arp command",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: arp")
					return nil
				},
			},
			{
				Name:    "ping",
				Aliases: []string{"p"},
				Usage:   "ping command",
				Action: func(cCtx *cli.Context) error {
					app.Ping()
					return nil
				},
			},
		},
	}
	return &app
}

func (a *App) Run(ctx context.Context) error {
	return a.Cli.RunContext(ctx, os.Args)
}

func (a *App) Capture() {
	iface, err := net.InterfaceByName("en0")
	if err != nil {
		log.Fatalf("Failed to get network interface: %v", err)
	}

	conn, err := raw.ListenPacket(iface, uint16(ethernet.EtherTypeIPv4), nil)
	if err != nil {
		log.Fatalf("Failed to create raw socket: %v", err)
	}

	defer conn.Close()
	for {
		conn.SetReadDeadline(time.Now().Add(100 * time.Second))
		reply := make([]byte, 1500)
		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			fmt.Printf("Failed to read packet: %v\n", err)
			os.Exit(1)
		}
		// 受信したデータの送信元IPアドレスとバイト数を出力
		fmt.Printf("Received %d bytes from %v\n", n, peer)

		ethFrame, err := parse.ParseEthernetFrame(reply[:n], true)
		if err != nil {
			fmt.Printf("Failed to parse Ethernet frame: %v\n", err)
		}

		ipFrame, err := parse.ParseIpFrame(ethFrame.Payload, true)
		if err != nil {
			fmt.Printf("Failed to parse Ip frame: %v\n", err)
		}
		if ipFrame.Protocol == 6 {
			_, err := parse.ParseUdpFrame(ipFrame.Payload, true)
			if err != nil {
				fmt.Printf("Failed to parse Udp frame: %v\n", err)
			}
		} else if ipFrame.Protocol == 17 {
			_, err := parse.ParseTcpFrame(ipFrame.Payload, true)
			if err != nil {
				fmt.Printf("Failed to parse Tcp frame: %v\n", err)
			}
		} else {
			fmt.Println("unlnown frame")
		}
	}
}

func (a *App) Arp() {

}

func (a *App) Ping() {
	iface, err := net.InterfaceByName("en0")
	if err != nil {
		log.Fatalf("Failed to get network interface: %v", err)
	}

	conn, err := raw.ListenPacket(iface, uint16(ethernet.EtherTypeIPv4), nil)
	if err != nil {
		log.Fatalf("Failed to create raw socket: %v", err)
	}

	defer conn.Close()

	icmpMsg := icmp.Message{
		Type:     ipv4.ICMPTypeEcho,
		Code:     0,
		Checksum: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Hello"),
		},
	}

	icmpBytes, err := icmpMsg.Marshal(nil)
	if err != nil {
		log.Fatalf("err")
	}
	ipHdr := &ipv4.Header{
		Version:  4,
		Len:      ipv4.HeaderLen,
		TTL:      64,
		Protocol: 1,
		Src:      net.ParseIP("192.168.3.4").To4(),
		Dst:      net.ParseIP("192.168.3.4").To4(),
	}
	ipHdrBytes, err := ipHdr.Marshal()
	if err != nil {
		log.Fatalf("Failed to marshal IP header: %v", err)
	}

	ethFrame := &ethernet.Frame{
		Destination: net.HardwareAddr{0x6c, 0x7e, 0x67, 0xcb, 0x97, 0xaa},
		Source:      iface.HardwareAddr,
		EtherType:   ethernet.EtherTypeIPv4,
		Payload:     append(ipHdrBytes, icmpBytes...),
	}
	ethFrameBytes, err := ethFrame.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to marshal Ethernet frame: %v", err)
	}
	if _, err := conn.WriteTo(ethFrameBytes, &raw.Addr{HardwareAddr: ethFrame.Destination}); err != nil {
		log.Fatalf("Failed to send Ethernet frame: %v", err)
	}

	log.Println("Ethernet frame sent successfully")

	conn.SetReadDeadline(time.Now().Add(100 * time.Second))
	// 応答パケットを受信するためのバッファ（最大1500バイト）
	reply := make([]byte, 1500)

	for {
		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			fmt.Printf("Failed to read ICMP reply: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Received %d bytes from %v\n", n, peer)

		ethFrameReply, err := parse.ParseEthernetFrame(reply[:n], false)
		if err != nil {
			fmt.Printf("Failed to parse Ethernet frame: %v\n", err)
		}

		ipFrameReply, err := parse.ParseIpFrame(ethFrameReply.Payload, true)
		if err != nil {
			fmt.Printf("Failed to parse Ip frame: %v\n", err)
		}

		if ipFrameReply.Protocol == 1 {
			fmt.Printf("Protocol: %d\n", ipFrameReply.Protocol)
			_, err := parse.ParseIcmpFrame(ipFrameReply.Payload, true)
			if err != nil {
				fmt.Printf("Failed to parse Udp frame: %v\n", err)
			}
			return
		}

	}
}
