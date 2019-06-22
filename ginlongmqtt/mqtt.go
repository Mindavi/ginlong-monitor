package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Mindavi/ginlong-monitor/dataformat"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"net"
	"os"
	"time"
)

type configuration struct {
	ListenPort        string
	MqttServerAddress string
	MqttServerPort    string
	MqttUsername      string
	MqttPassword      string
	MqttTopic         string
	MqttClientId      string
}

const (
	host      = ""
	conn_type = "tcp"
)

func defaultConfig() configuration {
	var config configuration

	config.ListenPort = "9999"
	config.MqttServerAddress = "tcp://127.0.0.1"
	config.MqttServerPort = "1883"
	config.MqttClientId = "ginlong-inverter-monitor"
	config.MqttTopic = "sensor/inverter/" + config.MqttClientId + "/status"

	return config
}

func readConfig() configuration {
	config := defaultConfig()

	setIfOk := func(key string, configkey *string) {
		val, ok := os.LookupEnv(key)
		if ok {
			*configkey = val
		}
	}

	errIfNok := func(key string, configkey *string) {
		val, ok := os.LookupEnv(key)
		if !ok {
			log.Fatal("Expected environment variable " + key + " to be set")
		} else {
			*configkey = val
		}
	}

	setIfOk("INVERTER_LISTENPORT", &config.ListenPort)
	setIfOk("MQTT_CLIENTID", &config.MqttClientId)
	setIfOk("MQTT_SERVERADDRESS", &config.MqttServerAddress)
	setIfOk("MQTT_SERVERPORT", &config.MqttServerPort)
	setIfOk("MQTT_INVERTER_TOPIC", &config.MqttTopic)

	errIfNok("MQTT_USERNAME", &config.MqttUsername)
	errIfNok("MQTT_PASSWORD", &config.MqttPassword)

	log.Print("INVERTER_LISTENPORT: ", config.ListenPort)
	log.Print("MQTT_CLIENTID: ", config.MqttClientId)
	log.Print("MQTT_SERVERADDRESS:MQTT_SERVERPORT: ", config.MqttServerAddress, ":", config.MqttServerPort)
	log.Print("MQTT_INVERTER_TOPIC: ", config.MqttTopic)

	return config
}

func setupMqtt(config configuration) mqtt.Client {
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(config.MqttServerAddress + ":" + config.MqttServerPort).SetClientID(config.MqttClientId)
	opts.SetKeepAlive(30 * time.Second)
	opts.Username = config.MqttUsername
	opts.Password = config.MqttPassword

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return c
}

func runServer(config configuration, client mqtt.Client) {
	server, err := net.Listen(conn_type, host+":"+config.ListenPort)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		log.Print("accepted connection:", conn.RemoteAddr())
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(conn, client, config)
	}
}

func main() {
	config := readConfig()
	client := setupMqtt(config)
	runServer(config, client)
}

func convertData(data []byte) (dataformat.InverterData, error) {
	reader := bytes.NewReader(data)
	var invData dataformat.RawInverterData
	err := binary.Read(reader, binary.BigEndian, &invData)
	if err != nil {
		log.Fatal("Invalid binary data", err)
	}
	return dataformat.ConvertInverterData(invData)
}

func postData(data dataformat.InverterData, client mqtt.Client, config configuration) {
	log.Print("Publishing to mqtt")
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error converting to json", err)
	}
	token := client.Publish(config.MqttTopic, 0, false, jsonData)
	token.Wait()
}

func convertAndPost(data []byte, client mqtt.Client, config configuration) {
	converted, err := convertData(data)
	if err != nil {
		log.Print(err.Error())
	} else {
		postData(converted, client, config)
	}
}

func handleRequest(conn net.Conn, client mqtt.Client, config configuration) {
	buf := make([]byte, 512)
	length, err := conn.Read(buf)
	if err != nil {
		log.Print(err)
	}
	if length != dataformat.ExpectedLength {
		log.Printf("Invalid length for received packet: %d, expected %d, %s", length, dataformat.ExpectedLength, conn.RemoteAddr().String())
		if length > 0 {
			fmt.Printf("The invalid message: %+q\n", string(buf[0:length]))
		}
	} else {
		go convertAndPost(buf, client, config)
	}
	conn.Close()
}
