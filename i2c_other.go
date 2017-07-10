// +build !linux

/*
 Package i2c provides low level control over the linux i2c bus.

 Before usage you should load the i2c-dev kernel module

       sudo modprobe i2c-dev

 Each i2c bus can address 127 independent i2c devices, and most
 linux systems contain several buses.

 Source: https://github.com/d2r2/go-i2c
 */
package adcpi

import (
	"fmt"
	"os"
	"syscall"
)

const (
	I2C_SLAVE = 0x0703
)

/*
 I2C represents a connection to an i2c device.
 */
type I2C struct {
	rc *os.File
}

/*
 New opens a connection to an i2c device.
 */
func NewI2C(addr uint8, bus int) (*I2C, error) {
	return nil, nil
}

/*
 Write sends buf to the remote i2c device. The interpretation of
 the message is implementation dependant.
 */
func (this *I2C) Write(buf []byte) (int, error) {
	return nil
}

func (this *I2C) WriteByte(b byte) (int, error) {
	return nil
}

func (this *I2C) Read(p []byte) (int, error) {
	return nil
}

func (this *I2C) Close() error {
	return nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Read byte from i2c device register specified in reg.
 */
func (this *I2C) ReadRegU8(reg byte) (byte, error) {
	return 0, nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Write byte to i2c device register specified in reg.
 */
func (this *I2C) WriteRegU8(reg byte, value byte) error {
	return nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Read unsigned big endian word (16 bits) from i2c device
 starting from address specified in reg.
 */
func (this *I2C) ReadRegU16BE(reg byte) (uint16, error) {
	return 0, nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Read unsigned little endian word (16 bits) from i2c device
 starting from address specified in reg.
 */
func (this *I2C) ReadRegU16LE(reg byte) (uint16, error) {
	return 0, nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Read signed big endian word (16 bits) from i2c device
 starting from address specified in reg.
 */
func (this *I2C) ReadRegS16BE(reg byte) (int16, error) {
	return 0, nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Read unsigned little endian word (16 bits) from i2c device
 starting from address specified in reg.
 */
func (this *I2C) ReadRegS16LE(reg byte) (int16, error) {
	return 0, nil

}

/*
 SMBus (System Management Bus) protocol over I2C.
 Write unsigned big endian word (16 bits) value to i2c device
 starting from address specified in reg.
 */
func (this *I2C) WriteRegU16BE(reg byte, value uint16) error {
	return nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Write unsigned big endian word (16 bits) value to i2c device
 starting from address specified in reg.
 */
func (this *I2C) WriteRegU16LE(reg byte, value uint16) error {
	return nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Write signed big endian word (16 bits) value to i2c device
 starting from address specified in reg.
 */
func (this *I2C) WriteRegS16BE(reg byte, value int16) error {
	return nil
}

/*
 SMBus (System Management Bus) protocol over I2C.
 Write signed big endian word (16 bits) value to i2c device
 starting from address specified in reg.
 */
func (this *I2C) WriteRegS16LE(reg byte, value int16) error {
	return nil
}

func ioctl(fd, cmd, arg uintptr) error {
	return nil
}