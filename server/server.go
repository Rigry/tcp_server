package server

import (
	"encoding/binary"
	"fmt"
	"net"
)

// Server структура сервера
type Server struct {
	listener         []net.Listener
	conn             net.Conn
	holdingRegisters []uint16
	modbus           *ModBus
}

// Make создание сервера
func Make() *Server {
	server := &Server{}

	server.holdingRegisters = make([]uint16, 100)
	// server.holdingRegisters[0] = 12
	// server.holdingRegisters[1] = 7

	return server
}

// Listen открытие соединения
func (server *Server) Listen() (err error) {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Failed to Listen\n", err)
		return err
	}
	server.listener = append(server.listener, listen)

	go server.accept(listen)

	return err
}

func (server *Server) accept(listen net.Listener) (err error) {
	for {
		server.conn, err = listen.Accept()
		if err != nil {
			fmt.Println(err)
			return err
		}

		go func(conn net.Conn) {
			defer server.conn.Close()
			for {
				packet := make([]byte, 512)
				bytes, _ := server.conn.Read(packet)
				packet = packet[:bytes]

				server.modbus, _ = getPacket(packet)
				server.HandleRequest()
			}
		}(server.conn)
	}
}

// Close закрытие всех соединений
func (server *Server) Close() {
	for _, listen := range server.listener {
		listen.Close()
	}
}

// func (server *Server) tcp_header() (header []byte) {

// }
// answer03 ответ на чтение регистров, лучше разместить в Modbus
func (server *Server) answer03() (answer []byte) {
	firstReg, qtyReg, lastReg := server.modbus.getFirstQtyRegs()
	data := make([]byte, 1)
	data[0] = byte(qtyReg * 2)
	data = append(data, uint16ToBytes(server.holdingRegisters[firstReg:lastReg])...)

	answer = make([]byte, 8)
	binary.BigEndian.PutUint16(answer[0:2], server.modbus.idTransaction)
	binary.BigEndian.PutUint16(answer[2:4], server.modbus.idProtocol)
	binary.BigEndian.PutUint16(answer[4:6], uint16(2+len(data)))
	answer[6] = server.modbus.idUnit
	answer[7] = server.modbus.function
	answer = append(answer, data...)

	return answer
}

// // answer16 ответ на запись регистров, лучше разместить в Modbus
func (server *Server) answer16() (answer []byte) {
	firstReg, qtyReg, lastReg := server.modbus.getFirstQtyRegs()
	values := bytesToUint16(server.modbus.getData()[5:])
	// ошибка по кол-ву байт
	copy(server.holdingRegisters[firstReg:lastReg], values)

	data := make([]byte, 4)
	binary.BigEndian.PutUint16(data[0:2], firstReg)
	binary.BigEndian.PutUint16(data[2:4], qtyReg)

	answer = make([]byte, 8)
	binary.BigEndian.PutUint16(answer[0:2], server.modbus.idTransaction)
	binary.BigEndian.PutUint16(answer[2:4], server.modbus.idProtocol)
	binary.BigEndian.PutUint16(answer[4:6], uint16(2+len(data)))
	answer[6] = server.modbus.idUnit
	answer[7] = server.modbus.function
	answer = append(answer, data...)

	return answer
}

func uint16ToBytes(values []uint16) []byte {
	bytes := make([]byte, len(values)*2)

	for i, value := range values {
		binary.BigEndian.PutUint16(bytes[i*2:(i+1)*2], value)
	}
	return bytes
}

func bytesToUint16(bytes []byte) []uint16 {
	values := make([]uint16, len(bytes)/2)

	for i := range values {
		values[i] = binary.BigEndian.Uint16(bytes[i*2 : (i+1)*2])
	}
	return values
}

// HandleRequest обработчик запросов
func (server *Server) HandleRequest() {
	switch server.modbus.function {
	case 0x03:
		server.conn.Write(server.answer03())
	case 0x10:
		server.conn.Write(server.answer16())
	default:
		fmt.Println("Function " + string(server.modbus.function) + " not realised")
	}
}
