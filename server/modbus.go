package server

import "fmt"
import "encoding/binary"

type ModBus struct {
	id_transaction uint16
	id_protocol    uint16
	length         uint16
	id_unit        uint8
	function       uint8
	first_reg      uint16
	last_reg       uint16
	qty_reg        uint16
	qty_byte       uint8
	data           []byte
}

func get_packet(packet []byte) (*ModBus) {
	if len(packet) < 12 {
		fmt.Println("Not full packet")
	}

	modbus := &ModBus {}

	modbus.id_transaction = binary.BigEndian.Uint16(packet[0:2])
	modbus.id_protocol    = binary.BigEndian.Uint16(packet[2:4])
	modbus.length         = binary.BigEndian.Uint16(packet[4:6])
	modbus.id_unit        = uint8(packet[6])
	modbus.function       = uint8(packet[7])

	switch modbus.function {
	case 3:	
		modbus.answer_03(packet)
	case 10:
		modbus.answer_16(packet)
	default:
		fmt.Println("Function " + string(modbus.function) + " not realised")
	}

	return modbus
}

func (modbus *ModBus) print() {
	fmt.Println(modbus.id_transaction, modbus.id_protocol, modbus.length, modbus.id_unit, modbus.function, modbus.first_reg)
}

func (modbus *ModBus) answer_03(packet []byte) {
	modbus.get_first_qty_regs(packet)
}

func (modbus *ModBus) answer_16(packet []byte) {
	modbus.get_first_qty_regs(packet)
	modbus.qty_byte  = uint8(packet[12])
	modbus.data      = packet[13:]
}

func (modbus *ModBus) get_first_qty_regs(packet []byte) {
	modbus.first_reg = binary.BigEndian.Uint16(packet[8:10])
	modbus.qty_reg   = binary.BigEndian.Uint16(packet[10:12])
}

