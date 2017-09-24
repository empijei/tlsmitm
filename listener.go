package main

import (
	"crypto/tls"
	"io"
	"net"
)

const encr = "encrypted"
const unen = "unencrypted"

type listener struct {
	Localport, Remoteport, Remoteip string
	//Tells the listener to setup a tls socket
	Secure bool `"json:,omitempty"`
	//Makes the tool change protocol, this can be used to downgrade a secure
	//connection or to upgrade an insecure one
	ProtoSwitch bool `"json:,omitempty"`

	// Used only during deserialization, name of private key of listener
	KeyName string `"json:,omitempty"`
	// Used only during deserialization, name of public key of listener
	CrtName string `"json:,omitempty"`

	//If false does not validate certificates upon connection to remote hosts
	Verify bool `"json:,omitempty"`

	// Unexported field, will be created on load
	certconf *tls.Config
	wwrapper func(io.Writer) io.Writer
	rwrapper func(io.Reader) io.Reader
	ll       net.Listener
}

func (wl *listener) WrapWriter(w io.WriteCloser) io.WriteCloser {
	if wl.wwrapper == nil {
		return w
	}
	return attachCloseToWriter(wl.wwrapper(w), w)
}

func (wl *listener) WrapReader(r io.ReadCloser) io.ReadCloser {
	if wl.rwrapper == nil {
		return r
	}
	return attachCloseToReader(wl.rwrapper(r), r)
}

func (wl *listener) Close() error {
	return wl.ll.Close()
}

func (l *listener) String() string {
	var localtext, remotetext string
	if l.Secure {
		localtext = encr
		if l.ProtoSwitch {
			remotetext = unen
		} else {
			remotetext = encr
		}
	} else {
		localtext = unen
		if l.ProtoSwitch {
			remotetext = encr
		} else {
			remotetext = unen
		}
	}
	suffix := ""
	if l.wwrapper != nil || l.rwrapper != nil {
		suffix = " (wrapped)"
	}
	return localtext + " " + l.Localport + " -> " + l.Remoteip + l.Remoteport + " " + remotetext + suffix
}
