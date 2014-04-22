package moo

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type Action struct {
	Name   string
	Target string
}

type TelnetMooClient struct {
	running    bool
	connection net.Conn
	server     string
	port       string
}

type MooClient interface {
	Read(con net.Conn) string
	Send(msg []byte)
	Receive(chan<- *Action)
	Init()
}


func (c *TelnetMooClient) Read(con net.Conn) string {
	var buf [4048]byte
	_, err := con.Read(buf[0:4048])
	if err != nil {
		con.Close()
		c.running = false
		return "Error in reading!"
	}
	str := string(buf[0:4048])
	return string(str)
}

func (c *TelnetMooClient) Send(msg []byte) {
	c.connection.Write(msg)
	/*
	   reader := bufio.NewReader(os.Stdin);
	   for {
	       input, err := reader.ReadBytes('\n')
	       if err == nil  && len(input) > 1 {
	           tokens := strings.Fields(string(input[0:len(input)-1]))

	           if tokens[0] == "/quit" {
	               c.connection.Write([]byte("is leaving..."))
	               c.running = false
	               break
	           } else if tokens[0] == "/command" {
	               if len(tokens) > 1 {
	                   out, err := exec.Command(tokens[1], tokens[2:]...).Output()
	                   if err != nil {
	                       fmt.Printf("Error: %s\n", err)
	                   } else {
	                       c.connection.Write(out) // send output to server
	                   }
	               } else {
	                   fmt.Printf("Usage:\n\t/command <exec> <arguments>\n\tEx: /command ls -l -a\n\n")
	               }
	               continue
	           }
	           c.connection.Write(input[0:len(input)-1])
	       }
	   }
	*/
}

func (c *TelnetMooClient) Receive(out chan<- *Action) {
	for {
		buf := make([]byte, 2048)
		_, err := c.connection.Read(buf)
		if err != nil {
			panic(err)
		}
		var rec Action
    //TODO: Convert buf from byte to string back to byte... hrm...
		json.Unmarshal([]byte(strings.Trim(string(buf), "\x00")), &rec)
    fmt.Printf("\n%v ->]\n ", rec.Target)
    out <- &rec

	}
}

func (c *TelnetMooClient) Init() {
	c.running = true
	destination := fmt.Sprintf("%s:%s", "127.0.0.1", "9988")
	cn, _ := net.Dial("tcp", destination)
	c.connection = cn
	//  defer cn.Close();
}
