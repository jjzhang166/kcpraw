package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/urfave/cli"
	"github.com/xtaci/kcp-go"
)

var VERSION = "SELFBUILD"

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

func checkError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	myApp := cli.NewApp()
	myApp.Name = "kcptun"
	myApp.Usage = "kcptun client"
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "localaddr,l",
			Value: ":12948",
			Usage: "local listen address",
		},
		cli.StringFlag{
			Name:  "remoteaddr, r",
			Value: "vps:29900",
			Usage: "kcp server address",
		},
		cli.StringFlag{
			Name:   "key",
			Value:  "it's a secrect",
			Usage:  "key for communcation, must be the same as kcptun server",
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
			Value: 128,
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
			Usage: "set FEC group size, must be the same as server",
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
		addr, err := net.ResolveTCPAddr("tcp", c.String("localaddr"))
		checkError(err)
		listener, err := net.ListenTCP("tcp", addr)
		checkError(err)
		log.Println("listening on:", listener.Addr())

	START_KCP:
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
		log.Println("communication mode:", c.String("mode"))
		// kcp server
		kcpconn, err := kcp.DialWithOptions(mode, c.Int("fec"), c.String("remoteaddr"), []byte(c.String("key")))
		checkError(err)
		log.Println("remote address:", c.String("remoteaddr"))
		log.Println("sndwnd:", c.Int("sndwnd"), "rcvwnd:", c.Int("rcvwnd"))
		log.Println("mtu:", c.Int("mtu"))
		log.Println("fec:", c.Int("fec"))
		log.Println("acknodelay:", c.Bool("acknodelay"))
		log.Println("dscp:", c.Int("dscp"))

		kcpconn.SetWindowSize(c.Int("sndwnd"), c.Int("rcvwnd"))
		kcpconn.SetMtu(c.Int("mtu"))
		kcpconn.SetACKNoDelay(c.Bool("acknodelay"))
		kcpconn.SetDSCP(c.Int("dscp"))

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
		session, err := yamux.Client(kcpconn, config)
		checkError(err)
		mux = session

		for {
			p1, err := listener.AcceptTCP()
			if err != nil {
				log.Println(err)
				continue
			}
			p2, err := mux.Open()
			if err != nil { // yamux failure
				log.Println(err)
				kcpconn.Close()
				p1.Close()
				goto START_KCP
			}
			go handleClient(p1, p2)
		}
	}
	myApp.Run(os.Args)
}
