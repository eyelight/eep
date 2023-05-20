package eep

import (
	"errors"
	"machine"
	"strconv"

	"tinygo.org/x/drivers/at24cx"
)

const (
	ERR_KEY_NOT_FOUND  = "key not found in list of eeprom entries"
	ERR_VALUE_TOO_LONG = "could not write value to eeprom; max length for key exceeded"
	ERR_NOT_BOOL_KEY   = "key is not a boolean style entry"
)

type EEPROMEntry struct {
	Key    string
	Offset uint16
	Length uint16
}

type eeper struct {
	entries []EEPROMEntry
	device  at24cx.Device
}

type Eeper interface {
	Read(key string) ([]byte, error)      // returns the byte slice of the value and/or an error
	Write(key string, value []byte) error // returns nil upon successful write of value for key; fails with error
	Length(key string) (uint16, error)    // returns the allocated length of the key / max length in bytes of the value; error if not found
	Is(key string) (*bool, error)         // for bool values, returns *bool + error if not found; for stored values other than 0|1, byte value is returned in the error
}

// New returns a new Eeper; the I2C bus must be configured prior to calling this
func New(bus *machine.I2C, vals []EEPROMEntry) Eeper {
	e := at24cx.New(bus)
	return &eeper{
		entries: vals,
		device:  e,
	}
}

func (e *eeper) Read(key string) ([]byte, error) {
	// check if the key exists; return early
	ee, err := e.findEntry(key)
	if err != nil {
		return nil, errors.New(ERR_KEY_NOT_FOUND)
	}
	// create a byte slice of the length of the key to be read
	b := make([]byte, 0, ee.Length)

	// use ReadAt to populate the byte slice from the entry's offset & return bytes & err
	_, err = e.device.ReadAt(b, int64(ee.Offset))
	return b, err
}

func (e *eeper) Write(key string, value []byte) error {
	// check the key exists; return early
	ee, err := e.findEntry(key)
	if err != nil {
		return errors.New(ERR_KEY_NOT_FOUND)
	}
	// check the value is within the max length for key; return early
	if len(value) > int(ee.Length) {
		return errors.New(ERR_VALUE_TOO_LONG)
	}
	// use WriteAt write the value at the entry's offset & return the error
	_, err = e.device.WriteAt(value, int64(ee.Offset))
	return err
}

// Length returns the max length for the given key, or fails with a key not found error
func (e *eeper) Length(key string) (uint16, error) {
	ee, err := e.findEntry(key)
	if err != nil {
		return 0, err
	}
	return ee.Length, nil
}

// Is returns the bool value for a key, or fails with an error
// if the byte has a value other than 0 or 1, it will be returned in the error
func (e *eeper) Is(key string) (*bool, error) {
	// check if key exists
	ee, err := e.findEntry(key)
	if err != nil {
		return nil, err
	}
	// fail if key is longer than 1; we don't want to read a partial value
	if ee.Length > 1 {
		return nil, errors.New("(" + key + ") " + ERR_NOT_BOOL_KEY)
	}

	b := uint8(0)
	b, err = e.device.ReadByte(ee.Offset)
	if b == 1 {
		f := true
		return &f, nil
	} else if b == 0 {
		f := false
		return &f, nil
	}
	return nil, errors.New(strconv.FormatUint(uint64(b), 10))
}

// findEntry finds and returns a pointer to the EEPROMEntry with Key 'k' or fails with an error
func (e *eeper) findEntry(k string) (EEPROMEntry, error) {
	for _, entry := range e.entries {
		if entry.Key == k {
			return entry, nil
		}
	}
	return EEPROMEntry{}, errors.New("(" + k + ") " + ERR_KEY_NOT_FOUND)
}
