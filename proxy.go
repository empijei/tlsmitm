package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
)

type listener struct {
	localport, remoteport, remoteip string
	secure                          bool
	certconf                        *tls.Config
}

func (l *listener) String() string {
	var text = "encrypted"
	//This is ugly but fun
	if !l.secure {
		text = "un" + text
	}
	return l.localport + " -> " + l.remoteip + l.remoteport + " " + text
}

func (l *listener) Listen() {
	//TODO check if the listener was properly constructed

	var ll net.Listener
	var err error
	if l.secure {
		ll, err = tls.Listen("tcp", l.localport, l.certconf)
	} else {
		ll, err = net.Listen("tcp", l.localport)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Running rule " + l.String())
	for {
		conn, err := ll.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(lc net.Conn) {
			log.Printf("Got incoming connection from ip %s on port %s", conn.RemoteAddr().String(), l.localport)
			defer func() { _ = lc.Close() }()
			//Dump traffic to stdout
			lt := io.TeeReader(lc, os.Stdout)
			var rc net.Conn
			if l.secure {
				rc, err = tls.Dial("tcp", l.remoteip+l.remoteport,
					&tls.Config{
						//TODO make this configurable
						InsecureSkipVerify: true,
					})
			} else {
				rc, err = net.Dial("tcp", l.remoteip+l.remoteport)
			}
			if err != nil {
				//TODO handle this
				panic(err)
			}
			defer func() { _ = rc.Close() }()
			//Dump traffic to stdout
			rt := io.TeeReader(rc, os.Stdout)
			go func() { _, _ = io.Copy(rc, lt) }()
			_, _ = io.Copy(lc, rt)
		}(conn)
	}

}
