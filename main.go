package main

import (
	"flag"
	"fmt"
	"time"

	serialWithWriteTimeout "github.com/CreatorsLab/serial"
	serial "github.com/tarm/serial"
)

type Serial interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Flush() error
	Close() (err error)
}

func main() {
	var (
		device  = flag.String("device", "", "Serial device path")
		timeout = flag.Duration("timeout", time.Second*5, "Read timeout")
	)
	flag.Parse()

	if device == nil || *device == "" {
		panic("device required")
	}

	testNonBlockingSerial(*device, *timeout)
	testBlockingSerial(*device, *timeout)
}

func read(port Serial, timeout time.Duration) {
	fmt.Printf("Read with timeout of %v\n", timeout)
	t := time.Now()

	b := make([]byte, 128)
	n, err := port.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d bytes: %v\n", n, b)
	fmt.Printf("Read took %v\n", time.Since(t))
}

func write(port Serial, timeout time.Duration) {
	fmt.Printf("Write with timeout of %v\n", timeout)
	t := time.Now()

	b := []byte{1, 2, 3, 4, 5, 6}
	n, err := port.Write(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Wrote %d bytes: %v\n", n, b)
	fmt.Printf("Write took %v\n", time.Since(t))
}

func testBlockingSerial(device string, timeout time.Duration) {
	fmt.Printf("Blocking serial conn to %s\n", device)

	cfg := serial.Config{
		Name:        device,
		Baud:        9600,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
		ReadTimeout: timeout,
	}
	port, err := serial.OpenPort(&cfg)
	if err != nil {
		panic(err)
	}
	defer port.Close()

	read(port, timeout)
	write(port, timeout)
}

func testNonBlockingSerial(device string, timeout time.Duration) {
	fmt.Printf("Non-blocking serial conn to %s\n", device)

	cfg := serialWithWriteTimeout.Config{
		Name:         device,
		Baud:         9600,
		Size:         8,
		Parity:       serialWithWriteTimeout.ParityNone,
		StopBits:     serialWithWriteTimeout.Stop1,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	port, err := serialWithWriteTimeout.OpenPort(&cfg)
	if err != nil {
		panic(err)
	}
	defer port.Close()

	read(port, timeout)
	write(port, timeout)
}
