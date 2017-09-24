package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

type header struct {
	size    uint
	version uint16
}

func (h *header) marshal() []byte {
	return []byte{
		//size
		byte(h.size >> 24),
		byte(h.size >> 16),
		byte(h.size >> 8),
		byte(h.size),
		//version
		byte(h.version >> 8),
		byte(h.version),
	}
}

func unmarshalHeader(b []byte) *header {
	if len(b) != 6 {
		return nil
	}
	h := &header{}
	h.size = uint(b[0]<<24) + uint(b[1]<<16) + uint(b[2]<<8) + uint(b[3])
	h.version = uint16(b[4]<<8) + uint16(b[5])
	return h
}

type gzipwrapper struct {
	rm             sync.Mutex
	deflatedReader io.Reader
	wm             sync.Mutex
	writeBuffer    []byte
	versionHeader  uint16
	wrapped        io.ReadWriteCloser
}

func (g *gzipwrapper) Read(p []byte) (n int, err error) {
	g.rm.Lock()
	defer g.rm.Unlock()
	if g.deflatedReader == nil {
		//This is a new packet, let's read the header and skip the version
		rawheader := make([]byte, 0, 6)
		n, err = io.ReadAtLeast(g.wrapped, rawheader, 6)
		if err != nil {
			return
		}
		header := unmarshalHeader(rawheader)
		readBuffer := make([]byte, 0, header.size-6)
		n, err = io.ReadAtLeast(g.wrapped, readBuffer, len(readBuffer))
		g.deflatedReader, err = gzip.NewReader(bytes.NewReader(readBuffer))
		if err != nil {
			return
		}
	}
	n, err = g.deflatedReader.Read(p)
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		g.deflatedReader = nil
	}
	return
}

func (g *gzipwrapper) Write(p []byte) (n int, err error) {
	g.wm.Lock()
	defer g.wm.Unlock()
	panic("not implemented")
	//TODO if finished, nil the buffer
}

func (g *gzipwrapper) Close() error {
	return g.wrapped.Close()
}
