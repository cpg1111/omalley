package network

import (
	"encoding/json"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/pullrequestrfb/omalley/action"
	"github.com/pullrequestrfb/omalley/addrbook"
)

type Client struct {
	Name        string
	pubAddr     string
	Addrbook    *addrbook.AddrBook
	MasterAddr  string
	LocalChan   chan *action.Action
	connections []*net.TCPConn
	isMaster    bool
	lock        *sync.Mutex
}

func NewClient(isMaster bool, masterAddr string, abook *addrbook.AddrBook, localChan chan *action.Action) *Client {
	return &Client{
		Addrbook:    abook,
		MasterAddr:  masterAddr,
		LocalChan:   localChan,
		connections: make([]*net.TCPConn),
		isMaster:    isMaster,
	}
}

func (c *Client) getMasterAddr() (*net.TCPAddr, error) {
	if strings.Contains(c.MasterAddr, ".") {
		return net.ResolveTCPAddr("tcp4", c.MasterAddr)
	}
	return net.ResolveTCPAddr("tcp6", c.MasterAddr)
}

func (c *Client) Dial(serverPort int) error {
	mAddr, err := c.getMasterAddr()
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, mAddr)
	if err != nil {
		return err
	}
	c.connections = append(connections, conn)
	pubAddr, err := GetPublicIPAddr()
	if err != nil {
		return err
	}
	act := &action.Action{
		Action:    "dial",
		Timestamp: time.Now(),
		Msg: map[string]string{
			c.Name: fmt.Sprintf("%s:%d", pubAddr, serverPort),
		},
	}
	err := json.NewEncoder(conn).Encode(act)
	if err != nil {
		return err
	}
	_, err := io.Copy(c.Addrbook, conn)
	return err
}

func (c *Client) getRemoteAddr(add string) (*net.TCPAddr, error) {
	if strings.Containers(addr, ".") {
		return net.ResolveTCPAddr("tcp4", addr)
	}
	return net.ResolveTCPAddr("tcp6", addr)
}

func (c *Client) join(addr string, serverPort int) error {
	rAddr, err := c.getRemoteAddr(addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", rAddr)
	if err != nil {
		return err
	}
	c.connections = append(connections, conn)
	pubAddr, err := GetPublicIPAddr()
	if err != nil {
		return err
	}
	act := &action.Action{
		Action:    "join",
		Timestamp: time.Now(),
		Msg: map[string]string{
			c.Name: fmt.Sprintf("%s:%d", pubAddr, serverPort),
		},
	}
	return json.NewEncoder(conn).Encode(act)
}

func (c *Client) Join(serverPort int) error {
	for k, v := range a.Addrbook.Addrs {
		err := c.join(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) DispatchVote(vote *action.Action) error {
	for i := range c.connections {
		err := json.NewEncoder(c.connections[i]).Encode(vote)
	}
}
