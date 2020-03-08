package dataformat

import (
	"errors"
	"math" // for sqrt
)

const (
	ready      uint16 = 0x0000
	generating uint16 = 0x0003
	// other codes might be found in the link(s) below
	// https://github.com/Crosenhain/ginlong_poller/blob/b966488e6371a93c9307268142a1ff3395a500e2/ginlong_rs485_protocol.pdf
	// https://github.com/XtheOne/Inverter-Data-Logger/files/1967100/Solis.Three.phase.communication.protocal.pdf
	noGrid         uint16 = 0x1015
	ExpectedLength int    = 103
	packetLength   byte   = 89
)

type InverterData struct {
	Temperature float64
	Vdc1        float64
	Vdc2        float64
	Adc1        float64
	Adc2        float64
	Aac         float64
	Vac         float64
	Fac         float64
	PNow        uint16
	Yesterday   float64
	Today       float64
	Total       float64
	Month       uint16
	LastMonth   uint16
	Status      string
}

type RawInverterData struct {
	Start        byte // always 0x68
	Length       byte // length of the data (from command to checksum)
	ControlCode  [2]byte
	Id1          [4]byte
	Id2          [4]byte
	Command      byte
	ProtocolType [2]byte
	Serial       [16]byte
	Temperature  uint16
	Vdc1         uint16
	Vdc2         uint16
	_            [2]byte
	Adc1         uint16
	Adc2         uint16
	_            [2]byte
	Aac          uint16
	_            [4]byte
	Vac          uint16
	_            [4]byte
	Fac          uint16
	PNow         uint16
	_            [6]byte
	Yesterday    uint16
	Today        uint16
	Total        uint32
	_            [4]byte
	Status       uint16
	_            [6]byte
	Month        uint16
	_            [2]byte
	LastMonth    uint16
	_            [8]byte
	Checksum     byte // sum of all bytes in data (from length to here, checksum) mod 256
	End          byte // always 0x16
}

func statusToString(status uint16) string {
	switch status {
	case ready:
		return "Ready"
	case generating:
		return "Generating"
	case noGrid:
		return "No grid"
	default:
		return "Unknown"
	}
}

func ConvertInverterData(rawData RawInverterData) (InverterData, error) {
	var data InverterData
	if rawData.Length != packetLength {
		return data, errors.New("Invalid packet length")
	}

	data.Temperature = float64(rawData.Temperature) / 10
	data.Vdc1 = float64(rawData.Vdc1) / 10
	data.Vdc2 = float64(rawData.Vdc2) / 10
	data.Adc1 = float64(rawData.Adc1) / 10
	data.Adc2 = float64(rawData.Adc2) / 10
	data.Aac = float64(rawData.Aac) / 10
	data.Vac = float64(rawData.Vac) / 10
	data.Fac = float64(rawData.Fac) / 100
	// Because this is a three-phase inverter, and the inverter only reports
	// data from one phase, we need to multiply this with sqrt(3).
	// That'll not be exactly correct, but at least a lot better than the incorrect value.
	data.PNow = uint16(float64(rawData.PNow) * math.Sqrt(3))
	data.Yesterday = float64(rawData.Yesterday) / 100
	data.Today = float64(rawData.Today) / 100
	data.Total = float64(rawData.Total) / 10
	data.Month = rawData.Month
	data.LastMonth = rawData.LastMonth
	if rawData.Status == generating && data.PNow == 0 {
		// assume the inverter is done for today
		data.Status = "Off"
	} else {
		data.Status = statusToString(rawData.Status)
	}

	return data, nil
}
