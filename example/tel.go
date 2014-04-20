/* **********N-CLIENT CHAT SERVER ***********************
   ********** run program in one  terminal********
   ******* in another one run :- telnet localhost 9988********/
package main

import (
  "net"
  "strings"
  "container/list"
  "bytes"
  "encoding/json"
       )

type ClientChat struct {
  Name string
  IN chan string
  OUT chan string
  Con net.Conn
  Quit chan bool
  ListChain *list.List
}

type Action struct {
  Name string
  Target string
}

func (c *ClientChat) Read(buf []byte) bool {
  _,err := c.Con.Read(buf)
  if err != nil {
    c.Close()
    return false
  }
  return true
}

func (c *ClientChat) Close() {
  c.Quit <- true
  c.Con.Close()
  c.deleteFromList()
}

func (c *ClientChat) Equal(cl *ClientChat) bool {
  if bytes.Equal([]byte(c.Name),[]byte(cl.Name)) {
    if c.Con == cl.Con {
      return true
    }
  }
  return false
}

func (c *ClientChat) deleteFromList() {
  for e:= c.ListChain.Front(); e != nil ; e = e.Next() {
    client := e.Value.(ClientChat)
    if c.Equal(&client) {
    c.ListChain.Remove(e)
    }
  }
}

func handlingINOUT(IN <-chan string,lst *list.List) {
  for {
    input := <-IN
    for val := lst.Front();val != nil;val = val.Next() {
      client := val.Value.(ClientChat)
      client.IN <- input
    }
  }
}

func clientreceiver(client *ClientChat) {
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

func clientsender(client *ClientChat) {
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

func clientHandling(con net.Conn,ch chan string,lst *list.List) {
  buf := make([]byte,1024)
  bytenum,_ := con.Read(buf)
  name := string(buf[0:bytenum])
  newclient := &ClientChat{name,make(chan string),ch,con,make(chan bool),lst}
  go clientsender(newclient)
  go clientreceiver(newclient)
  lst.PushBack(*newclient)
  r := &Action{Name:name, Target:"JOINED"}
  send,_ := json.Marshal(r)
  ch <- string(send)
}

func main() {
  clientlist := list.New()
  in := make(chan string)
  go handlingINOUT(in,clientlist)
  netlisten,_ := net.Listen("tcp","127.0.0.1:9988")
  defer netlisten.Close()
  for {
    conn,_ := netlisten.Accept()
    go clientHandling(conn,in,clientlist)//&conn..
  }
}
