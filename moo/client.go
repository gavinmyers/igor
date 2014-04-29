package moo

/*
what should the server do

#1 - Receive messages from client
server.Receive(Message)
#2 - Broadcast messages to all clients... really this should be a subset depending on need
server.Broadcast(Message, []clients)
#3 - Connect a new client
server.Connect(client)
#4 - Disconnect a client
server.Disconnect(client)

at some point in time the server will actually need to _do_ _something_ ... which means take client message X and apply it to logic engine Y... how in the hell do you want to do that.

Moo world

X <interaction> <against> Y

player.doAttack(lawnmower)
"you can't attack a lawnmower"

client.send("Player:Attack:Lawnmower")
server.receive(above)

Engine.get(Player).do(Player.Attack).to(Lawnmower)



*/


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
	Send(msg *Action)
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

func (c *TelnetMooClient) Send(act *Action) {
  send, _ := json.Marshal(act)
  fmt.Printf("\nSending: %s\n", send)
	c.connection.Write(send)
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
