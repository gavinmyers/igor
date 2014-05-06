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
  Session string
  Conn net.Conn
}

type Message struct {
  Session string
  Action string
  Commands map[string]string
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
    for _, msg := range arr {
      if(len(msg) > 0) {
        var message Message
        json.Unmarshal([]byte(msg), &message)
        fmt.Printf("\n Svr Received %d > %v %v %v", bytenum, message.Session, message.Action, message.Commands)
        switch message.Action {
          case "CONNECT":
            fmt.Printf("\n WELCOME TO IGOR [%v] ", message.Session)
            client := &Client{Session:message.Session,Conn:con}
            clients[message.Session] = client
          case "DISCONNECT":
            fmt.Printf("\n GOODBYE [%v] ", message.Session)
            clients[message.Session].Conn.Close()
            delete(clients, message.Session)
          default:
            go Write(msg)
        }
      }
    }
  }
}

func Write(msg string) {
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

func clientListener(con net.Conn,who string) {
  buf := make([]byte, 2048)
  for {
    bytenum, err := con.Read(buf)
    strin := string(buf[0:bytenum])
    if err == io.EOF || err != nil || len(strin) == 0 {
      con.Close()
      return
    }
    arr := strings.Split(strin, "\r\n")
    for _, msg := range arr {
      if(len(msg) > 0) {
        var message Message
        json.Unmarshal([]byte(msg), &message)
        if(message.Session != who) {
          fmt.Printf("\n %v from %v to %v '%v'",message.Action, message.Session,who, message.Commands["5"])
        }
      }
    }
  }
}

func clientBroadcast(con net.Conn, who string, msg string) {
  send, _ := json.Marshal(&Message{Session:who, Action:"SPEAK", Commands:map[string]string{"5":"five", "55":"z", "555":"b"}})
  con.Write([]byte(string(send)+"\r\n"))
}

func clientJoin(con net.Conn, who string) {
  send, _ := json.Marshal(&Message{Session:who, Action:"CONNECT"})
  con.Write([]byte(string(send)+"\r\n"))
}
func clientDisconnect(con net.Conn, who string) {
  send, _ := json.Marshal(&Message{Session:who, Action:"DISCONNECT"})
  con.Write([]byte(string(send)+"\r\n"))
}

func client(who string) {
  con, _ := net.Dial("tcp", "127.0.0.1:9988")
  go clientListener(con,who)
  clientJoin(con,who)
  clientBroadcast(con,who,"Hello")
  clientBroadcast(con,who,"How are you today?")
  clientBroadcast(con,who,"Goodbye")
  //clientDisconnect(con,who)
}

func main() {
  go Listen()
  go client("Joe")
  go client("Frank")
  go client("Lisa")
  var i string
  fmt.Scanf("%v", &i)
}
