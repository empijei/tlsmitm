package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
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
			//Dump client traffic to stdout
			lt := io.TeeReader(lc, os.Stdout)

			//Dump client traffic to file
			outclientfile, filerr := os.OpenFile(
				fmt.Sprintf("%d_client.log", time.Now().UnixNano()),
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0666)
			if filerr != nil {
				log.Printf("Could not open logfile: %s\n, proceeding without it", filerr.Error())
			} else {
				lt = io.TeeReader(lt, outclientfile)
			}

			var rc net.Conn
			//This means (l.secure && !l.protoSwitch) || (!l.secure && l.protoSwitch)
			log.Printf("Contacting remote server @%s", l.remoteip+l.remoteport)
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
				log.Println(err)
				return
			}
			log.Println("Connected to remote server")

			defer func() { _ = rc.Close() }()
			//Dump server traffic to stdout
			rt := io.TeeReader(rc, os.Stdout)

			//Dump server traffic to file
			outserverfile, filerr := os.OpenFile(
				fmt.Sprintf("%d_server.log", time.Now().UnixNano()),
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0666)
			if filerr != nil {
				log.Printf("Could not open logfile: %s\n, proceeding without it", filerr.Error())
			} else {
				rt = io.TeeReader(rt, outserverfile)
			}

			go func() {
				_, err := io.Copy(rc, lt)
				log.Printf("Error while sending to remote server: %s", err)
				_ = conn.Close()
				_ = rc.Close()
			}()
			_, err = io.Copy(lc, rt)
			_ = conn.Close()
			_ = rc.Close()
			log.Printf("Error while receiving from remote server: %s", err)
			log.Printf("Connection with %s closed", conn.RemoteAddr().String())
		}(conn)
	}

}
