package server

import (
	"log"
	"net"
)

type Server struct {
	listener         []net.Listener
	conn             net.Conn
	HoldingRegisters []uint16
}

func Make() *Server {
	s := &Server{}

	s.HoldingRegisters = make([]uint16, 65536)

	return s
}

func (s *Server) Listen() (err error) {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Printf("Failed to Listen\n", err)
		return err
	}

	s.accept(listen)

	return err
}

func (s *Server) accept(listen net.Listener) (err error) {
	s.conn, err = listen.Accept()
	if err != nil {
		log.Println(err)
		return err
	}

	s.listener = append(s.listener, listen)

	return err
}

func (s *Server) Close() {
	for _, listen := range s.listener {
		listen.Close()
	}
}

func (s *Server) HandleRequest() {
	s.conn.Write([]byte("Message received." + "\n"))
}
