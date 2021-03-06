package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Server struct
type Server struct {
	listener         []net.Listener
	conn             net.Conn
	holdingRegisters []uint16
	timeStamp        []int64
	modbus           *ModBus
	lifeTime         int64
}

// ValueTimestamp contents value and time of register
// type ValueTimestamp struct {
// 	value     uint16
// 	timeStamp int64
// }

// Make server
func Make() *Server {
	server := &Server{}

	server.holdingRegisters = make([]uint16, 100)
	server.timeStamp = make([]int64, 100)
	server.holdingRegisters[0] = 12
	server.holdingRegisters[1] = 7

	return server
}

// Listen opening
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

// Close listening
func (server *Server) Close() {
	for _, listen := range server.listener {
		listen.Close()
	}
}

// SetLifeTime sets time of life in seconds
func (server *Server) SetLifeTime(t int64) {
	server.lifeTime = t
}

// answer03 read holding registers, to place in ModBus
func (server *Server) answer03() (answer []byte) {
	firstReg, qtyReg, lastReg := server.modbus.getFirstQtyRegs()

	if int(lastReg) > len(server.holdingRegisters) {
		answer = server.answerError(wrongReg)
		return
	}

	for i := firstReg; i < lastReg; i++ {
		if time.Now().Unix()-server.timeStamp[i] >= server.lifeTime {
			answer = server.answerError(wrongReg)
			return
		}
	}

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

	return
}

// answer16 write holding registers, to place in ModBus
func (server *Server) answer16() (answer []byte) {
	firstReg, qtyReg, lastReg := server.modbus.getFirstQtyRegs()

	if int(lastReg) > len(server.holdingRegisters) {
		answer = server.answerError(wrongReg)
		return
	}

	values := bytesToUint16(server.modbus.getData()[5:])
	if len(values) != int(qtyReg) {
		answer = server.answerError(wrongReg)
		return
	}
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

	for i := firstReg; i < lastReg; i++ {
		server.timeStamp[i] = time.Now().Unix()
	}

	return
}

const (
	wrongFunc byte = 1 << iota
	wrongReg
)

func (server *Server) answerError(errorCode byte) (answer []byte) {
	answer = make([]byte, 9)
	binary.BigEndian.PutUint16(answer[0:2], server.modbus.idTransaction)
	binary.BigEndian.PutUint16(answer[2:4], server.modbus.idProtocol)
	binary.BigEndian.PutUint16(answer[4:6], 3) // length of message in error response
	answer[6] = server.modbus.idUnit
	answer[7] = server.modbus.function + 0b10000000
	answer[8] = errorCode

	return
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

// HandleRequest handler
func (server *Server) HandleRequest() {
	switch server.modbus.function {
	case 0x03:
		server.conn.Write(server.answer03())
	case 0x10:
		server.conn.Write(server.answer16())
	default:
		server.conn.Write(server.answerError(wrongFunc))
	}
}
