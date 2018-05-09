package main

import (
  "sync"
  "github.com/gilmae/fractal"
)

var counter = struct{
    sync.RWMutex
    visited map[fractal.Key]bool
}{visited: make(map[fractal.Key]bool)}

func find_escapee(points_map map[fractal.Key]fractal.Point) fractal.Key {
  for k, v := range points_map {
    if v.Escape < maxIterations {
      return k
    }
  }

  return fractal.Key{-1,-1}
}

func check_is_edge(k fractal.Key, edgePoints chan fractal.Point, points_map map[fractal.Key]fractal.Point) {
  if k.X < 0 || k.X >= width || k.Y < 0 || k.Y > height {
    return
  }

  if _, ok := counter.visited[k]; ok {
    return
  }
  counter.Lock()
  counter.visited[k] = true
  counter.Unlock()

  if (points_map[k].Escape >= maxIterations) {
    edgePoints <- points_map[k]
    return
  }

  for ii:= -1;ii<2;ii++ {
    for jj:= -1; jj<2; jj++ {
      if ii != 0 || jj != 00 {
        check_is_edge(fractal.Key{k.X+ii, k.Y+jj}, edgePoints, points_map)
      }
    }
  }
}

func find_edges(edgePoints chan fractal.Point, points_map map[fractal.Key]fractal.Point) {
  // scan for a non-escaped pixel
  var p = find_escapee(points_map)

  if (p.X >= 0) {
    check_is_edge(p, edgePoints, points_map)
  }
}
