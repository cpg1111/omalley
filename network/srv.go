package network

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"

	"github.com/pullrequestrfb/omalley/action"
	"github.com/pullrequestrfb/omalley/addrbook"
)

type SRV struct {
	Addrbook    *addrbook.AddrBook
	Host        string
	Port        int
	Name        string
	MasterAddr  string
	LocalChan   chan *action.Action
	Elector     *elect.Elector
	listener    *net.TCPListener
	connections []*net.TCPConn
	isMaster    bool
}

func New(isMaster bool, master, name, host string, port int, abook addrbook.AddrBook, elector *elect.Elector) *SRV {
	localChan := make(chan *action.Action)
	return &SRV{
		Addrbook:    abook,
		Host:        host,
		Port:        port,
		Name:        name,
		MasterAddr:  master,
		LocalChan:   localChan,
		Elector:     elect.New(localChan),
		connections: make([]*net.TCPConn),
		isMaster:    isMaster,
	}
}

func (s *SRV) getHostAddr() (*net.TCPAddr, error) {
	if strings.Contains(s.Host, ":") {
		return net.ResolveTCPAddr("tcp6", fmt.Sprintf("[%s]:%d", s.Host, s.Port))
	}
	return net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.Host, s.Port))
}

func (s *SRV) saveAddr(conn *net.TCPConn, msg map[string]string) (bool, error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return false, err
	}
	_, err := s.Addrbook.Write(msgBytes)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SRV) delegateAction(conn *net.TCPConn, act *action.Action) (bool, error) {
	switch act.Action {
	case "dial":
		if !s.isMaster {
			return false, nil
		}
		return s.saveAddr(conn, act.Msg)
	case "join":
		return s.saveAddr(conn, act.Msg)
	case "vote":
		return s.Elector.Recv(conn, act.Msg)
	case "elect":
		return s.Elector.Confirm(conn, act.Msg)
	}
}

func (s *SRV) handleConn(conn *net.TCPConn) {
	buf := make([]byte, 4096)
	read := true
	defer s.CloseConn(conn)
	for read {
		act := &action.Action{}
		decoder := json.NewDecoder(conn)
		err := decoder.Decode(act)
		if err != nil {
			log.Error(err)
			return
		}
		read, err = s.delegateAction(conn, act)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *SRV) Run() error {
	addr, err := s.getHostAddr()
	if err != nil {
		return err
	}
	s.listener, err = net.ListenerTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		s.connections = append(s.connections, conn)
		go s.handleConn(conn)
	}
}
