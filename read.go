package lineio

import (
	"errors"
	"io"
	"os"

	"github.com/l4go/task"
)

const bufSize = 4096

type Reader struct {
	io     io.Reader
	send   chan []byte
	err    error
	cancel task.Canceller

	buf   []byte
	char  []byte
	delim byte
}

func NewReader(rio io.Reader) *Reader {
	return NewReaderByDelim(rio, '\n')
}

func NewReaderByDelim(rio io.Reader, delim byte) *Reader {
	r := &Reader{
		io:     rio,
		send:   make(chan []byte),
		err:    nil,
		cancel: task.NewCancel(),

		char:  make([]byte, 1),
		buf:   nil,
		delim: delim,
	}

	go r.start()
	return r
}

func (r *Reader) Err() error {
	return r.err
}

func (r *Reader) Recv() <-chan []byte {
	return r.send
}

func (r *Reader) Close() {
	r.cancel.Cancel()
}

func (r *Reader) init_buf() {
	r.buf = make([]byte, 0, bufSize)
}

func (r *Reader) flush() {
	select {
	case <-r.cancel.RecvCancel():
		r.err = ErrCancel
	case r.send <- r.buf:
		r.init_buf()
	}
}

func (r *Reader) push() {
	r.buf = append(r.buf, r.char[0])
}

func (r *Reader) start() {
	defer close(r.send)

	r.init_buf()

read_loop:
	for r.err == nil {
		_, e := io.ReadFull(r.io, r.char)

		select {
		default:
		case <-r.cancel.RecvCancel():
			r.err = ErrCancel
			break read_loop
		}

		switch {
		case e == nil:
			if r.char[0] == r.delim {
				r.flush()
			} else {
				r.push()
			}
		case e == io.EOF:
			if len(r.buf) > 0 {
				r.flush()
			}
			break read_loop
		case errors.Is(e, os.ErrClosed):
			break read_loop
		default:
			r.err = e
			break read_loop
		}
	}
}
