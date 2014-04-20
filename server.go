package main

import ("./moo")

var server moo.MooServer
func main() {
  server = &moo.TelnetMooServer {}
  server.Init()
}
