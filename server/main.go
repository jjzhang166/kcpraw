package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/yamux"
	"github.com/xtaci/kcp-go"
)

var VERSION = "SELFBUILD"

// handle multiplex-ed connection
func handleMux(conn *kcp.UDPSession, key, target string, mtu, sndwnd, rcvwnd int, acknodelay bool, dscp int) {
	conn.SetWindowSize(1024, 1024)
	conn.SetMtu(mtu)
	conn.SetWindowSize(sndwnd, rcvwnd)
	conn.SetACKNoDelay(acknodelay)
	conn.SetDSCP(dscp)

	// stream multiplex
	var mux *yamux.Session
	config := &yamux.Config{
		AcceptBacklog:          256,
		EnableKeepAlive:        true,
		KeepAliveInterval:      30 * time.Second,
		ConnectionWriteTimeout: 30 * time.Second,
		MaxStreamWindowSize:    1048576,
		LogOutput:              os.Stderr,
	}
	m, err := yamux.Server(conn, config)
	if err != nil {
		log.Println(err)
		return
	}
	mux = m
	defer mux.Close()

	for {
		p1, err := mux.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		p2, err := net.Dial("tcp", target)
		if err != nil {
			log.Println(err)
			return
		}
		go handleClient(p1, p2)
	}
}

func handleClient(p1, p2 net.Conn) {
	log.Println("stream opened")
	defer log.Println("stream closed")
	defer p1.Close()
	defer p2.Close()

	// start tunnel
	p1die := make(chan struct{})
	go func() {
		io.Copy(p1, p2)
		close(p1die)
	}()

	p2die := make(chan struct{})
	go func() {
		io.Copy(p2, p1)
		close(p2die)
	}()

	// wait for tunnel termination
	select {
	case <-p1die:
	case <-p2die:
	}
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	myApp := cli.NewApp()
	myApp.Name = "kcptun"
	myApp.Usage = "kcptun server"
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen,l",
			Value: ":29900",
			Usage: "kcp server listen address",
		},
		cli.StringFlag{
			Name:  "target, t",
			Value: "127.0.0.1:12948",
			Usage: "target server address",
		},
		cli.StringFlag{
			Name:   "key",
			Value:  "it's a secrect",
			Usage:  "key for communcation, must be the same as kcptun client",
			EnvVar: "KCPTUN_KEY",
		},
		cli.StringFlag{
			Name:  "mode",
			Value: "fast",
			Usage: "mode for communication: fast2, fast, normal, default",
		},
		cli.IntFlag{
			Name:  "mtu",
			Value: 1350,
			Usage: "set MTU of UDP packets, suggest 'tracepath' to discover path mtu",
		},
		cli.IntFlag{
			Name:  "sndwnd",
			Value: 1024,
			Usage: "set send window size(num of packets)",
		},
		cli.IntFlag{
			Name:  "rcvwnd",
			Value: 1024,
			Usage: "set receive window size(num of packets)",
		},
		cli.IntFlag{
			Name:  "fec",
			Value: 4,
			Usage: "set FEC group size, must be the same as client",
		},
		cli.BoolFlag{
			Name:  "acknodelay",
			Usage: "flush ack immediately when a packet is received",
		},
		cli.IntFlag{
			Name:  "dscp",
			Value: 46,
			Usage: "set DSCP(6bit)",
		},
	}
	myApp.Action = func(c *cli.Context) {
		log.Println("version:", VERSION)
		// KCP listen
		var mode kcp.Mode
		switch c.String("mode") {
		case "normal":
			mode = kcp.MODE_NORMAL
		case "default":
			mode = kcp.MODE_DEFAULT
		case "fast":
			mode = kcp.MODE_FAST
		case "fast2":
			mode = kcp.MODE_FAST2
		default:
			log.Println("unrecognized mode:", c.String("mode"))
			return
		}

		lis, err := kcp.ListenEncrypted(mode, c.Int("fec"), c.String("listen"), []byte(c.String("key")))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("listening on ", lis.Addr())
		log.Println("communication mode:", c.String("mode"))
		log.Println("sndwnd:", c.Int("sndwnd"), "rcvwnd:", c.Int("rcvwnd"))
		log.Println("mtu:", c.Int("mtu"))
		log.Println("fec:", c.Int("fec"))
		log.Println("acknodelay:", c.Bool("acknodelay"))
		log.Println("dscp:", c.Int("dscp"))
		for {
			if conn, err := lis.Accept(); err == nil {
				log.Println("remote address:", conn.RemoteAddr())
				go handleMux(conn,
					c.String("key"),
					c.String("target"),
					c.Int("mtu"),
					c.Int("sndwnd"), c.Int("rcvwnd"),
					c.Bool("acknodelay"),
					c.Int("dscp"),
				)
			} else {
				log.Println(err)
			}
		}
	}
	myApp.Run(os.Args)
}
