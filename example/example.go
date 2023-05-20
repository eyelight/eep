package main

import (
	"machine"

	"github.com/eyelight/eep"
)

var entries = []eep.EEPROMEntry{
	{Key: "DEVICE_NICKNAME", Offset: 0, Length: 32}, // 32-byte value of the device's nickname
	{Key: "TINYGO_COOL", Offset: 32, Length: 1},     // 1-byte bool of whether TinyGo is cool
	{Key: "BYTE", Offset: 33, Length: 1},            // 1-byte value
}

var ee eep.Eeper

func init() {
	machine.I2C0.Configure(machine.I2CConfig{})
}

func main() {
	ee := eep.New(machine.I2C0, entries)

	// Write the device nickname
	if err := ee.Write("DEVICE_NICKNAME", []byte("Example eeprom toy")); err != nil {
		println(err)
	}

	// Write whether TinyGo is cool
	if err := ee.Write("TINYGO_IS_COOL", []byte{1}); err != nil {
		println(err)
	}

	// Write a byte
	b := make([]byte, 0)
	b = append(b, 42)
	if err := ee.Write("BYTE", b); err != nil {
		println(err)
	}

	// Read & print the device nickname
	print("Device Nickname: ")
	d, err := ee.Read("DEVICE_NICKNAME")
	if err != nil {
		println(err)
	}
	println(string(d)) // 'Example eeprom toy'

	// Read & print whether TinyGo is cool
	print("TinyGo is cool: ")
	s, err := ee.Is("TINYGO_COOL")
	if err != nil {
		println(err)
	}
	println(*s) // 'true'

	// Read the random byte
	print("Byte: ")
	_, err = ee.Is("BYTE")
	if err != nil {
		println(err) // '*' or [42]
	}

	// Try to read a nonexistent key
	print("Nonexistent key: ")
	_, err = ee.Read("NONEXISTENT")
	println(err)
}
