package main

import ("./moo"
  "flag"
 )

var gui moo.GUI
var client moo.MooClient

func main() {
  useQml := flag.Bool("gui", true, "use the graphical interface")
  client = &moo.TelnetMooClient {}
  client.Init()
  if(*useQml == true) {
    gui = &moo.QmlGUI {}
  } else {
    gui = &moo.TermboxGUI {}
  }
  gui.Main()
}
