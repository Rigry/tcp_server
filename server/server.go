package server

import (
	"log"
	"net"
	"fmt"
	"encoding/binary"
)

type Server struct {
	listener         []net.Listener
	conn             net.Conn
	holding_registers []uint16
	modbus           *ModBus
}

func Make() *Server {
	server := &Server{}

	server.holding_registers = make([]uint16, 10)
	server.holding_registers[0] = 12
	server.holding_registers[1] = 7

	return server
}

func (server *Server) Listen() (err error) {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Printf("Failed to Listen\n", err)
		return err
	}
	server.listener = append(server.listener, listen)

	go server.accept(listen)

	return err
}

func (server *Server) accept(listen net.Listener) (err error) {
	server.conn, err = listen.Accept()
	if err != nil {
		log.Println(err)
		return err
	}

	go func(conn net.Conn) {
		defer conn.Close()

		for {
			packet := make([]byte, 512)
			bytes, _ := server.conn.Read(packet)
			packet = packet[:bytes]

			server.modbus = get_packet(packet)

			server.HandleRequest()
		}
	}(server.conn)

	return err
}



func (server *Server) Close() {
	for _, listen := range server.listener {
		listen.Close()
	}
}

// func (server *Server) tcp_header() (header []byte) {
	
// }

func (server *Server) answer_03() (answer []byte) {
	first_reg, qty_reg, last_reg := server.modbus.get_first_qty_regs()
	data := make([]byte,1)
	data[0] = byte(qty_reg * 2)
	data = append(data, uint16_to_bytes(server.holding_registers[first_reg:last_reg])...)

	answer = make([]byte, 8)
	binary.BigEndian.PutUint16(answer[0:2], server.modbus.id_transaction)
	binary.BigEndian.PutUint16(answer[2:4], server.modbus.id_protocol)
	binary.BigEndian.PutUint16(answer[4:6], uint16(2+len(data)))
	answer[6] = server.modbus.id_unit
	answer[7] = server.modbus.function
	answer = append(answer, data...)

	return answer
}

func (server *Server) answer_16() (answer []byte){
	first_reg, qty_reg, last_reg := server.modbus.get_first_qty_regs()
	values := bytes_to_uint16(server.modbus.get_data()[5:])
	// ошибка по кол-ву байт
	copy(server.holding_registers[first_reg:last_reg], values)
	
	data := make([]byte,4)
	binary.BigEndian.PutUint16(data[0:2], first_reg)
	binary.BigEndian.PutUint16(data[2:4], qty_reg)

	answer = make([]byte, 8)
	binary.BigEndian.PutUint16(answer[0:2], server.modbus.id_transaction)
	binary.BigEndian.PutUint16(answer[2:4], server.modbus.id_protocol)
	binary.BigEndian.PutUint16(answer[4:6], uint16(2+len(data)))
	answer[6] = server.modbus.id_unit
	answer[7] = server.modbus.function
	answer = append(answer, data...)
	
	return answer
}

func uint16_to_bytes(values []uint16) []byte {
	bytes := make([]byte, len(values)*2)

	for i, value := range values {
		binary.BigEndian.PutUint16(bytes[i*2:(i+1)*2], value)
	}
	return bytes
}

func bytes_to_uint16(bytes []byte) []uint16 {
	values := make([]uint16, len(bytes)/2)

	for i := range values {
		values[i] = binary.BigEndian.Uint16(bytes[i*2 : (i+1)*2])
	}
	return values
}

func (server *Server) HandleRequest() {
	for {
		if server.modbus.id_unit == 1 {
			switch server.modbus.function {
			case 0x03:	
				server.conn.Write(server.answer_03())
			case 0x10:
				server.conn.Write(server.answer_16())
			default:
				fmt.Println("Function " + string(server.modbus.function) + " not realised")
			}
		}
	}
}
