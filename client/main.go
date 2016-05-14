package main

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
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

type secureConn struct {
	encoder cipher.Stream
	decoder cipher.Stream
	conn    net.Conn
}

func newSecureConn(key string, conn net.Conn, iv []byte) *secureConn {
	sc := new(secureConn)
	sc.conn = conn
	commkey := sha256.Sum256([]byte(key))

	// encoder
	block, err := aes.NewCipher(commkey[:])
	if err != nil {
		log.Println(err)
		return nil
	}
	sc.encoder = cipher.NewCFBEncrypter(block, iv[:aes.BlockSize])

	// decoder
	block, err = aes.NewCipher(commkey[:])
	if err != nil {
		log.Println(err)
		return nil
	}
	sc.decoder = cipher.NewCFBDecrypter(block, iv[aes.BlockSize:])
	return sc
}

func (sc *secureConn) Read(p []byte) (n int, err error) {
	n, err = sc.conn.Read(p)
	if err == nil {
		sc.decoder.XORKeyStream(p[:n], p[:n])
	}
	return
}

func (sc *secureConn) Write(p []byte) (n int, err error) {
	sc.encoder.XORKeyStream(p, p)
	return sc.conn.Write(p)
}

func (sc *secureConn) Close() (err error) {
	return sc.conn.Close()
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
			Usage: "local listen addr:",
		},
		cli.StringFlag{
			Name:  "remoteaddr, r",
			Value: "vps:29900",
			Usage: "kcp server addr",
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
			Usage: "mode for communication: fast, normal, default",
		},
		cli.BoolFlag{
			Name:  "tuncrypt",
			Usage: "enable tunnel encryption, adds extra secrecy for data transfer",
		},
		cli.IntFlag{
			Name:  "mtu",
			Value: 1400,
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
		default:
			log.Println("unrecognized mode:", c.String("mode"))
			return
		}
		log.Println("communication mode:", c.String("mode"))
		// kcp server
		kcpconn, err := kcp.DialEncrypted(mode, c.String("remoteaddr"), []byte(c.String("key")))
		checkError(err)
		kcpconn.SetRetries(50)
		log.Println("remote address:", c.String("remoteaddr"))
		kcpconn.SetWindowSize(c.Int("sndwnd"), c.Int("rcvwnd"))
		log.Println("sndwnd:", c.Int("sndwnd"), "rcvwnd:", c.Int("rcvwnd"))
		kcpconn.SetMtu(c.Int("mtu"))
		log.Println("mtu:", c.Int("mtu"))

		// generate & send iv
		iv := make([]byte, 2*aes.BlockSize)
		io.ReadFull(crand.Reader, iv)
		_, err = kcpconn.Write(iv)
		checkError(err)

		// stream multiplex
		var mux *yamux.Session
		if c.Bool("tuncrypt") {
			scon := newSecureConn(c.String("key"), kcpconn, iv)
			session, err := yamux.Client(scon, nil)
			checkError(err)
			mux = session
		} else {
			session, err := yamux.Client(kcpconn, nil)
			checkError(err)
			mux = session
		}
		log.Println("tunnel encryption:", c.Bool("tuncrypt"))

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
