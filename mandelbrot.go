package main

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "math/cmplx"
    "os"
    "strconv"
    "math"
    "math/rand"
    "time"
    "flag"
    "sync"
    "log"
)
var  maxIterations float64 = 2000.0
var  bailout float64 = 4.0

type Point struct {
   c complex128
   x int
   y int
   escape float64
}

const (
    rMin   = -2.25
    rMax   = 0.75
    iMin   = -1.5
    iMax   = 1.5
    width  = 1600
    usage  = "mandelbot output_path real imaginary zoom\n\n Plots the mandelbrot set, centered at point indicated by real,imaginary and at the given zoom level.\n\nSaves the output into the given path.\n\n"
)

func calculate_escape(c complex128) float64 {
  iteration := 0.0

  var z complex128
  for z= c;cmplx.Abs(z) < bailout && iteration < maxIterations; iteration+=1 {
    z = z*z+c;
  }

  if (iteration >= maxIterations) {
    return maxIterations;
  }

  z = z*z+c
  z = z*z+c
  iteration += 2
  reZ := real(z)
  imZ := imag(z)
  magnitude := math.Sqrt(reZ * reZ + imZ * imZ)
  mu := iteration + 1 - (math.Log(math.Log(magnitude)))/math.Log(2.0)
  return mu
}

func plot(midX float64, midY float64, scale float64, width int, height int, calculated chan Point) {
  points := make(chan Point, 64)

  // spawn four worker goroutines
  var wg sync.WaitGroup
  for i := 0; i < 4; i++ {
    wg.Add(1)
    go func() {
      for p := range points {
        p.escape = calculate_escape(p.c)
        calculated <- p
      }
      wg.Done()
    }()
  }

  for x:=0; x < width; x += 1 {
    for y:=0; y < height; y += 1 {
      points <- Point{complex(float64(x - width/2)/scale + midX, float64((height-y) - height/2)/scale+midY),x,y, 0}
    }
  }

  close(points)

  wg.Wait()
}

func get_colour(esc float64, smooth bool) color.NRGBA {
  if esc >= maxIterations{
    return color.NRGBA{0, 0, 0, 255}
  }

  palette := [16]color.NRGBA{
    color.NRGBA{66, 30, 15, 255},
    color.NRGBA{25, 7, 26, 255},
    color.NRGBA{9, 1, 47, 255},
    color.NRGBA{4, 4, 73, 255},
    color.NRGBA{0, 7, 100, 255},
    color.NRGBA{12, 44, 138, 255},
    color.NRGBA{24, 82, 177, 255},
    color.NRGBA{57, 125, 209, 255},
    color.NRGBA{134, 181, 229, 255},
    color.NRGBA{211, 236, 248, 255},
    color.NRGBA{241, 233, 191, 255},
    color.NRGBA{248, 201, 95, 255},
    color.NRGBA{255, 170, 0, 255},
    color.NRGBA{204, 128, 0, 255},
    color.NRGBA{153, 87, 0, 255},
    color.NRGBA{106, 52, 3, 255}}

  if (smooth) {
    clr1 := int(esc)
    t2 :=  esc - float64(clr1);
    t1 := 1 - t2;

    clr1 = clr1 % len(palette)
    clr2 := (clr1 + 1) % len(palette)

    r := float64(palette[clr1].R) * t1 + float64(palette[clr2].R) * t2
    g := float64(palette[clr1].G) * t1 + float64(palette[clr2].G) * t2
    b := float64(palette[clr1].B) * t1 + float64(palette[clr2].B) * t2

    return color.NRGBA{uint8(r),uint8(g),uint8(b),255};
  } else {
    return palette[int(esc) % len(palette)]
  }
}

func main() {
  start := time.Now()

  var midX float64
  var midY float64
  var zoom float64
  var output string
  var filename string
  var smooth bool

  rand.Seed(time.Now().UnixNano())
  flag.Float64Var(&midX, "r", -0.75, "Real component of the midpoint. Defaults tp -0.75.")
  flag.Float64Var(&midY, "i", 0.0, "Imaginary component of the midpoint. Defautls to 0.0.")
  flag.Float64Var(&zoom, "z", 1, "Zoom level. Defaults to 1.0.")
  flag.StringVar(&output, "o", ".", "Output path. Defaults to current path.")
  flag.StringVar(&filename, "f", "", "Output file name.")
  flag.BoolVar(&smooth, "s", true, "Smooth colours.")
  flag.Float64Var(&bailout, "b", 4.0, "Bailout value.")
  flag.Float64Var(&maxIterations, "m", 2000.0, "Maximum Iterations.")
  flag.Parse()

  if (filename == "") {
    filename = "/mb_" + strconv.FormatFloat(midX, 'E', -1, 64) + "_" + strconv.FormatFloat(midY, 'E', -1, 64) + "_" +  strconv.FormatFloat(zoom, 'E', -1, 64) + ".jpg"
  }

  filename = output + "/" + filename

  scale := (width / (rMax - rMin))
  scale = scale * zoom

  height := int(scale * (iMax-iMin))
  bounds := image.Rect(0,0,width,height)
  b := image.NewNRGBA(bounds)
  draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

  calculatedChan := make(chan Point)

  go func(points<-chan Point, targetImage *image.NRGBA) {
    for p := range points {
      targetImage.Set(p.x,p.y, get_colour(p.escape, smooth))
    }
  }(calculatedChan, b)

  plot(midX, midY, scale, width, height, calculatedChan)

  file, err := os.Create(filename)
  if err != nil {
    fmt.Println(err)
  }

  if err = jpeg.Encode(file,b, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
    fmt.Println(err)
  }

  if err = file.Close();err != nil {
    fmt.Println(err)
  }

  elapsed := time.Since(start)
  log.Printf("Plot took %s", elapsed)

}
