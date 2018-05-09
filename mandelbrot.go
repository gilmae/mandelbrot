package main

import (
    "fmt"
    "math/cmplx"
    "strconv"
    "math"
    "math/rand"
    "time"
    "flag"
    "github.com/gilmae/fractal"
    "github.com/gilmae/rescale"
)

var  maxIterations float64 = 1000.0
var  bailout float64 = 4.0
var  width int = 1600
var height int = 1600
var x int = 0
var y int = 0

var default_gradient string = `[["0.0", "000764"],["0.16", "026bcb"],["0.42", "edffff"],["0.6425", "ffaa00"],["0.8675", "000200"],["1.0","000764"]]`

const (
    rMin   = -2.25
    rMax   = 0.75
    iMin   = -1.5
    iMax   = 1.5
    usage  = "mandelbot OPTIONS\n\nPlots the mandelbrot set, centered at a point indicated by the provided real and imaginary, and at the given zoom level.\n\nSaves the output into the given path.\n\n"
)

func get_cordinates(midX float64, midY float64, zoom float64, width int, height int, x int, y int) complex128 {
  new_r_start, new_r_end := rescale.GetZoomedBounds(rMin, rMax, midX, zoom)
  scaled_r := rescale.Rescale(new_r_start, new_r_end, width, x)

  new_i_start, new_i_end := rescale.GetZoomedBounds(iMin, iMax, midY, zoom)
  scaled_i := rescale.Rescale(new_i_start, new_i_end, height, height-y)

  return complex(scaled_r, scaled_i)
}


func main() {
  //start := time.Now()

  var midX float64
  var midY float64
  var zoom float64
  var output string
  var filename string
  var gradient string
  var mode string

  rand.Seed(time.Now().UnixNano())
  flag.Float64Var(&midX, "r", -0.75, "Real component of the midpoint.")
  flag.Float64Var(&midY, "i", 0.0, "Imaginary component of the midpoint.")
  flag.Float64Var(&zoom, "z", 1, "Zoom level.")
  flag.StringVar(&output, "o", ".", "Output path.")
  flag.StringVar(&filename, "f", "", "Output file name.")
  flag.StringVar(&colour_mode, "c", "none", "Colour mode: true, smooth, banded, none.")
  flag.Float64Var(&bailout, "b", 4.0, "Bailout value.")
  flag.IntVar(&width, "w", 1600, "Width of render.")
  flag.IntVar(&height, "h", 1600, "Height of render.")
  flag.Float64Var(&maxIterations, "m", 2000.0, "Maximum Iterations before giving up on finding an escape.")
  flag.StringVar(&gradient, "g", default_gradient, "Gradient to use.")
  flag.StringVar(&mode, "mode", "image", "Mode: edge, image, coordsAt")
  flag.IntVar(&x, "x", 0, "x cordinate of a pixel, used for translating to the real component. 0,0 is top left.")
  flag.IntVar(&y, "y", 0, "y cordinate of a pixel, used for translating to the real component. 0,0 is top left.")
  flag.Parse()
  
  var calculator fractal.EscapeCalculator = func(z complex128) (float64, complex128, bool) {
    iteration := 0.0
    c := z

    for z= c;cmplx.Abs(z) < bailout && iteration < maxIterations; iteration+=1 {
      z = z*z+c;
    }

    if (iteration >= maxIterations) {
      return maxIterations, z, false
    }
    
    if (mode=="image" && colour_mode=="smooth") {
      z = z*z+c
      z = z*z+c
      iteration += 2
      reZ := real(z)
      imZ := imag(z)
      magnitude := math.Sqrt(reZ * reZ + imZ * imZ)
      mu := iteration + 1 - (math.Log(math.Log(magnitude)))/math.Log(2.0)
    
      return mu, z, true
    }
    
    return iteration, z, true
  }

  base := fractal.Base{rMin, rMax, iMin, iMax}

  
  if (mode == "image") {
    var points_map = fractal.Escape_Time_Calculator(base, midX, midY, zoom, width, height, calculator);
    if (filename == "") {
      filename = "mb_" + strconv.FormatFloat(midX, 'E', -1, 64) + "_" + strconv.FormatFloat(midY, 'E', -1, 64) + "_" +  strconv.FormatFloat(zoom, 'E', -1, 64) + ".jpg"
    }

    filename = output + "/" + filename

    draw_image(filename, points_map, width, height, gradient)
    fmt.Printf("%s\n", filename)
  } else if (mode == "edge") {
    var points_map = fractal.Escape_Time_Calculator(base, midX, midY, zoom, width, height, calculator);
    var edgePoints = make(chan fractal.Point)

    var found_edges []fractal.Point = make([]fractal.Point, 0)

    go func(edge<-chan fractal.Point) {
      for p := range edge {
        found_edges = append(found_edges, p)
      }
    }(edgePoints)

    find_edges(edgePoints, points_map)

    if (len(found_edges) == 0) {
      return
    }

    var index = int(rand.Float64() * float64(len(found_edges)))

    var p = found_edges[index].C
    fmt.Printf("%18.17e, %18.17e\n", real(p), imag(p))
  } else if (mode == "raw") {
    var points_map = fractal.Escape_Time_Calculator(base, midX, midY, zoom, width, height, calculator);
    
    if (filename == "") {
      filename = "/mb_" + strconv.FormatFloat(midX, 'E', -1, 64) + "_" + strconv.FormatFloat(midY, 'E', -1, 64) + "_" +  strconv.FormatFloat(zoom, 'E', -1, 64) + ".json"
    }

    filename = output + "/" + filename

    write_raw(points_map, filename)
  } else if (mode == "coordsAt") {
    var p = get_cordinates(midX, midY, zoom, width, height, x, y)
    fmt.Printf("%18.17e, %18.17e\n", real(p), imag(p))
  }
}
