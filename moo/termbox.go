package moo

import "github.com/limetext/termbox-go"
import "fmt"

type TermboxGUI struct {
  color [256]termbox.Attribute
  width int
  height int
}

func (c *TermboxGUI) Main() {
	defer termbox.Close()
  c.Init()
loop:
	for {
    c.Draw(25,25,125,11,"$")
    c.Flush()
    key := c.PollEvent()
    if key == termbox.KeyCtrlC {
      break loop
    }
  }
}


func (c *TermboxGUI) Width() int {
  return c.width
}

func (c *TermboxGUI) Height() int {
  return c.height
}

func (c *TermboxGUI) Init() {
  err := termbox.Init()
	termbox.SetColorMode(termbox.ColorMode256)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 255; i++ {
		c.color[i] = termbox.Attribute(i)
	}
	termbox.SetInputMode(termbox.InputEsc)
  termbox.Clear(c.color[0], c.color[0])
  c.width, c.height = termbox.Size()
}

func (c *TermboxGUI) Flush() {
  termbox.Flush()
}

func (c *TermboxGUI) PollEvent() termbox.Key {
  ev := termbox.PollEvent()
  return ev.Key
}

func (c TermboxGUI) Draw(x, y, fg, bg int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	c.print_tb(x, y, c.color[fg], c.color[bg], s)
}

func (c *TermboxGUI) print_tb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}


