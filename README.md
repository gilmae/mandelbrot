# Mandelbrot

A program that generates images of the Mandelbrot set.

Given a co-ordinate on the complex plane and a zoom level, generates a 1600 x 1600px PNG depicting a snapshot of the mandelbrot set.

## Usage

```
$ cd $GOPATH
$ git clone git@github.com:gilmae/mandelbrot.git ./src/github.com/gilmae/mandelbrot
$ git install github.com/gilmae/mandelbrot
$ bin/mandelbrot -r=<real> -i=<imaginary> -z=<zoom>
```
