package server

import (
	"encoding/binary"
	"fmt"
)

// ModBus struct of message 
type ModBus struct {
	idTransaction uint16
	idProtocol    uint16
	length        uint16
	idUnit        uint8
	function      uint8
	data          []byte
}

func getPacket(packet []byte) (*ModBus, error){
	if len(packet) < 12 {
		return nil, fmt.Errorf("Not full packet %v", packet)
	}

	modbus := &ModBus{
		idTransaction: binary.BigEndian.Uint16(packet[0:2]),
		idProtocol:    binary.BigEndian.Uint16(packet[2:4]),
		length:        binary.BigEndian.Uint16(packet[4:6]),
		idUnit:        uint8(packet[6]),
		function:      uint8(packet[7]),
		data:          packet[8:],
	}

	return modbus, nil
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

// 	return answer
// }

// func (modbus *ModBus) answer16() (answer []byte) {

// 	return answer
// }


// func (modbus *ModBus) get_function() (uint8 function) {

// 	return
// }
