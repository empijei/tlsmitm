package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
)

const encr = "encrypted"
const unen = "unencrypted"

type listener struct {
	localport, remoteport, remoteip string
	secure                          bool
	certconf                        *tls.Config
	protoSwitch                     bool
}

func (l *listener) String() string {
	var localtext, remotetext string
	if l.secure {
		localtext = encr
		if l.protoSwitch {
			remotetext = unen
		} else {
			remotetext = encr
		}
	} else {
		localtext = unen
		if l.protoSwitch {
			remotetext = encr
		} else {
			remotetext = unen
		}
	}
	return localtext + " " + l.localport + " -> " + l.remoteip + l.remoteport + " " + remotetext
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
			//This means (l.secure && !l.protoSwitch) || (!l.secure && l.protoSwitch)
			if l.secure != l.protoSwitch {
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
