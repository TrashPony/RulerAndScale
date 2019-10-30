package transport_data

import (
	"github.com/jacobsa/go-serial/serial"
	"strconv"
	"sync"
	"time"
)

var Ports = InitPortStorage()

type PortStorage struct {
	Ports map[string]*Port
	mx    sync.Mutex
}

func InitPortStorage() *PortStorage {
	return &PortStorage{
		Ports: make(map[string]*Port),
	}
}

func (p *PortStorage) GetPort(device string) *Port {
	p.mx.Lock()
	defer p.mx.Unlock()

	return p.Ports[device]
}

func (p *PortStorage) ResetPort(device string) {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.Ports[device] = nil
}

func (p *PortStorage) SetPort(port *Port, device string) {

	if port == nil {
		return
	}

	p.mx.Lock()
	defer p.mx.Unlock()

	p.Ports[device] = port
}

func SelectPort() {

	// горутина следит за тем что бы были найдены все йстройства
	println("Поиск портов")
	portClass := []string{"/dev/ttyS", "/dev/ttyACM", "/dev/ttyUSB"}

	for {
		time.Sleep(1000 * time.Millisecond)

		for _, nameClass := range portClass {
			for i := 0; i < 20; i++ {

				portName := nameClass + strconv.Itoa(i)

				if Ports.GetPort("ruler") == nil {
					Ports.SetPort(FindRuler(portName), "ruler")
				}

				if Ports.GetPort("scale") == nil {
					Ports.SetPort(FindScale(portName), "scale")
				}
			}
		}
	}
}

func FindScale(portName string) *Port {
	weightConfig := serial.OpenOptions{
		PortName:              portName,
		BaudRate:              4800,
		ParityMode:            serial.PARITY_EVEN,
		DataBits:              8,
		StopBits:              1,
		InterCharacterTimeout: 300,
	}

	// игнорируем если этот порт уже принадлежит другому устройству
	if Ports.GetPort("ruler") != nil && portName == Ports.GetPort("ruler").Config.PortName {
		return nil
	}

	connect, err := serial.Open(weightConfig)
	if err != nil {
		//println("serial.Open: %v", err.Error())
		return nil
	}

	_, err = connect.Write([]byte{0x48})
	if err != nil {
		connect.Close()
		return nil
	}

	buf := make([]byte, 2)
	n, err := connect.Read(buf)

	if err != nil {
		connect.Close()
		return nil
	} else {
		if n == 2 && (buf[0] == 128 || buf[0] == 192) {
			println("Весы подключены к порту " + portName)
			return &Port{Name: portName, Config: &weightConfig, Connection: connect}
		} else {
			return nil
		}
	}
}

func FindRuler(portName string) *Port {

	rulerConfig := serial.OpenOptions{
		PortName:              portName,
		BaudRate:              115200,
		ParityMode:            serial.PARITY_EVEN,
		DataBits:              8,
		StopBits:              1,
		InterCharacterTimeout: 600,
	}

	// игнорируем если этот порт уже принадлежит другому устройству
	if Ports.GetPort("scale") != nil && portName == Ports.GetPort("scale").Config.PortName {
		return nil
	}

	connect, err := serial.Open(rulerConfig)
	if err != nil {
		//println("serial.Open: %v", err.Error())
		return nil
	}

	_, err = connect.Write([]byte{0x95, 0x95})
	if err != nil {
		connect.Close()
		return nil
	}

	buf := make([]byte, 999)
	_, err = connect.Read(buf)

	if err != nil {
		connect.Close()
		return nil
	} else {
		if buf[0] == 127 {
			println("Линейка подключена к порту " + portName)
			port := &Port{Name: portName, Config: &rulerConfig, Connection: connect}
			return port
		} else {
			return nil
		}
	}
}
