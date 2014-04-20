package main

import (
	"fmt"
	"gopkg.in/qml.v0"
	"os"
)

var Root qml.Object

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type Player struct {
	qml.Object
}
type Game struct {
	qml.Object
}
type Floor struct {
	qml.Object
}
func (r *Player) HandleClick(xPos, yPos int) {
	r.Set("x", r.Int("x")+40)
	r.Set("y", r.Int("y")+40)
}

func (r *Player) Update() {
	r.Set("x", r.Int("x")-2)
	r.Set("y", r.Int("y")-2)
}
func (r *Game) Update() {
}
func (r *Game) Init() {
  t := Root.Object("floor_0_1")
  for x := -24; x < 1260; x=x+16 {
    for y := -24; y < 960; y=y+16 {
      c := t.Create(nil)
      c.Set("x",x)
      c.Set("y",y)
      c.Set("enabled", true)
      c.Set("parent",Root)
    }
  }
}
func (r *Player) Paint(p *qml.Painter) {
}

func (r *Game) Paint(p *qml.Painter) {
}

func (r *Floor) Paint(p *qml.Painter) {
}

func run() error {
	qml.Init(nil)

	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{
		{Init: func(g *Game, obj qml.Object) { g.Object = obj }},
		{Init: func(g *Player, obj qml.Object) { g.Object = obj }},
		{Init: func(g *Floor, obj qml.Object) { g.Object = obj }},
	})

	engine := qml.NewEngine()
	component, err := engine.LoadFile("painting.qml")
	if err != nil {
		return err
	}

	win := component.CreateWindow(nil)
  Root = win.Root()
	win.Show()
	win.Wait()

	return nil
}
