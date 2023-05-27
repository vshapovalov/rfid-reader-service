package readers

import (
	"errors"
	"github.com/tarm/serial"
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/models"
	"time"
)

const (
	commandSearchCard     = 0x02
	commandReadCardNumber = 0x03
	//commandCheckModule    = 0x12
	commandBuzzer = 0x13
	//commandSelectCard     = 0x04
	//commandStopCard       = 0x09
)

type m302Packet struct {
	Command byte
	Data    []byte
}

const portReadTimeout = 2000 * time.Millisecond

func parseRawM302Packet(data []byte) (*m302Packet, error) {
	length := data[0]
	crcChecksum := data[len(data)-1]
	if len(data)-1 != int(length) {
		return nil, errors.New("invalid packet length")
	}
	if crcChecksum != calcCRCChecksum(data[:len(data)-1]) {
		return nil, errors.New("invalid packet checksum")
	}

	return &m302Packet{
		Command: data[1],
		Data:    data[2 : len(data)-1],
	}, nil
}

func newM302Packet(command byte, data []byte) *m302Packet {
	return &m302Packet{
		Command: command,
		Data:    data,
	}
}

func (packet *m302Packet) Serialize() []byte {
	var result []byte
	result = append(result, 2+byte(len(packet.Data)))
	result = append(result, packet.Command)
	result = append(result, packet.Data...)
	result = append(result, calcCRCChecksum(result))
	return result
}

type M302ReaderModule struct {
	port   *serial.Port
	logger infrastructure.ILogger
}

func (m *M302ReaderModule) GetReaderInfo() string {
	return "{}"
}

func NewM302ReaderModule(deviceSettings models.ConfigM302Settings, logger infrastructure.ILogger) (*M302ReaderModule, error) {
	config := &serial.Config{
		Name:        deviceSettings.Port,
		Baud:        deviceSettings.Baud,
		ReadTimeout: deviceSettings.ReadTimeout,
		Size:        deviceSettings.Size,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return &M302ReaderModule{
		port:   port,
		logger: logger,
	}, nil
}

func readPacket(port *serial.Port, timeout time.Duration) (*m302Packet, error) {
	var data []byte
	var packet *m302Packet
	var hasData bool
	var readingInterval = 30 * time.Millisecond

	buf := make([]byte, 100)

	timeoutThreshold := time.Now().Add(timeout)
	for {
		n, err := port.Read(buf)
		if err != nil {
			return nil, err
		}
		if n > 0 {
			hasData = true
			data = append(data, buf[:n]...)
		} else {
			if hasData {
				hasData = false
				packet, err = parseRawM302Packet(data)
				if err != nil {
					return nil, err
				}
				data = nil
				break
			}
		}
		time.Sleep(readingInterval)
		if time.Now().After(timeoutThreshold) {
			if packet != nil {
				break
			} else {
				return nil, errors.New("timeout")
			}
		}
	}

	return packet, nil
}

func (m *M302ReaderModule) Close() error {
	return m.port.Close()
}

func (m *M302ReaderModule) Buzz() error {
	_, err := m.port.Write(newM302Packet(commandBuzzer, nil).Serialize())
	if err != nil {
		return err
	}
	_, _ = readPacket(m.port, portReadTimeout)
	return nil
}

func (m *M302ReaderModule) ReadCards() ([][]byte, error) {
	var lastPacket *m302Packet
	var err error

	_, err = m.port.Write(newM302Packet(commandSearchCard, []byte{0}).Serialize())
	if err != nil {
		return nil, err
	}
	lastPacket, err = readPacket(m.port, portReadTimeout)
	if err == nil && lastPacket.Command == commandSearchCard {
		_, err = m.port.Write(newM302Packet(commandReadCardNumber, nil).Serialize())
		if err != nil {
			return nil, err
		}
		lastPacket, err = readPacket(m.port, portReadTimeout)
		if err == nil && lastPacket.Command == commandReadCardNumber {
			return [][]byte{lastPacket.Data}, nil
		}
	}
	return nil, nil
}
