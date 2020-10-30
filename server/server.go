package server

import (
	"log"
	"net"
)

type Server struct {
	listener         []net.Listener
	conn             net.Conn
	HoldingRegisters []uint16
	modbus_slave      *ModBus
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

	packet := make([]byte, 512)
	bytesRead, _ := s.conn.Read(packet)
	packet = packet[:bytesRead]

	s.modbus_slave = get_packet(packet)
	// s.modbus_slave.print()

	// go s.HandleRequest()

	return err
}

func (s *Server) Close() {
	for _, listen := range s.listener {
		listen.Close()
	}
}

func (s *Server) HandleRequest() {
	// s.conn.Write([]byte("Message received." + "\n"))
	// defer s.conn.Close()
	input := make([]byte, 20)
	n, _ := s.conn.Read(input)
	input = input[0:n]
	if input[6] == 1 {
		s.modbus_slave.print()
	}
}
