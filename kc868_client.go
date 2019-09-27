package main

import (
	"bufio"
	"fmt"
	"github.com/amaxlab/go-lib/log"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const PingMinutesTime = 1

type Relay struct {
	Value   bool      `json:"value"`
	OnTime  time.Time `json:"on_time"`
	OffTime time.Time `json:"off_time"`
}

type KC868Client struct {
	Host      string
	Port      int
	Relays    map[string]*Relay
	Connect   net.Conn
	Connected bool
}

func NewKC868Client(Host string, Port int) *KC868Client {
	return &KC868Client{Host: Host, Port: Port, Relays: make(map[string]*Relay), Connected: false}
}

func (c *KC868Client) connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		log.Error.Printf("%s", err)
		return
	}

	c.Connect = conn
	c.Connected = true
	go c.reader()
	go c.pinger()
	c.send("RELAY-SCAN_DEVICE-NOW")
}

func (c *KC868Client) disconnect() {
	_ = c.Connect.Close()
}

func (c *KC868Client) pinger() {
	for {
		time.Sleep(time.Minute * PingMinutesTime)
		c.send("PING")
	}
}

func (c *KC868Client) reader() {
	buf := bufio.NewReader(c.Connect)

	for {
		str, err := buf.ReadString(0)
		if len(str) > 0 {
			c.handle(str)
		}
		if err != nil {
			if err == io.EOF {

			}
			break
		}
	}
}

func (c *KC868Client) send(text string) {
	log.Debug.Printf("Send command -> %s", text)
	if !c.Connected {
		log.Warning.Printf("Client not connected")
		return
	}

	_, err := c.Connect.Write([]byte(text))
	if err != nil {
		log.Error.Printf("%s", err)
	}
}

func (c *KC868Client) handle(text string) {
	log.Debug.Printf("Response from server -> %s", text)
	response := strings.Split(text, "-")
	if len(response) < 2 || response[0] != "RELAY" {
		log.Warning.Printf("Wrong response format from server -> %s", response[0])
		return
	}

	switch response[1] {
	case "SCAN_DEVICE":
		dev := strings.Split(strings.Split(response[2], "_")[1], ",")
		count, _ := strconv.Atoi(dev[0])
		go c.StartScan(count)
	case "READ":
		if len(response) < 3 {
			log.Error.Printf("Wrong READ command parameters -> %s", response)
			return
		}
		r := strings.Split(response[2], ",")
		c.setRelayState(r[1], r[2])
	case "SET":
		if len(response) < 3 {
			log.Error.Printf("Wrong SET command parameters -> %s", response)
			return
		}
		r := strings.Split(response[2], ",")
		c.setRelayState(r[1], r[2])
	default:
		log.Warning.Printf("Wrong response command -> %s", response)
	}
}

func (c *KC868Client) ChangeRelayState(id string, newState bool) {
	state := "0"
	if newState {
		state = "1"
	}
	c.send(fmt.Sprintf("RELAY-SET-1,%s,%s", id, state))
}

func (c *KC868Client) StartScan(count int) {
	for i := 1; i <= count; i++ {
		c.send(fmt.Sprintf("RELAY-READ-1,%d", i))
		time.Sleep(time.Millisecond * 500)
	}
}

func (c *KC868Client) setRelayState(id, value string) {
	if c.Relays[id] == nil {
		c.Relays[id] = &Relay{}
	}

	if value == "1" {
		c.Relays[id].Value = true
		c.Relays[id].OnTime = time.Now()
	} else {
		c.Relays[id].Value = false
		c.Relays[id].OffTime = time.Now()
	}
}
