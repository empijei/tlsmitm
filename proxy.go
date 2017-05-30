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

func (l *listener) Listen() {
	//TODO check if the listener was properly constructed

	var ll net.Listener
	var err error
	if l.secure {
		//TODO setup conf
		ll, err = tls.Listen("tcp", l.localport, nil)
	} else {
		ll, err = net.Listen("tcp", l.localport)
	}
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ll.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(lc net.Conn) {
			defer func() { _ = lc.Close() }()
			//Dump traffic to stdout
			lt := io.TeeReader(lc, os.Stdout)
			var rc net.Conn
			if l.secure {
				rc, err = tls.Dial("tcp", l.remoteip+l.remoteport, nil)
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
