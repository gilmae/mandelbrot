# Mandelbrot

A program that generates images of the Mandelbrot set.

Given a co-ordinate on the complex plane and a zoom level, generates a 1600 x 1037px PNG depicting a snapshot of the mandelbrot set.

If no particular value is given for the real or imaginary components of the coplex number, or for the zoom level, the program will randomise that argument, within the following ranges:

* Real: -2.6 - 1.1
* Imaginary: -1.1 - 1.1
* Zoom: 1 - 2^10 + 1

Because generally a randomly generated snapshot of the mandelbrot set is dreadfully boring - often deep past the event horizon or every point on the image escapes at the same time. 
So this program has a 'boringness' filter, which really just looks at a histogram and dimisses an image as boring if there is too little variation within the histogram results.

## Usage

```
$ cd $GOPATH
$ git clone git@github.com:gilmae/mandelbrot.git ./src/github.com/gilmae/mandelbrot
$ git install github.com/gilmae/mandelbrot
$ bin/mandelbrot <output_path> <real> <imaginary> <zoom>
```

To randomise any of real, imaginary, or zoom, enter a value of .

