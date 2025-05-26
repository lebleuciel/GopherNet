package utils

import (
	"math"

	"gophernet/pkg/db/ent"
)

func CalculateVolume(burrow *ent.Burrow) float64 {
	// Volume of a cylinder: V = π * r^2 * h (where r is radius, h is height)
	radius := burrow.Width / 2
	volume := math.Pi * math.Pow(radius, 2) * burrow.Depth
	return volume
}
