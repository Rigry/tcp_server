package server

import (
	"encoding/binary"
	"fmt"
)

// ModBus структура модбас-сообщения
type ModBus struct {
	idTransaction uint16
	idProtocol    uint16
	length        uint16
	idUnit        uint8
	function      uint8
	data          []byte
}

func getPacket(packet []byte) *ModBus {
	if len(packet) < 12 {
		fmt.Println("Not full packet")
	}

	modbus := &ModBus{
		idTransaction: binary.BigEndian.Uint16(packet[0:2]),
		idProtocol:    binary.BigEndian.Uint16(packet[2:4]),
		length:        binary.BigEndian.Uint16(packet[4:6]),
		idUnit:        uint8(packet[6]),
		function:      uint8(packet[7]),
		data:          packet[8:],
	}

	return modbus
}

func (modbus *ModBus) getFirstQtyRegs() (firstReg, qtyReg, lastReg uint16) {
	firstReg = binary.BigEndian.Uint16(modbus.data[0:2])
	qtyReg = binary.BigEndian.Uint16(modbus.data[2:4])
	lastReg = firstReg + qtyReg
	return firstReg, qtyReg, lastReg
}

func (modbus *ModBus) getData() (data []byte) {
	return modbus.data
}

// func (modbus *ModBus) answer03() (answer []byte) {
// 	firstReg, qtyReg, lastReg := modbus.getFirstQtyRegs()
// 	data := make([]byte, 1)
// 	data[0] = byte(qtyReg * 2)
// 	data = append(data, uint16ToBytes(regs[firstReg:lastReg])...)

// 	answer = make([]byte, 8)
// 	binary.BigEndian.PutUint16(answer[0:2], modbus.idTransaction)
// 	binary.BigEndian.PutUint16(answer[2:4], modbus.idProtocol)
// 	binary.BigEndian.PutUint16(answer[4:6], uint16(2+len(data)))
// 	answer[6] = modbus.idUnit
// 	answer[7] = modbus.function
// 	answer = append(answer, data...)

// 	return answer
// }

// func (modbus *ModBus) answer16() (answer []byte) {
// 	firstReg, qtyReg, lastReg := modbus.getFirstQtyRegs()
// 	values := bytesToUint16(modbus.getData()[5:])
// 	// ошибка по кол-ву байт
// 	// copy(server.holdingRegisters[firstReg:lastReg], values)

// 	data := make([]byte, 4)
// 	binary.BigEndian.PutUint16(data[0:2], firstReg)
// 	binary.BigEndian.PutUint16(data[2:4], qtyReg)

// 	answer = make([]byte, 8)
// 	binary.BigEndian.PutUint16(answer[0:2], modbus.idTransaction)
// 	binary.BigEndian.PutUint16(answer[2:4], modbus.idProtocol)
// 	binary.BigEndian.PutUint16(answer[4:6], uint16(2+len(data)))
// 	answer[6] = modbus.idUnit
// 	answer[7] = modbus.function
// 	answer = append(answer, data...)

// 	return answer
// }


// func (modbus *ModBus) get_function() (uint8 function) {

// 	return
// }
