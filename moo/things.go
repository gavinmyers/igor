package moo

//#0
type Thing struct {
  Id int
  Parent *Thing
  World *World
  Area *Area
  Map *Map
  Tile *Tile
  Children []*Thing
  Sprite Sprite
}
