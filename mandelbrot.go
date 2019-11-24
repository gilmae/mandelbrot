package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"strconv"
	"time"

	"github.com/gilmae/fractal"
	"github.com/gilmae/rescale"
)

var maxIterations = 1000
var bailout = 4.0
var width = 1600
var height = 1600
var x = 0
var y = 0

var defaultGradient = `[["0.0", "000764"],["0.16", "026bcb"],["0.42", "edffff"],["0.6425", "ffaa00"],["0.8675", "000200"],["1.0","000764"]]`

const (
	rMin  = -2.25
	rMax  = 0.75
	iMin  = -1.5
	iMax  = 1.5
	usage = "mandelbot OPTIONS\n\nPlots the mandelbrot set, centered at a point indicated by the provided real and imaginary, and at the given zoom level.\n\nSaves the output into the given path.\n\n"
)

func getCordinates(midX float64, midY float64, zoom float64, width int, height int, x int, y int) complex128 {
	newRStart, newREnd := rescale.GetZoomedBounds(rMin, rMax, midX, zoom)
	scaledR := rescale.Rescale(newRStart, newREnd, width, x)

	newIStart, newIEnd := rescale.GetZoomedBounds(iMin, iMax, midY, zoom)
	scaledI := rescale.Rescale(newIStart, newIEnd, height, height-y)

	return complex(scaledR, scaledI)
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
	var colourMode string

	rand.Seed(time.Now().UnixNano())
	flag.Float64Var(&midX, "r", -0.75, "Real component of the midpoint.")
	flag.Float64Var(&midY, "i", 0.0, "Imaginary component of the midpoint.")
	flag.Float64Var(&zoom, "z", 1, "Zoom level.")
	flag.StringVar(&output, "o", ".", "Output path.")
	flag.StringVar(&filename, "f", "", "Output file name.")
	flag.StringVar(&colourMode, "c", "none", "Colour mode: true, smooth, banded, none.")
	flag.Float64Var(&bailout, "b", 4.0, "Bailout value.")
	flag.IntVar(&width, "w", 1600, "Width of render.")
	flag.IntVar(&height, "h", 1600, "Height of render.")
	flag.IntVar(&maxIterations, "m", 2000, "Maximum Iterations before giving up on finding an escape.")
	flag.StringVar(&gradient, "g", defaultGradient, "Gradient to use.")
	flag.StringVar(&mode, "mode", "image", "Mode: image, coordsAt")
	flag.IntVar(&x, "x", 0, "x cordinate of a pixel, used for translating to the real component. 0,0 is top left.")
	flag.IntVar(&y, "y", 0, "y cordinate of a pixel, used for translating to the real component. 0,0 is top left.")
	flag.Parse()

	var calculator fractal.EscapeCalculator = func(z complex128) (int, complex128, bool) {
		iteration := 0
		c := z

		for z = c; cmplx.Abs(z) < bailout && iteration < maxIterations; iteration++ {
			z = z*z + c
		}

		if iteration >= maxIterations {
			return maxIterations, z, false
		}

		if mode == "image" && colourMode == "smooth" {
			z = z*z + c
			z = z*z + c
			iteration += 2
			reZ := real(z)
			imZ := imag(z)
			magnitude := math.Sqrt(reZ*reZ + imZ*imZ)
			mu := float64(iteration+1) - (math.Log(math.Log(magnitude)))/math.Log(2.0)

			return int(mu), z, true
		}

		return iteration, z, true
	}

	base := fractal.Base{rMin, rMax, iMin, iMax}

	if mode == "image" {
		var points_map = fractal.EscapeTimeCalculator(base, midX, midY, zoom, width, height, calculator)
		if filename == "" {
			filename = "mb_" + strconv.FormatFloat(midX, 'E', -1, 64) + "_" + strconv.FormatFloat(midY, 'E', -1, 64) + "_" + strconv.FormatFloat(zoom, 'E', -1, 64) + ".jpg"
		}

		filename = output + "/" + filename

		fractal.Draw_Image(filename, points_map, width, height, gradient, maxIterations, colourMode)
		fmt.Printf("%s\n", filename)
	} else if mode == "coordsAt" {
		var p = getCordinates(midX, midY, zoom, width, height, x, y)
		fmt.Printf("%18.17e, %18.17e\n", real(p), imag(p))
	}
}
