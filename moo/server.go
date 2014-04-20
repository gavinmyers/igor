package moo

import (
  "net"
  "strings"
  "container/list"
  "bytes"
  "encoding/json")

type TelnetMooServer struct {
}

type MooServer interface {
  Init()
}

type ActiveClient struct {
  Name string
  IN chan string
  OUT chan string
  Con net.Conn
  Quit chan bool
  ListChain *list.List
}

func (c *ActiveClient) Read(buf []byte) bool {
  _,err := c.Con.Read(buf)
  if err != nil {
    c.Close()
    return false
  }
  return true
}

func (c *ActiveClient) Close() {
  c.Quit <- true
  c.Con.Close()
  c.deleteFromList()
}

func (c *ActiveClient) Equal(cl *ActiveClient) bool {
  if bytes.Equal([]byte(c.Name),[]byte(cl.Name)) {
    if c.Con == cl.Con {
      return true
    }
  }
  return false
}

func (c *ActiveClient) deleteFromList() {
  for e:= c.ListChain.Front(); e != nil ; e = e.Next() {
    client := e.Value.(ActiveClient)
    if c.Equal(&client) {
    c.ListChain.Remove(e)
    }
  }
}

func (c *TelnetMooServer) handlingINOUT(IN <-chan string,lst *list.List) {
  for {
    input := <-IN
    for val := lst.Front();val != nil;val = val.Next() {
      client := val.Value.(ActiveClient)
      client.IN <- input
    }
  }
}

func (c *TelnetMooServer) clientreceiver(client *ActiveClient) {
  buf := make([]byte,2048)
  for client.Read(buf) {
    if bytes.Equal(buf,[]byte("quit")) {
      client.Close()
      break
    }
    r := &Action{Name:client.Name, Target:strings.TrimRight(string(buf), "\x00")}
    send,_ := json.Marshal(r)
    client.OUT <-string(send)
  }
  r := &Action{Name:client.Name, Target:"LEFT"}
  send,_ := json.Marshal(r)
  client.OUT <-string(send)
}

func (c *TelnetMooServer) clientsender(client *ActiveClient) {
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

func (c *TelnetMooServer) clientHandling(con net.Conn,ch chan string,lst *list.List) {
  buf := make([]byte,1024)
  bytenum,_ := con.Read(buf)
  name := string(buf[0:bytenum])
  newclient := &ActiveClient{name,make(chan string),ch,con,make(chan bool),lst}
  go c.clientsender(newclient)
  go c.clientreceiver(newclient)
  lst.PushBack(*newclient)
  r := &Action{Name:name, Target:"JOINED"}
  send,_ := json.Marshal(r)
  ch <- string(send)
}

func (c *TelnetMooServer) Init() {
  clientlist := list.New()
  in := make(chan string)
  go c.handlingINOUT(in,clientlist)
  netlisten,_ := net.Listen("tcp","127.0.0.1:9988")
  defer netlisten.Close()
  for {
    conn,_ := netlisten.Accept()
    go c.clientHandling(conn,in,clientlist)//&conn..
  }
}
