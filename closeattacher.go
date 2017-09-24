package main

import (
	"errors"
	"io"
)

func attachCloseToWriter(w io.Writer, c io.Closer) io.WriteCloser {
	if wc, ok := w.(io.WriteCloser); ok {
		return &writeCloseWrapper{w, []io.Closer{wc, c}}
	}
	return &writeCloseWrapper{w, []io.Closer{c}}
}

func attachCloseToReader(r io.Reader, c io.Closer) io.ReadCloser {
	if rc, ok := r.(io.ReadCloser); ok {
		return &readCloseWrapper{r, []io.Closer{rc, c}}
	}
	return &readCloseWrapper{r, []io.Closer{c}}
}

type readCloseWrapper struct {
	r io.Reader
	c []io.Closer
}

func (rcw *readCloseWrapper) Read(p []byte) (n int, err error) {
	return rcw.r.Read(p)
}

func (rcw *readCloseWrapper) Close() error {
	return combineClose(rcw.c)
}

type writeCloseWrapper struct {
	w io.Writer
	c []io.Closer
}

func (rcw *writeCloseWrapper) Write(p []byte) (n int, err error) {
	return rcw.w.Write(p)
}

func (rcw *writeCloseWrapper) Close() error {
	return combineClose(rcw.c)
}

func combineClose(cs []io.Closer) error {
	var err error
	for _, c := range cs {
		suberr := c.Close()
		if suberr != nil {
			if err != nil {
				err = errors.New(suberr.Error() + "," + err.Error())
			}
			err = suberr
		}
	}
	return err
}
