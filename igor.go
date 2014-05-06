package main

import (
	"./moo"
	"flag"
  "time"
  "fmt"
)

var gui moo.GUI
var client moo.MooClient

func main() {
	useQml := flag.Bool("gui", true, "use the graphical interface")
	client = &moo.TelnetMooClient{}
	if *useQml == true {
		gui = &moo.QmlGUI{}
	} else {
		gui = &moo.TermboxGUI{}
	}
	act := make(chan *moo.Action)
	client.Init()
	go client.Receive(act)
	go gui.Receive(act)
  go send()
	gui.Init()
  time.Sleep(10 * time.Millisecond)
}

func send() {
  fmt.Printf("\nGeneric Send\n")
  r1 := &moo.Action{Name: "Joe", Target: "LOOK"}
	go client.Send(r1)
  r2 := &moo.Action{Name: "Joe", Target: "SIT"}
	go client.Send(r2)
  r3 := &moo.Action{Name: "Joe", Target: "STAND"}
	go client.Send(r3)
}
