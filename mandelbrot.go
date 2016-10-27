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
    "github.com/gilmae/interpolation"
)
var  maxIterations float64 = 1000.0
var  bailout float64 = 4.0
var  width int = 1600
var height int = 1600
var paletteLength int = 16;
var colour_mode string = "";

/*var xSequence = []float64{0.0, .16, .42, .6425, .8675, 1}
var redpoints =  []float64{0.0,32.0,237.0,255.0,0.0,0.0}
var greenpoints =  []float64{7.0,107.0,255.0, 170.0, 2.0, 7.0}
var bluepoints =  []float64{100.0, 203.0,255.0, 0.0, 0.0, 100.0}*/

var xSequence = []float64{0.0, .06666, .13333, .2, .26666, .33333, .4, .466666, .53333, .6, .66666, .73333, .8, .86666, .93333, 1.0}
var redpoints =  []float64{66.0, 25.0, 9.0, 4.0, 0.0, 12.0, 24.0, 57.0, 134.0, 211.0, 241.0, 248.0, 255.0, 204.0, 153.0, 106.0, }
var greenpoints =  []float64{30.0, 7.0, 1.0, 4.0, 7.0, 44.0, 82.0, 125.0, 181.0, 236.0, 233.0, 201.0, 170.0, 128.0, 87.0, 52.0}
 var bluepoints =  []float64{15.0, 26.0, 47.0, 73.0, 100.0, 138.0, 177.0, 209.0, 229.0, 248.0, 191.0, 95.0, 0.0, 0.0, 0.0, 3.0}

var palette = make([]color.NRGBA, paletteLength)

var redInterpolant = interpolation.CreateMonotonicCubic(xSequence, redpoints)
var greenInterpolant = interpolation.CreateMonotonicCubic(xSequence, greenpoints)
var blueInterpolant = interpolation.CreateMonotonicCubic(xSequence, bluepoints)

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
    usage  = "mandelbot output_path real imaginary zoom\n\n Plots the mandelbrot set, centered at point indicated by real,imaginary and at the given zoom level.\n\nSaves the output into the given path.\n\n"
)

func calculate_escape(c complex128) float64 {
  iteration := 0.0

  var z complex128
  for z= c;cmplx.Abs(z) < bailout && iteration < maxIterations; iteration+=1 {
    z = z*z+c;
  }

  if (iteration >= maxIterations) {
    return maxIterations
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

func fill_palette() {
  for i:= 0; i < paletteLength; i++ {
    var point = 1.0 * float64(i) / float64(paletteLength)
    var redpoint = redInterpolant(point)
    var greenpoint = greenInterpolant(point)
    var bluepoint = blueInterpolant(point)

    palette[i] = color.NRGBA{uint8(redpoint), uint8(greenpoint), uint8(bluepoint), 255}

  }
}

func get_colour(esc float64) color.NRGBA {
  if esc >= maxIterations{
    return color.NRGBA{0, 0, 0, 255}
  }

  if (colour_mode == "true") {
    var point = esc/float64(maxIterations)
    var redpoint = redInterpolant(point)
    var greenpoint = greenInterpolant(point)
    var bluepoint = blueInterpolant(point)

    return color.NRGBA{uint8(redpoint), uint8(greenpoint), uint8(bluepoint), 255}
  } else if (colour_mode == "smooth") {
    clr1 := int(esc)
    t2 :=  esc - float64(clr1);
    t1 := 1 - t2;

    clr1 = clr1 % len(palette)
    clr2 := (clr1 + 1) % len(palette)

    r := float64(palette[clr1].R) * t1 + float64(palette[clr2].R) * t2
    g := float64(palette[clr1].G) * t1 + float64(palette[clr2].G) * t2
    b := float64(palette[clr1].B) * t1 + float64(palette[clr2].B) * t2

    return color.NRGBA{uint8(r),uint8(g),uint8(b),255};
  } else if (colour_mode == "banded") {
    return palette[int(esc) % len(palette)]
  } else {
    return color.NRGBA{255, 255, 255, 255};
  }

}

func main() {
  start := time.Now()

  var midX float64
  var midY float64
  var zoom float64
  var output string
  var filename string

  rand.Seed(time.Now().UnixNano())
  flag.Float64Var(&midX, "r", -0.75, "Real component of the midpoint. Defaults tp -0.75.")
  flag.Float64Var(&midY, "i", 0.0, "Imaginary component of the midpoint. Defaults to 0.0.")
  flag.Float64Var(&zoom, "z", 1, "Zoom level. Defaults to 1.0.")
  flag.StringVar(&output, "o", ".", "Output path. Defaults to current path.")
  flag.StringVar(&filename, "f", "", "Output file name.")
  flag.StringVar(&colour_mode, "c", "", "Colour mode: true, smooth, banded, none. Defaults to none.")
  flag.Float64Var(&bailout, "b", 4.0, "Bailout value. Defaults to 4.0")
  flag.IntVar(&width, "w", 1600, "Width. Defaults to 1600.")
  flag.IntVar(&height, "h", 1600, "Height. Defaults to 1600.")
  flag.Float64Var(&maxIterations, "m", 2000.0, "Maximum Iterations. Defaults to 2000")
  flag.Parse()

  fill_palette()

  if (filename == "") {
    filename = "/mb_" + strconv.FormatFloat(midX, 'E', -1, 64) + "_" + strconv.FormatFloat(midY, 'E', -1, 64) + "_" +  strconv.FormatFloat(zoom, 'E', -1, 64) + ".jpg"
  }

  filename = output + "/" + filename

  scale := (float64(width) / (rMax - rMin))
  scale = scale * zoom

  bounds := image.Rect(0,0,width,height)
  b := image.NewNRGBA(bounds)
  draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

  calculatedChan := make(chan Point)

  go func(points<-chan Point, targetImage *image.NRGBA) {
    for p := range points {
      targetImage.Set(p.x,p.y, get_colour(p.escape))
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
