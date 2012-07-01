// The RS232 package lets you access old serial junk from go
package rs232

/*
#include <stdlib.h>
#include <fcntl.h>
#include <termios.h>
*/
import "C"

import (
	"os"
	"syscall"
)

// This is your serial port handle.
type SerialPort struct {
	port *os.File
}

func baudConversion(rate int) (flag _Ctype_speed_t) {
	if C.B9600 != 9600 {
		panic("Baud rates may not map directly.")
	}
	return _Ctype_speed_t(rate)
}

// SerConf represents the basic serial configuration to provide to OpenPort.
type SerConf int

const (
	S_8N1 = iota
	S_7E1
	S_7O1
)

// Opens and returns a non-blocking serial port.
// The device, baud rate, and SerConf is specified.
//
// Example:  rs232.OpenPort("/dev/ttyS0", 115200, rs232.S_8N1)
func OpenPort(port string, baudRate int, serconf SerConf) (rv SerialPort, err error) {
	f, open_err := os.OpenFile(port,
		syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY,
		0666)
	if open_err != nil {
		err = open_err
		return
	}
	rv.port = f

	fd := rv.port.Fd()

	var options C.struct_termios
	if C.tcgetattr(C.int(fd), &options) < 0 {
		panic("tcgetattr failed")
	}

	if C.cfsetispeed(&options, baudConversion(baudRate)) < 0 {
		panic("cfsetispeed failed")
	}
	if C.cfsetospeed(&options, baudConversion(baudRate)) < 0 {
		panic("cfsetospeed failed")
	}
	switch serconf {
	case S_8N1:
		{
			options.c_cflag &^= C.PARENB
			options.c_cflag &^= C.CSTOPB
			options.c_cflag &^= C.CSIZE
			options.c_cflag |= C.CS8
		}
	case S_7E1:
		{
			options.c_cflag |= C.PARENB
			options.c_cflag &^= C.PARODD
			options.c_cflag &^= C.CSTOPB
			options.c_cflag &^= C.CSIZE
			options.c_cflag |= C.CS7
		}
	case S_7O1:
		{
			options.c_cflag |= C.PARENB
			options.c_cflag |= C.PARODD
			options.c_cflag &^= C.CSTOPB
			options.c_cflag &^= C.CSIZE
			options.c_cflag |= C.CS7
		}
	}
	// Local
	options.c_cflag |= (C.CLOCAL | C.CREAD)
	// no hardware flow control
	options.c_cflag &^= C.CRTSCTS

	if C.tcsetattr(C.int(fd), C.TCSANOW, &options) < 0 {
		panic("tcsetattr failed")
	}

	if syscall.SetNonblock(int(fd), false) != nil {
		panic("Error disabling blocking")
	}

	return
}

// Read from the port.
func (port *SerialPort) Read(p []byte) (n int, err error) {
	return port.port.Read(p)
}

// Write to the port.
func (port *SerialPort) Write(p []byte) (n int, err error) {
	return port.port.Write(p)
}

// Close the port.
func (port *SerialPort) Close() error {
	return port.port.Close()
}
