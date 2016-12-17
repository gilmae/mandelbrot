package main

import (
  "sync"
)

var counter = struct{
    sync.RWMutex
    visited map[Key]bool
}{visited: make(map[Key]bool)}

func find_escapee() Key {
  for k, v := range points_map {
    if v.Escape < maxIterations {
      return k
    }
  }

  return Key{-1,-1}
}

func check_is_edge(k Key, edgePoints chan Point) {
  if k.x < 0 || k.x >= width || k.y < 0 || k.y > height {
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
        check_is_edge(Key{k.x+ii, k.y+jj}, edgePoints)
      }
    }
  }
}

func find_edges(edgePoints chan Point) {
  // scan for a non-escaped pixel
  var p = find_escapee()

  if (p.x >= 0) {
    check_is_edge(p, edgePoints)
  }
}
