# Ginlong tools

This folder contains tools for analyzing inverter data packets from Ginlong inverters (mostly sold under the Solis brand).

## Usage

To receive packets the inverter must be set up to send packets to the ip address you're going to run these tools on. The inverters sends updates about every 6 minutes. It happens more often in the morning when it's powering on.

## Dataformat

The data format to deserialize packets from the inverter. Contains a helper function to transform the packet into a more useful representation.

## Ginlongmonitor

Monitors for packets and writes them to a file with the current timestamp as name. Useful for debugging problems with the parsing.

## Ginlongparse

Usage: `ginlongparse -f filename.bin`. Reads in a binary packet (such as received with ginlongmonitor), parses it and prints some information contained in the packet. Mostly useful for debugging, possibly this can be adapted for your use case.

## Ginlongmqtt

This service posts data received from the inverter to mqtt as json.

Some environment variables must be set to make this work.


Required

- MQTT_USERNAME: the username for the mqtt server
- MQTT_PASSWORD: the password for the mqtt server

Optional

- INVERTER_LISTENPORT: on which port to listen for connections (default 9999)
- MQTT_CLIENTID: clientid for use on mqtt server (default TODO)
- MQTT_SERVERADDRESS: server address for mqtt (default 127.0.0.1)
- MQTT_SERVERPORT: server port for mqtt server (default 1883)
- MQTT_INVERTER_TOPIC: topic on which to post data (default TODO)

## Alternatives

An alternative software that does some of these things too is [ginlong-wifi](https://github.com/graham0/ginlong-wifi).

