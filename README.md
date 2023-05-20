# eep

Eep makes it easy to set up keys with values of known length for storing on an `at24cx` I2C [eeprom chip](https://ww1.microchip.com/downloads/en/DeviceDoc/doc0336.pdf) *(pdf)* and the [driver](https://github.com/tinygo-org/drivers/tree/release/at24cx) included with TinyGo.

Values on an eeprom chip have a length and an offset. Eep allows you to create a collection of keys with their offsets and lengths. The sum of the current key's offset & length is the next key's offset. Lengths shouldn't be exceeded, so pad as necessary.

The Eeper interface has 4 methods:
```go
Read(key string) ([]byte, error)
Write(key string, value []byte) error
Length(key string) (uint16, error)
Is(key string) (*bool, error) // here be dragons
```

### API

**Read** will read the eeprom at the offset for `key` and return the stored value as a byte slice. In the event of an error, such as a key not found in the list of entries, you'll get an empty byte slice and an error.

**Write** will store a value in the eeprom after validating the soundness of the `key` and its `length` against the length of `value`. A successful write will return a nil error.

**Length** will return the `length` for the `key` passed, or an error if it isn't found.

**Is** has some thorns on it. `Is` should be used for reading bool values and optionally single-byte values. Because ambiguity is introduced between `0`, `false`, and `nil` if you do not check your errors, a pointer is returned which will be nil if an error is present (there are several possibilities), rather than to perpetuate the ambiguity in code that doesn't check the error here. Don't panic. Lastly, any retreived byte other than 0 or 1 will be returned as the error alongside a nil pointer. 

### Usage
First, install the package:
```go
go get github.com/eyelight/eep
```

Next, set up your entries with appropriate offsets and expected (eg, padded) length:
```go
entries := []eep.EEPROMEntry{
    {Key: "Name", Offset: 0, Length: 32}, // 32-byte value
    {Key: "MyBool", Offset: 32, Length: 1}, // 1-byte value
    {Key: "MyBytes", Offset: 33, Length: 7}, // 7-byte value
}
```

You are now able to create an Eeper and add your entries into it.

```go
// you must configure the i2c bus first
machine.I2C0.Configure(machine.I2CConfig{})

// create your Eeper
ee := eep.New(&machine.I2C0, entries)
```
