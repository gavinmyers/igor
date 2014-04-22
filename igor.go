package main

import (
	"./moo"
	"flag"
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
  client.Send([]byte("bleagh"))
	gui.Init()
}
