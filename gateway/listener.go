package gateway

import (
	"errors"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type MultiAddr struct {
	addrs []net.Addr
}

func NewMultiAddr(listeners ...net.Listener) net.Addr {
	addrs := []net.Addr{}
	for _, l := range listeners {
		addrs = append(addrs, l.Addr())
	}
	return &MultiAddr{addrs: addrs}
}

func (a *MultiAddr) Addrs() []net.Addr {
	return a.addrs
}

func (a *MultiAddr) Network() string {
	if len(a.addrs) == 0 {
		return ""
	}
	return a.addrs[0].Network()
}

func (a *MultiAddr) String() string {
	if len(a.addrs) == 0 {
		return ""
	}
	return a.addrs[0].String()
}

type MultiListener struct {
	closed     atomic.Bool
	acceptChan chan *acceptResp

	addr      net.Addr
	listeners []net.Listener
}

type acceptResp struct {
	conn net.Conn
	err  error
}

func NewMultiListener(listeners ...net.Listener) net.Listener {
	ac := make(chan *acceptResp, 64)
	for _, listener := range listeners {
		go func(l net.Listener) {
			for {
				conn, err := l.Accept()
				select {
				case ac <- &acceptResp{conn: conn, err: err}:
				default:
				}
				if errors.Is(err, net.ErrClosed) {
					break
				}
			}
		}(listener)
	}
	return &MultiListener{
		acceptChan: ac,
		addr:       NewMultiAddr(listeners...),
		listeners:  listeners,
	}
}

func (l *MultiListener) Accept() (net.Conn, error) {
	for {
		if l.closed.Load() {
			return nil, net.ErrClosed
		}
		resp := <-l.acceptChan
		if resp.err != nil {
			return nil, resp.err
		} else if resp.conn != nil {
			resp.conn.SetReadDeadline(time.Now())
			if _, err := resp.conn.Read([]byte{}); err == io.EOF {
				resp.conn.Close()
			} else {
				resp.conn.SetReadDeadline(time.Time{})
				return resp.conn, resp.err
			}
		}
	}
}

func (l *MultiListener) Close() error {
	if l.closed.Load() {
		return nil
	}
	var err error
	for _, listener := range l.listeners {
		if cerr := listener.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}
	l.closed.Store(true)
	return err
}

func (l *MultiListener) Addr() net.Addr {
	return l.addr
}
