package main

import (
  "fmt"
  "net"
  "strings"
  "time"
  "io"
	"encoding/json"

  "math/rand"
  //"math/big"
)

func randString(n int) string {
    r := rand.New(rand.NewSource(99))
    var bytes = make([]byte, n)
    for i := range bytes {
        bytes[i] = byte(r.Intn(10))
    }
    return string(bytes)
}
var r string

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
          case "connect":
            client := &Client{Conn:con,Data:map[string]string{"session":message.Token}}
            clients[message.Token] = client
            send, _ := json.Marshal(&Message{Token:message.Token, Action:"ack",Data:map[string]string{"rand":randString(30000)}})
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
        if(message.Token != who) {
          fmt.Printf("\n %v from %v to %v '%v'",message.Action, message.Token,who, message.Data["message"])
        }
      }
    }
  }
}

func clientBroadcast(con net.Conn, who string, msg string) {
  message := &Message{Token:who, Action:"writeto",
                      Data:map[string]string{"target":"Rbt1", "message":msg}}
  send, _ := json.Marshal(message)
  con.Write([]byte(string(send)+"\r\n"))
}

func clientJoin(con net.Conn, who string) {
  send, _ := json.Marshal(&Message{Token:who, Action:"connect"})
  con.Write([]byte(string(send)+"\r\n"))
}
func clientDisconnect(con net.Conn, who string) {
  send, _ := json.Marshal(&Message{Token:who, Action:"disconnect"})
  con.Write([]byte(string(send)+"\r\n"))
}

func client(who string) {
  con, _ := net.Dial("tcp", "127.0.0.1:9988")
  go clientListener(con,who)
  clientJoin(con,who)
  time.Sleep(100 * time.Millisecond)
  clientBroadcast(con,who,"Service test1")
  time.Sleep(50 * time.Millisecond)
  clientBroadcast(con,who,"Service test2")
  time.Sleep(50 * time.Millisecond)
  clientBroadcast(con,who,"Goodbye")
  time.Sleep(100 * time.Millisecond)
  clientDisconnect(con,who)
}

func main() {
  r = randString(1200000)
  go Listen()
  go client("Rbt1")
  go client("Rbt2")
  go client("Rbt3")
  var i string
  fmt.Scanf("%v", &i)
}
