package main

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "math/cmplx"
    "os"
    "strconv"
)

const (
    maxEsc = 255.0
    rMin   = -2.6
    rMax   = 1.1
    iMin   = -1.2
    iMax   = 1.2
    width  = 1600
    red    = 230
    green  = 235
    blue   = 255
    usage  = "mandelbot real imaginary zoom output_path\n\nPlots the mandelbrot set, centered at point indicated by real,imaginary and at the given zoom level.\n\nSaves the output into the given path.\n"
)

func plot(c complex128) int {
  i := 0

  for z:= c;cmplx.Abs(z) < 2.0 && i < maxEsc; i+=1 {
    z = z*z+c;
  }

  return i;
}

func get_colour(esc int) color.NRGBA {
  if esc == maxEsc{
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
  return palette[esc % 16]
}

func main() {
  if (len(os.Args) < 5) {
    fmt.Println(usage)
    return
  }
  
  midX,err := strconv.ParseFloat(os.Args[1], 64)
  if (err != nil) {
    fmt.Println(err)
    return
  }
  midY, err := strconv.ParseFloat(os.Args[2], 64)
  if (err != nil) {
    fmt.Println(err)
    return
  }


  zoom, err := strconv.ParseFloat(os.Args[3], 64)
  if (err != nil) {
    fmt.Println(err)
    return
  }

  out_path := os.Args[4]

  scale := (width / (rMax - rMin))
  height := int(scale * (iMax-iMin))
  scale = scale * zoom
  bounds := image.Rect(0,0,width,height)
  b := image.NewNRGBA(bounds)
  draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

  for x:=0; x < width; x += 1 {
    for y:=0; y < height; y += 1 {
      esc := plot(complex(float64(x - width/2)/scale + midX, float64(y - height/2)/scale+midY))
      b.Set(x,y, get_colour(esc))
    }
  }

  filename := out_path + "/mb_" + os.Args[1] + "_" + os.Args[2] + "_" + os.Args[3] + ".png"
  f, err := os.Create(filename)
  if err != nil {
    fmt.Println(err)
  }

  if err = png.Encode(f,b); err != nil {
    fmt.Println(err)
  }

  if err = f.Close();err != nil {
    fmt.Println(err)
  }
}
