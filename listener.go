package main

import (
	"crypto/tls"
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
	ll       net.Listener
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
	return localtext + " " + l.Localport + " -> " + l.Remoteip + l.Remoteport + " " + remotetext + suffix
}
