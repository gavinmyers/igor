package igor
import "math/rand"

func Split(w,h int) []*Tile {
  tiles :=  make([]*Tile, 25)
  rwidth := rand.Intn(w/2) + w/4

  rwidth_c1 := rand.Intn(rwidth/2) + rwidth/4
  rwidth_c2 := rand.Intn((w-rwidth)/2) +
    rwidth + ((w - rwidth) / 4)

  rheight := rand.Intn(h/2) + h/4

  rheight_c1 := rand.Intn(rheight/2) + rheight/4
  rheight_c2 := rand.Intn((h-rheight)/2) +
    rheight + ((h - rheight) / 4)

  tiles[0] = &Tile {X:0, Y:rheight_c1}
  tiles[1] = &Tile {X:rwidth_c1, Y:rheight_c1}
  tiles[2] = &Tile {X:rwidth, Y:rheight_c1}
  tiles[3] = &Tile {X:rwidth_c2, Y:rheight_c1}
  tiles[4] = &Tile {X:w, Y:rheight_c1}

  tiles[5] = &Tile {X:0, Y:rheight}
  tiles[6] = &Tile {X:rwidth_c1, Y:rheight}
  tiles[7] = &Tile {X:rwidth, Y:rheight}
  tiles[8] = &Tile {X:rwidth_c2, Y:rheight}
  tiles[9] = &Tile {X:w, Y:rheight}

  tiles[10] = &Tile {X:0, Y:rheight_c2}
  tiles[11] = &Tile {X:rwidth_c1, Y:rheight_c2}
  tiles[12] = &Tile {X:rwidth, Y:rheight_c2}
  tiles[13] = &Tile {X:rwidth_c2, Y:rheight_c2}
  tiles[14] = &Tile {X:w, Y:rheight_c2}

  tiles[15] = &Tile {X:0, Y:0}
  tiles[16] = &Tile {X:rwidth_c1, Y:0}
  tiles[17] = &Tile {X:rwidth, Y:0}
  tiles[18] = &Tile {X:rwidth_c2, Y:0}
  tiles[19] = &Tile {X:w, Y:0}

  tiles[20] = &Tile {X:0, Y:h}
  tiles[21] = &Tile {X:rwidth_c1, Y:h}
  tiles[22] = &Tile {X:rwidth, Y:h}
  tiles[23] = &Tile {X:rwidth_c2, Y:h}
  tiles[24] = &Tile {X:w, Y:h}



  return tiles
}

func Clear(w,h int) *Map {
  m := &Map {Tiles:make([][]*Tile, w * 2)}
  for i := 0; i < w * 2; i++ {
    m.Tiles[i] = make([]*Tile, h * 2)
    for j := 0; j < h * 2; j++ {
      t := &Tile {X:i, Y:j}
      m.Tiles[i][j] = t
    }
  }

  return m
}

func Connect(t1,t2 *Tile, m *Map) {
  startX := t1.X
  startY := t1.Y
  endX := t2.X
  endY := t2.Y
	lenX := startX - endX
	if lenX < 0 {
		lenX = endX - startX
	}
	lenY := startY - endY
	if lenY < 0 {
		lenY = endY - startY
	}
	lenT := lenX + lenY
	path := make([]int, lenT, lenT)
	for i := 0; i < lenT; i++ {
    if(i < lenX) {
      if(startX > startY) {
        path[i] = DirWest
      } else {
        path[i] = DirEast
      }
    } else {
      if(startY > endY) {
        path[i] = DirNorth
      } else {
        path[i] = DirSouth
      }
    }
	}
  perm := rand.Perm(len(path))
  tiles := make([]*Tile, lenT)
  pathX := startX
  pathY := startY
  for i, v := range perm {
    if path[v] == DirWest {
      pathX--
    } else if path[v] == DirEast {
      pathX++
    } else if path[v] == DirNorth {
      pathY--
    } else if path[v] == DirSouth {
      pathY++
    }
    t := &Tile {X:pathX, Y:pathY}
    tiles[i] = t
    if(pathX > 0 && pathY > 0) {
      m.Tiles[pathX][pathY] = t
    }
  }
}

