package physics

import (
	"fmt"
	"math"
)

func predictWKBHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
	if z <= 0 || a <= 0 || a < z {
		return 0, fmt.Errorf("invalid nuclide: Z=%d A=%d", z, a)
	}
	if qAlphaMeV <= 0 {
		return 0, fmt.Errorf("Q_alpha must be positive, got %f", qAlphaMeV)
	}

	zf := float64(z)
	af := float64(a)

	par := classifyParity(z, a)
	var aC, bC, cC float64
	switch par {
	case ParityEvenEven:
		aC, bC, cC = -25.31, -1.1629, 1.5864
	case ParityOddZ:
		aC, bC, cC = -26.65, -1.0859, 1.5848
	case ParityOddN:
		aC, bC, cC = -25.68, -1.1423, 1.5920
	default:
		aC, bC, cC = -29.48, -1.1130, 1.6971
	}

	a6th := math.Pow(af, 1.0/6.0)
	zSqrt := math.Sqrt(zf)
	logT := aC + bC*a6th*zSqrt + cC*zf/math.Sqrt(qAlphaMeV)

	return logT, nil
}
