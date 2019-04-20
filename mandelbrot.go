package main

import (
	"fmt"
	"math"
	mbig "math/big"

	//"math/cmplx"
	"flag"
	"math/rand"
	"strconv"
	"time"

	"github.com/gilmae/fractal"
	"github.com/gilmae/rescale/big"
)

var maxIterations float64 = 1000.0
var bailout float64 = 4.0
var width int = 1600
var height int = 1600
var x int = 0
var y int = 0
var prec uint = 128

var default_gradient string = `[["0.0", "000764"],["0.16", "026bcb"],["0.42", "edffff"],["0.6425", "ffaa00"],["0.8675", "000200"],["1.0","000764"]]`

const (
	rMin  = -2.25
	rMax  = 0.75
	iMin  = -1.5
	iMax  = 1.5
	usage = "mandelbot OPTIONS\n\nPlots the mandelbrot set, centered at a point indicated by the provided real and imaginary, and at the given zoom level.\n\nSaves the output into the given path.\n\n"
)

func get_big_float(val float64) mbig.Float {
	bigVal := mbig.NewFloat(val).SetPrec(prec)
	return *bigVal
}

func get_cordinates(midX mbig.Float, midY mbig.Float, zoom float64, width int, height int, x int, y int) (mbig.Float, mbig.Float) {
	bigRMin := get_big_float(rMin)
	bigRMax := get_big_float(rMax)
	bigIMin := get_big_float(iMin)
	bigIMax := get_big_float(iMax)

	new_r_start, new_r_end := big.GetZoomedBounds(&bigRMin, &bigRMax, &midX, zoom)
	scaled_r := big.Rescale(new_r_start, new_r_end, width, x)

	new_i_start, new_i_end := big.GetZoomedBounds(&bigIMin, &bigIMax, &midY, zoom)
	scaled_i := big.Rescale(new_i_start, new_i_end, height, height-y)

	return scaled_r, scaled_i
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
	var colour_mode string = ""

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
	flag.UintVar(&prec, "p", 128, "Precision of floats")
	flag.Parse()

	bigX := get_big_float(midX)
	bigY := get_big_float(midY)

	var calculator fractal.EscapeCalculator = func(seed fractal.BigComplex) (float64, fractal.BigComplex, bool) {
		iteration := 0.0
		bigBailout := mbig.NewFloat(bailout)

		c := new(fractal.BigComplex)
		c.Set(&seed)
		z := new(fractal.BigComplex)
		z.Set(&seed)

		for ; fractal.Abs(z).Cmp(bigBailout) < 0 && iteration < maxIterations; iteration += 1 {
			z.Mul(z, z)
			z.Add(z, c)
		}

		if iteration >= maxIterations {
			return maxIterations, *z, false
		}

		if mode == "image" && colour_mode == "smooth" {
			z.Mul(z, z)
			z.Add(z, c)
			z.Mul(z, z)
			z.Add(z, c)
			iteration += 2
			reZ := fractal.Real(z)
			imZ := fractal.Imag(z)
			bigMagnitude := new(mbig.Float).Set(reZ)
			bigMagnitude.Mul(bigMagnitude, bigMagnitude)

			imagsquared := new(mbig.Float).Set(imZ)
			imagsquared.Mul(imagsquared, imagsquared)

			bigMagnitude.Add(bigMagnitude, imagsquared)
			bigMagnitude.Sqrt(bigMagnitude)

			magnitude, _ := bigMagnitude.Float64()

			mu := iteration + 1 - (math.Log(math.Log(magnitude)))/math.Log(2.0)

			return mu, *z, true
		}

		return iteration, *z, true
	}

	base := fractal.Base{get_big_float(rMin), get_big_float(rMax), get_big_float(iMin), get_big_float(iMax)}

	if mode == "image" {
		var points_map = fractal.Escape_Time_Calculator(base, midX, midY, zoom, width, height, calculator)
		if filename == "" {
			filename = "mb_" + strconv.FormatFloat(midX, 'E', -1, 64) + "_" + strconv.FormatFloat(midY, 'E', -1, 64) + "_" + strconv.FormatFloat(zoom, 'E', -1, 64) + ".jpg"
		}

		filename = output + "/" + filename

		fractal.Draw_Image(filename, points_map, width, height, gradient, maxIterations, colour_mode)
		fmt.Printf("%s\n", filename)
	} else if mode == "raw" {
		var points_map = fractal.Escape_Time_Calculator(base, midX, midY, zoom, width, height, calculator)

		if filename == "" {
			filename = "/mb_" + strconv.FormatFloat(midX, 'E', -1, 64) + "_" + strconv.FormatFloat(midY, 'E', -1, 64) + "_" + strconv.FormatFloat(zoom, 'E', -1, 64) + ".json"
		}

		filename = output + "/" + filename

		fractal.Write_Raw(points_map, filename)
	} else if mode == "coordsAt" {
		var x, y = get_cordinates(bigX, bigY, zoom, width, height, x, y)
		fmt.Printf("%s, %s\n", x.Text('e', 100), y.Text('e', 100))
	}
}
