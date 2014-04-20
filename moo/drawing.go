package moo


type GUI interface {
  Width() int
  Height() int
  Draw(int, int, int, int, string, ...interface{})
  Flush()
  Init()
  Main()
}


