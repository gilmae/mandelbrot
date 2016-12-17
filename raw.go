package main

import (
  "encoding/json"
  "os"
  "fmt"
)

func (p *Point) MarshalJSON() ([]byte, error) {
	type Alias Point
	return json.Marshal(&struct {
    X int
    Y int
    Escape float64
    Real float64
    Imaginary float64
	}{
    X : p.X,
    Y : p.Y,
    Escape : p.Escape,
    Real : real(p.C),
    Imaginary : imag(p.C),
	})
}

func write_raw(points map[Key]Point, filename string) {
  file, err := os.Create(filename)
  if err != nil {
    fmt.Println(err)
  }


  for _, v := range points_map {
    json.NewEncoder(file).Encode(&v)
  }

  file.Close()
}
