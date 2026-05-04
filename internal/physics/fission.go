package physics

import (
	"fmt"
	"math"
)

func SwiateckiSFHalfLife(z, a int) (float64, error) {
	if z <= 0 || a <= 0 || a < z {
		return 0, fmt.Errorf("invalid nuclide: Z=%d A=%d", z, a)
	}

	zf := float64(z)
	af := float64(a)

	const criticalZ2A = 50.883
	const (
		k = 120.0
		n = 3.0
		c = 10.0
	)

	z2a := zf * zf / af
	if z2a >= criticalZ2A {
		return -20.0, nil
	}

	x := z2a / criticalZ2A
	logT := c + k*math.Pow(1.0-x, n)
	return logT, nil
}

func CompetitiveDecayMode(z, a int, qAlphaMeV, alphaLogT, sfLogT float64) string {
	if sfLogT < alphaLogT {
		return "SF"
	}
	if qAlphaMeV <= 0 {
		return "SF"
	}
	return "alpha"
}
