package moo

type Sprite struct {
  Id int
  Code string
}
type SpriteBuilder struct {
  Sprites map[string]*Sprite
}
func (c *SpriteBuilder) Init() {
  c.Sprites = make(map[string]*Sprite)
  c.Sprites["WALL"] =       &Sprite {Id:100, Code:"#"}
  c.Sprites["VOID"] =       &Sprite {Id:101, Code:"."}
  c.Sprites["PLAYER"] =     &Sprite {Id:102, Code:"@"}
}
