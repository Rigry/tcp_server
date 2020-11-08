package server

import "testing"
import "time"

func TestAnswer16 (t *testing.T) {
	server := Make()
	packet := []byte{0,1,0,0,0,0xB,1,0x10,0,0,0,2,4,0,8,0,10}
	server.modbus, _ = getPacket(packet)
	expected := []byte{0,1,0,0,0,6,1,0x10,0,0,0,2}
	got := server.answer16()
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], got[i])
		}
	}
}

func TestAnswer16BrokenPacket (t *testing.T) {
	server := Make()
	packet := []byte{0,1,0,0,0,0xB,1,0x10,0,0,0,2,4,0,8,0}
	server.modbus, _ = getPacket(packet)
	expected := []byte{0,1,0,0,0,3,1,0x90,2}
	got := server.answer16()
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], got[i])
		}
	}
}

func TestAnswer16WrongReg (t *testing.T) {
	server := Make()
	packet := []byte{0,1,0,0,0,0xB,1,0x10,0,0,0,101,4,0,8,0}
	server.modbus, _ = getPacket(packet)
	expected := []byte{0,1,0,0,0,3,1,0x90,2}
	got := server.answer16()
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], got[i])
		}
	}
}

func TestAnswer03 (t *testing.T) {
	server := Make()
	server.holdingRegisters[0] = 12
	server.holdingRegisters[1] = 7
	server.SetLifeTime(3)
	server.timeStamp[0] = time.Now().Unix()
	server.timeStamp[1] = time.Now().Unix()
	packet := []byte{0,1,0,0,0,6,1,0x3,0,0,0,2}
	server.modbus, _ = getPacket(packet)
	expected := []byte{0,1,0,0,0,7,1,0x3,4,0,12,0,7}
	got := server.answer03()
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], got[i])
		}
	}
	time.Sleep(3 * time.Second)
	expected = []byte{0,1,0,0,0,3,1,0x83,2}
	got = server.answer03()
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], got[i])
		}
	}
}

func TestAnswer03WrongReg (t *testing.T) {
	server := Make()
	packet := []byte{0,1,0,0,0,6,1,0x3,0,0,0,101}
	server.modbus, _ = getPacket(packet)
	expected := []byte{0,1,0,0,0,3,1,0x83,2}
	got := server.answer03()
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], got[i])
		}
	}
}

func TestAnsweError (t *testing.T) {
	server := Make()
	packet := []byte{0,1,0,0,0,6,1,6,0,0,0,2}
	server.modbus, _ = getPacket(packet)
	expected := []byte{0,1,0,0,0,3,1,0x86,1}
	got := server.answerError(1)
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], got[i])
		}
	}
}

func TestSetLifeTime (t *testing.T) {
	server := Make()
	server.SetLifeTime(10)
	if server.lifeTime != 10 {
		t.Errorf("expected %v, got %v", 10, server.lifeTime)
	}
}

func TestUint16ToBytes (t *testing.T) {
	values := []uint16{1207, 1990}
	bytes := uint16ToBytes(values)
	if bytes[0] != byte(values[0] >> 8) {
		t.Errorf("expected %v, got %v", byte(values[0] >> 8), bytes[0])
	}
	if bytes[1] != byte(values[0]) {
		t.Errorf("expected %v, got %v", byte(values[0]), bytes[1])
	}
	if bytes[2] != byte(values[1] >> 8) {
		t.Errorf("expected %v, got %v", byte(values[1] >> 8), bytes[0])
	}
	if bytes[3] != byte(values[1]) {
		t.Errorf("expected %v, got %v", byte(values[1]), bytes[1])
	}
}

func TestBytesToUint16 (t *testing.T) {
	bytes := []byte{4, 183, 7, 198}
	values := bytesToUint16(bytes)
	if values[0] != uint16(bytes[0]) << 8 + uint16(bytes[1]) {
		t.Errorf("expected %v, got %v", uint16(bytes[0]) << 8 + uint16(bytes[1]), values[0])
	}
	if values[1] != uint16(bytes[2]) << 8 + uint16(bytes[3]) {
		t.Errorf("expected %v, got %v", uint16(bytes[2]) << 8 + uint16(bytes[3]), values[1])
	}
}
