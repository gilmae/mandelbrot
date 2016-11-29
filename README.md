# Mandelbrot

A program that generates images of the Mandelbrot set.

## Usage

mandelbrot OPTIONS

```
-b float
    Bailout value. Defaults to 4.0 (default 4)
-c string
    Colour mode: true, smooth, banded, none. Defaults to none.
-f string
    Output file name.
-g string
    Gradient to use. (default "[[\"0.0\", \"000764\"],[\"0.16\", \"026bcb\"],[\"0.42\", \"edffff\"],[\"0.6425\", \"ffaa00\"],[\"0.8675\", \"000200\"],[\"1.0\",\"000764\"]]")
-h int
    Height. Defaults to 1600. (default 1600)
-i float
    Imaginary component of the midpoint. Defaults to 0.0.
-m float
    Maximum Iterations. Defaults to 2000 (default 2000)
-o string
    Output path. Defaults to current path. (default ".")
-r float
    Real component of the midpoint. Defaults to -0.75. (default -0.75)
-w int
    Width. Defaults to 1600. (default 1600)
-z float
    Zoom level. Defaults to 1.0. (default 1)
```
