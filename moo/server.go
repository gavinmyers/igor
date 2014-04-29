package moo

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type TelnetMooServer struct {
}

type MooServer interface {
	Init()
}

type ActiveClient struct {
	Name      string
	IN        chan string
	OUT       chan string
	Con       net.Conn
	Quit      chan bool
	ListChain *list.List
}

func (c *ActiveClient) Read(buf []byte) bool {
	_, err := c.Con.Read(buf)
	if err != nil {
		c.Close()
		return false
	}
	return true
}

func (c *ActiveClient) receive() {
	buf := make([]byte, 2048)
	for c.Read(buf) {
    fmt.Printf("\nv\n",string(buf))
		if bytes.Equal(buf, []byte("quit")) {
			c.Close()
			break
		}
		var r Action
		json.Unmarshal([]byte(strings.Trim(string(buf), "\x00")), &r)
    r.Target = "UPDATED"
		send, _ := json.Marshal(r)
		fmt.Printf("\n%v ->]\n ", r)
		c.OUT <- string(send)
	}
	r := &Action{Name: c.Name, Target: "LEFT"}
	send, _ := json.Marshal(r)
	c.OUT <- string(send)
}

func (client *ActiveClient) broadcast() {
	for {
		select {
		case buf := <-client.IN:
			client.Con.Write([]byte(buf))
		case <-client.Quit:
			client.Con.Close()
			break
		}
	}
}

func (c *ActiveClient) Close() {
	c.Quit <- true
	c.Con.Close()
	for e := c.ListChain.Front(); e != nil; e = e.Next() {
		c2 := e.Value.(ActiveClient)
    if bytes.Equal([]byte(c.Name), []byte(c2.Name)) {
      if c.Con == c2.Con {
        c.ListChain.Remove(e)
      }
		}
	}
}


func (c *TelnetMooServer) receive(IN <-chan string, lst *list.List) {
	for {
		input := <-IN
		for val := lst.Front(); val != nil; val = val.Next() {
			client := val.Value.(ActiveClient)
			client.IN <- input
		}
	}
}

func (c *TelnetMooServer) broadcast(con net.Conn, ch chan string, lst *list.List) {
	buf := make([]byte, 1024)
	bytenum, _ := con.Read(buf)
	name := string(buf[0:bytenum])
	newclient := &ActiveClient{name, make(chan string), ch, con, make(chan bool), lst}
	go newclient.broadcast()
	go newclient.receive()
	lst.PushBack(*newclient)
	r := &Action{Name: name, Target: "JOINED"}
	send, _ := json.Marshal(r)
	fmt.Printf("\n%v\n",string(send))
	ch <- string(send)
}

func (c *TelnetMooServer) Init() {
	clients := list.New()
	in := make(chan string)
	go c.receive(in, clients)
	netlisten, _ := net.Listen("tcp", "127.0.0.1:9988")
	defer netlisten.Close()
	for {
		conn, _ := netlisten.Accept()
		go c.broadcast(conn, in, clients) //&conn..
	}
}
