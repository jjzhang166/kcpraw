package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/urfave/cli"
	"github.com/xtaci/kcp-go"
)

var VERSION = "SELFBUILD"

// handle multiplex-ed connection
func handleMux(conn *kcp.UDPSession, target string) {
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
	go sig_handler()
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
			Usage: "mode for communication: fast3, fast2, fast, normal",
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
			Value: 0,
			Usage: "set DSCP(6bit)",
		},
	}
	myApp.Action = func(c *cli.Context) {
		log.Println("version:", VERSION)
		nodelay, interval, resend, nc := 0, 40, 0, 0
		switch c.String("mode") {
		case "normal":
			nodelay, interval, resend, nc = 0, 30, 2, 1
		case "fast":
			nodelay, interval, resend, nc = 0, 20, 2, 1
		case "fast2":
			nodelay, interval, resend, nc = 1, 20, 2, 1
		case "fast3":
			nodelay, interval, resend, nc = 1, 10, 2, 1
		}

		lis, err := kcp.ListenWithOptions(c.Int("fec"), c.String("listen"), []byte(c.String("key")))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("listening on ", lis.Addr())
		log.Println("nodelay parameters:", nodelay, interval, resend, nc)
		log.Println("sndwnd:", c.Int("sndwnd"), "rcvwnd:", c.Int("rcvwnd"))
		log.Println("mtu:", c.Int("mtu"))
		log.Println("fec:", c.Int("fec"))
		log.Println("acknodelay:", c.Bool("acknodelay"))
		log.Println("dscp:", c.Int("dscp"))
		for {
			if conn, err := lis.Accept(); err == nil {
				log.Println("remote address:", conn.RemoteAddr())
				conn.SetNoDelay(nodelay, interval, resend, nc)
				conn.SetMtu(c.Int("mtu"))
				conn.SetWindowSize(c.Int("sndwnd"), c.Int("rcvwnd"))
				conn.SetACKNoDelay(c.Bool("acknodelay"))
				conn.SetDSCP(c.Int("dscp"))
				go handleMux(conn, c.String("target"))
			} else {
				log.Println(err)
			}
		}
	}
	myApp.Run(os.Args)
}
func sig_handler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1)

	for {
		switch <-ch {
		case syscall.SIGUSR1:
			log.Printf("KCP SNMP:%+v", kcp.DefaultSnmp.Copy())
		}
	}
}
