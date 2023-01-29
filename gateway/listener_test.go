package gateway

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"golang.org/x/net/nettest"
)

func TestMultiListener(t *testing.T) {
	builder := func() net.Listener {
		l1, err := nettest.NewLocalListener("tcp")
		if err != nil {
			t.Error(err)
		}
		l2, err := nettest.NewLocalListener("tcp")
		if err != nil {
			t.Error(err)
		}
		return NewMultiListener(l1, l2)
	}
	tests := []struct {
		Name    string
		Fn      func(l net.Listener, c1, c2 net.Conn) error
		WantErr bool
	}{
		{
			Name: "Close",
			Fn: func(l net.Listener, c1, c2 net.Conn) error {
				if err := l.Close(); err != nil {
					return err
				}
				return l.Close()
			},
		},
		{
			Name: "Close then accept",
			Fn: func(l net.Listener, c1, c2 net.Conn) error {
				if err := l.Close(); err != nil {
					return err
				}
				if _, err := l.Accept(); !errors.Is(err, net.ErrClosed) {
					return err
				}
				return nil
			},
		},
		{
			Name: "Accept then close",
			Fn: func(l net.Listener, c1, c2 net.Conn) error {
				l.Accept()
				l.Accept()
				errChan := make(chan error, 1)
				go func() {
					_, err := l.Accept()
					errChan <- err
				}()
				time.Sleep(time.Second / 10)
				l.Close()
				if err := <-errChan; !errors.Is(err, net.ErrClosed) {
					return err
				}
				return nil
			},
		},
		{
			Name: "OK",
			Fn: func(l net.Listener, c1, c2 net.Conn) error {
				lc1, err := l.Accept()
				if err != nil {
					return err
				}
				lc2, err := l.Accept()
				if err != nil {
					return err
				}
				if _, err := c1.Write([]byte{0x01}); err != nil {
					return err
				} else if n, err := lc1.Read(make([]byte, 64)); err != nil {
					return err
				} else if n != 1 {
					return fmt.Errorf("read: wrong byte count (%v)", n)
				}
				if _, err := c2.Write([]byte{0x01}); err != nil {
					return err
				} else if n, err := lc2.Read(make([]byte, 64)); err != nil {
					return err
				} else if n != 1 {
					return fmt.Errorf("read: wrong byte count (%v)", n)
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			l := builder()
			addr1 := l.Addr().(*MultiAddr).Addrs()[0]
			addr2 := l.Addr().(*MultiAddr).Addrs()[0]
			c1, err := net.Dial(addr1.Network(), addr1.String())
			if err != nil {
				t.Error(err)
			}
			c2, err := net.Dial(addr2.Network(), addr2.String())
			if err != nil {
				t.Error(err)
			}
			if err := tt.Fn(l, c1, c2); (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}

}
