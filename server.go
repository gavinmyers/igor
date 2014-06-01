package main

import (
  "fmt"
  "net"
  "strings"
  "io"
	"encoding/json"
)

var clients map[string]*Client
var listener net.Listener

type Client struct {
  Data map[string]string
  Conn net.Conn
}

type Message struct {
  Token string
  Action string
  Data map[string]string
}

func Read(con net.Conn) {
  buf := make([]byte, 2048)
  for {
    bytenum, err := con.Read(buf)
    strin := string(buf[0:bytenum])
    if err == io.EOF || err != nil || len(strin) == 0 {
      con.Close()
      return
    }
    arr := strings.Split(strin, "\r\n")
    //fmt.Printf("\n Svr Received %v",strin) 
    for _, msg := range arr {
      if(len(msg) > 0) {
        var message Message
        json.Unmarshal([]byte(msg), &message)
        //fmt.Printf("\n Svr Received %d > %v %v %v", bytenum, message.Token, message.Action, message.Data)
        switch message.Action {
          case "ack":
            send, _ := json.Marshal(&Message{Token:message.Token, Action:"ack",Data:map[string]string{"ack":"ack"}})
            fmt.Printf("\n ACK [%v] ", message.Token)
            go WriteTo(message.Token,message.Token,string(send) + "/r/n")
          case "connect":
            client := &Client{Conn:con,Data:map[string]string{"session":message.Token}}
            clients[message.Token] = client
            send, _ := json.Marshal(&Message{Token:message.Token, Action:"ack",Data:map[string]string{"key":"value"}})
            fmt.Printf("\n SYSTEM ONLINE [%v] ", message.Token)
            go WriteTo(message.Token,message.Token,string(send) + "/r/n")
          case "disconnect":
            fmt.Printf("\n GOODBYE [%v] ", message.Token)
            clients[message.Token].Conn.Close()
            delete(clients, message.Token)
          case "writeto":
            go WriteTo(message.Token, message.Data["target"],msg)
          case "who":
            go Who(message.Token)
          default:
            go WriteToAll(msg)
        }
      }
    }
  }
}
func Who(from string) {
  client := clients[from]
  client.Conn.Write([]byte("LIST" + "\r\n"))
}

func WriteTo(from string, to string, msg string) {
  client := clients[to]
  if(client != nil) {
    client.Conn.Write([]byte(msg + "\r\n"))
  }
}

func WriteToAll(msg string) {
  for _, client := range clients {
    client.Conn.Write([]byte(msg + "\r\n"))
  }
}

func Listen() {
  net, _ := net.Listen("tcp", "127.0.0.1:9988")
  clients = make(map[string]*Client)
  listener = net
  for {
    con, _ := listener.Accept()
    go Read(con)
  }
}

func main() {
  Listen()
}

