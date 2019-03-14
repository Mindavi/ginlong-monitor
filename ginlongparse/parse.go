package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/Mindavi/ginlong-monitor/dataformat"
	"io/ioutil"
	"log"
)

func main() {
	filenamePtr := flag.String("f", "", "Filename of the binary inverter packet")
	flag.Parse()
	data, err := ioutil.ReadFile(*filenamePtr)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}
	reader := bytes.NewReader(data)
	var invData dataformat.RawInverterData
	err = binary.Read(reader, binary.BigEndian, &invData)
	if err != nil {
		log.Fatal("Invalid binary data", err)
	}

	fmt.Printf("Start: %#x\n", invData.Start)
	fmt.Printf("Length: %#x\n", invData.Length)
	fmt.Printf("Control code: %#x\n", invData.ControlCode)
	fmt.Printf("Id: %#x\n", invData.Id1)
	fmt.Printf("Id2: %#x\n", invData.Id2)
	fmt.Printf("Command: %#x\n", invData.Command)
	fmt.Printf("Protocol: %#x\n", invData.ProtocolType)
	fmt.Printf("Serial: %s\n", invData.Serial)
	fmt.Printf("Status: %#x\n", invData.Status)
	fmt.Printf("End: %#x\n", invData.End)

	newData := dataformat.ConvertInverterData(invData)

	// fmt.Printf("Header:%#x\n", invData.Header)
	fmt.Println("Temperature:", newData.Temperature, "*C")
	fmt.Println("Vdc1:", newData.Vdc1, "V")
	fmt.Println("Vdc2:", newData.Vdc2, "V")
	fmt.Println("Adc1:", newData.Adc1, "A")
	fmt.Println("Adc2:", newData.Adc2, "A")
	fmt.Println("Vac:", newData.Vac, "V")
	fmt.Println("Aac:", newData.Aac, "A")
	fmt.Println("Fac:", newData.Fac, "Hz")
	fmt.Println("Power now:", newData.PNow, "W")
	fmt.Println("Yesterday energy:", newData.Yesterday, "kWh")
	fmt.Println("Today energy:", newData.Today, "kWh")
	fmt.Println("Total energy:", newData.Total, "kWh")
	fmt.Println("Month energy:", newData.Month, "kWh")
	fmt.Println("Last month energy:", newData.LastMonth, "kWh")
	fmt.Println("Inverter status:", newData.Status)
}
