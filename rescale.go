package main

import (
	"math/big"
)

func RescaleArrayToBounds(bound_start big.Float, bound_end big.Float, points []big.Float) []big.Float {
	old_min := new(big.Float)
	old_max := new(big.Float)
	old_min.Copy(&points[0])
	old_max.Copy(&points[len(points)-1])

	original_length := new(big.Float)
	original_length.Sub(old_max, old_min)

	new_length := new(big.Float)
	new_length.Sub(&bound_end, &bound_start)

	ratio := new(big.Float)
	ratio.Quo(original_length, new_length)

	mapped := make([]big.Float, len(points))

	for i, v := range points {
		mapped[i].Set(&bound_start)
		tmp := new(big.Float).Set(&v)
		tmp.Sub(tmp, old_min)
		tmp.Quo(tmp, ratio)
		mapped[i].Add(&mapped[i], tmp)
	}
	return mapped
}

func Rescale(line_start float64, line_end float64, length int, point int) float64 {
	step := (line_end - line_start) / float64(length)
	line_focal_point := (line_end + line_start) / 2.0

	return line_focal_point + (float64(point)-(float64(length)/2.0))*step
}

func GetZoomedBounds(original_start float64, original_end float64, focal_point float64, zoom float64) (float64, float64) {
	original_width := original_end - original_start
	original_width_to_midpoint := original_width / 2.0

	new_width_to_mid_point := original_width_to_midpoint / zoom

	return focal_point - new_width_to_mid_point, focal_point + new_width_to_mid_point
}
