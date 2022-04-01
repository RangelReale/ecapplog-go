package ecapplog

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/RangelReale/ecapplog-go/internal"
)

type Client struct {
	appname      string
	address      string
	bufferSize   int
	isOpen       bool
	ringBuffer   *internal.RingBuffer
	inChan       chan interface{}
	outChan      chan interface{}
	cmdCtx       context.Context
	cmdCtxCancel context.CancelFunc
}

func NewClient(options ...Option) *Client {
	ret := &Client{
		appname:    "ECAPPLOG-GO",
		address:    "127.0.0.1:13991",
		bufferSize: 1000,
		isOpen:     false,
	}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func (c *Client) Open() {
	dropFn := func(m interface{}) {
		c.handleError(fmt.Errorf("Dropped older message"))
	}

	if !c.isOpen {
		c.inChan = make(chan interface{})
		c.outChan = make(chan interface{}, c.bufferSize)
		c.cmdCtx, c.cmdCtxCancel = context.WithCancel(context.Background())
		c.isOpen = true

		c.ringBuffer = internal.NewRingBufferWithDropFn(c.inChan, c.outChan, dropFn)
		go c.ringBuffer.Run()
		go c.handleConnection(c.cmdCtx, c.outChan)
	}
}

func (c *Client) Close() {
	if c.isOpen {
		c.cmdCtxCancel()
		close(c.inChan)

		c.ringBuffer = nil
		c.outChan = nil
		c.inChan = nil
		c.cmdCtx = nil
		c.cmdCtxCancel = nil
		c.isOpen = false
	}
}

func (c *Client) handleConnection(c_cmdCtx context.Context, c_cmdChan chan interface{}) {
rfor:
	for {
		err := func() error {
			var d net.Dialer
			cctx, ccancel := context.WithTimeout(c_cmdCtx, time.Second*10)
			defer ccancel()
			conn, err := d.DialContext(cctx, "tcp", c.address)
			if err != nil {
				return err
			}
			defer conn.Close()

			if conntcp, ok := conn.(*net.TCPConn); ok {
				err = conntcp.SetNoDelay(false)
				if err != nil {
					return err
				}
			}

			// write banner
			err = c.handleBanner(conn)
			if err != nil {
				return err
			}

			for {
				var err error
				select {
				case <-c_cmdCtx.Done():
					return nil
				case cmd := <-c_cmdChan:
					switch xcmd := cmd.(type) {
					case *cmdLog:
						err = c.handleCmdLog(conn, xcmd)
					}
				}
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return nil
					}
					return err
				}
			}
		}()
		if err != nil {
			c.handleError(err)
		}

		select {
		case <-c_cmdCtx.Done():
			break rfor
		case <-time.After(time.Second * 5):
			break
		}
	}
}

func (c *Client) handleBanner(conn net.Conn) error {
	// write command
	err := binary.Write(conn, binary.BigEndian, command_Banner)
	if err != nil {
		return err
	}

	data := []byte(fmt.Sprintf("ECAPPLOG %s", c.appname))

	// write size
	size := int32(len(data))
	err = binary.Write(conn, binary.BigEndian, size)
	if err != nil {
		return err
	}

	// write data
	err = binary.Write(conn, binary.BigEndian, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) handleCmdLog(conn net.Conn, cmd *cmdLog) error {
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	// write command
	err = binary.Write(conn, binary.BigEndian, command_Log)
	if err != nil {
		return err
	}

	// write size
	size := int32(len(jcmd))
	err = binary.Write(conn, binary.BigEndian, size)
	if err != nil {
		return err
	}

	// write encoded json
	err = binary.Write(conn, binary.BigEndian, jcmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) handleError(err error) {

}

func (c *Client) Log(time time.Time, priority Priority, category string, message string,
	options ...LogOption) {
	var lo logOptions
	for _, opt := range options {
		opt(&lo)
	}

	if c.isOpen {
		c.inChan <- &cmdLog{
			Time:            cmdTime{time},
			Priority:        priority,
			Category:        category,
			Message:         message,
			Source:          lo.source,
			ExtraCategories: lo.extraCategories,
		}
	}
}
