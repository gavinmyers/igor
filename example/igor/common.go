package igor

var WinWidth = 0
var WinHeight = 0
var DirNorth = 8
var DirSouth = 2
var DirWest = 4
var DirEast = 6
var DirNorthWest = 7
var DirNorthEast = 9
var DirSouthWest = 1
var DirSouthEast = 3

type Tile struct {
  Z int
  X int
  Y int
}


type Map struct {
  Z int
  Tiles [][]*Tile
}

