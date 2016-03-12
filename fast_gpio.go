package fast_gpio

import (
	"sync"
	"os"
	"syscall"
	"errors"
)

type Direction uint8
type Pin uint8
type State uint8
type Pull uint8

// Memory offsets for gpio, see the spec for more details
const (
	REGISTER_BLOCK_ADDR = 0x18040000
	REGISTER_BLOCK_SIZE = 0x30
	REGISTER_OE_OFFSET = 0
	REGISTER_IN_OFFSET = 1
	REGISTER_OUT_OFFSET = 2
	REGISTER_SET_OFFSET = 3
	REGISTER_CLEAR_OFFSET = 4
)

// Pin direction, a pin can be set in Input or Output mode
const (
	Input Direction = iota
	Output
)

// State of pin, High / Low
const (
	Low State = iota
	High
)

// Pull Up / Down / Off
const (
	PullOff Pull = iota
	PullDown
	PullUp
)

// Arrays for 8 / 32 bit access to memory and a semaphore for write locking
type FastGPIO struct {
	registerLock sync.Mutex
	register     [] byte
}

func(gpio* FastGPIO) Open() error{
	var file *os.File

	file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0)

	if err != nil{
		return err
	}

	defer file.Close()

	gpio.registerLock.Lock()
	defer gpio.registerLock.Unlock()

	gpio.register, err = syscall.Mmap(
		int(file.Fd()),
		REGISTER_BLOCK_ADDR,
		REGISTER_BLOCK_SIZE,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED)

	if err != nil || len(gpio.register) != REGISTER_BLOCK_SIZE{
		return errors.New("Error in reading memory")
	}

	return nil
}