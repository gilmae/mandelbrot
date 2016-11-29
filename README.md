# Mandelbrot

A program that generates images of the Mandelbrot set.

## Usage
```
mandelbrot OPTIONS

-b float
    Bailout value. (default 4)
-c string
    Colour mode: true, smooth, banded, none. (default "none")
-f string
    Output file name.
-g string
    Gradient to use. (default "[[\"0.0\", \"000764\"],[\"0.16\", \"026bcb\"],[\"0.42\", \"edffff\"],[\"0.6425\", \"ffaa00\"],[\"0.8675\", \"000200\"],[\"1.0\",\"000764\"]]")
-h int
    Height of render. (default 1600)
-i float
    Imaginary component of the midpoint.
-m float
    Maximum Iterations before giving up on finding an escape. (default 2000)
-o string
    Output path. (default ".")
-r float
    Real component of the midpoint. (default -0.75)
-w int
    Width of render. (default 1600)
-z float
    Zoom level. (default 1)
```
