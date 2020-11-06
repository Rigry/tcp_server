package server

import "testing"
import "encoding/binary"

func TestGetPacket(t *testing.T) {
	packet := []byte{0,1,0,0,0,0xB,1,0x10,0,0,0,2,4,0,8,0,10}
	modbus, _ := getPacket(packet)
	if modbus.idTransaction != binary.BigEndian.Uint16(packet[0:2]) {
		t.Errorf("expected %v, got %v", binary.BigEndian.Uint16(packet[0:2]), modbus.idTransaction)
	}
	if modbus.idProtocol != binary.BigEndian.Uint16(packet[2:4]) {
		t.Errorf("expected %v, got %v", binary.BigEndian.Uint16(packet[2:4]), modbus.idProtocol)
	}
	if modbus.length != binary.BigEndian.Uint16(packet[4:6]) {
		t.Errorf("expected %v, got %v", binary.BigEndian.Uint16(packet[4:6]), modbus.length)
	}
	if modbus.idUnit != uint8(packet[6]) {
		t.Errorf("expected %v, got %v", uint8(packet[6]), modbus.idUnit)
	}
	if modbus.function != uint8(packet[7]) {
		t.Errorf("expected %v, got %v", uint8(packet[7]), modbus.function)
	}
	for i := range modbus.data {
		if modbus.data[i] != packet[i + 8] {
			t.Errorf("expected %v, got %v", uint8(packet[i + 8]), modbus.data[i])
		}
	}

}

func TestGetShortPacket(t *testing.T) {
	packet := []byte{0,1,0,0,0,0xB}
	_, err := getPacket(packet)
	if err == nil {
		t.Errorf("expected error not nil, got %v", err)
	}
}

func TestGetFirstQtyRegs(t *testing.T) {
	var modbus ModBus
	modbus.data = []byte{0,8,0,2}
	firstReg, qtyRegs, lastReg := modbus.getFirstQtyRegs()
	if firstReg != 8 {
		t.Errorf("expected %v, got %v", 8, firstReg)
	}
	if qtyRegs != 2 {
		t.Errorf("expected %v, got %v", 2, qtyRegs)
	}
	if lastReg != 10 {
		t.Errorf("expected %v, got %v", 10, lastReg)
	}
}

func TestGetData(t *testing.T) {
	var modbus ModBus
	modbus.data = []byte{0,8,0,2}
	expect := []byte{0,8,0,2}
	got := modbus.getData()
	
	for i := range got {
		if expect[i] != got[i] {
			t.Errorf("expected %v, got %v", expect[i], got[i])
		}
	}
}