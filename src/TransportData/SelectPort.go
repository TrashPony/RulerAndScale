package TransportData

import (
	"github.com/tarm/serial"
	"strconv"
	"time"
)

func SelectPort() (scalePort *serial.Port, rulerPort *serial.Port) {

	portClass := []string{"/dev/ttyS", "/dev/ttyACM", "/dev/ttyUSB"}

	for {
		for _, nameClass := range portClass {
			for i := 0; i < 10; i++ {

				portName := nameClass + strconv.Itoa(i)

				if scalePort == nil {
					scalePort = FindScale(portName)
				}

				if rulerPort == nil {
					rulerPort = FindRuler(portName)
				}

				if scalePort != nil && rulerPort != nil {
					println("Все устройства подключены.")
					return
				}
			}
		}
	}
}

func FindScale(portName string) (port *serial.Port) {
	weightConfig := &serial.Config{Name: portName,
		Baud: 4800,
		Parity: 'E',
		ReadTimeout: time.Millisecond * 100}

	port, err := serial.OpenPort(weightConfig)

	if err != nil {
		return nil
	}

	_, err = port.Write([]byte{0x45})
	if err != nil {
		port.Close()
		return nil
	}

	time.Sleep(time.Millisecond * 100)

	buf := make([]byte, 2)
	n, err := port.Read(buf)

	if err != nil {
		port.Close()
		return nil
	} else {
		if n == 2 && buf[0] == 128 {
			println("Весы подключены к порту " + portName)
			return port
		} else {
			return nil
		}
	}
}

func FindRuler(portName string) (port *serial.Port) {

	rulerConfig := &serial.Config{Name: portName,
		Baud: 115200,
		ReadTimeout: time.Millisecond * 100}

	port, err := serial.OpenPort(rulerConfig)

	if err != nil {
		return nil
	}

	_, err = port.Write([]byte{0x95})
	if err != nil {
		port.Close()
		return nil
	}

	time.Sleep(time.Millisecond * 100)

	buf := make([]byte, 5)
	n, err := port.Read(buf)

	if err != nil {
		port.Close()
		return nil
	} else {
		if n == 5 && buf[0] == 127 {
			println("Линейка подключена к порту " + portName)
			return port
		} else {
			return nil
		}
	}
}
