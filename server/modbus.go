package server

import "fmt"
import "encoding/binary"

type ModBus struct {
	id_transaction uint16
	id_protocol    uint16
	length         uint16
	id_unit        uint8
	function       uint8
	data           []byte
}

func get_packet(packet []byte) (*ModBus) {
	if len(packet) < 12 {
		fmt.Println("Not full packet")
	}

	modbus := &ModBus {
		id_transaction : binary.BigEndian.Uint16(packet[0:2]),
		id_protocol    : binary.BigEndian.Uint16(packet[2:4]),
		length         : binary.BigEndian.Uint16(packet[4:6]),
		id_unit        : uint8(packet[6]),
		function       : uint8(packet[7]),
		data           : packet[8:],
	}

	return modbus
}

func (modbus *ModBus) get_first_qty_regs() (first_reg, qty_reg, last_reg uint16) {
	first_reg = binary.BigEndian.Uint16(modbus.data[0:2])
	qty_reg   = binary.BigEndian.Uint16(modbus.data[2:4])
	last_reg  = first_reg + qty_reg
	return first_reg, qty_reg, last_reg
}

func (modbus *ModBus) get_data() (data []byte) {
	return modbus.data
}


