package physics

import (
	"fmt"
	"math"
)

func BetheWeizsackerBindingEnergy(z, a int) (float64, error) {
	if z <= 0 || a <= 0 || a < z {
		return 0, fmt.Errorf("invalid nuclide: Z=%d A=%d", z, a)
	}
	n := a - z

	const (
		aV = 15.75
		aS = 17.8
		aC = 0.711
		aA = 23.7
		aP = 11.18
	)

	af := float64(a)
	zf := float64(z)
	nf := float64(n)

	volume := aV * af
	surface := -aS * math.Pow(af, 2.0/3.0)
	coulomb := -aC * zf * (zf - 1) / math.Pow(af, 1.0/3.0)
	asymmetry := -aA * (nf - zf) * (nf - zf) / af

	var pairing float64
	zEven := z%2 == 0
	nEven := n%2 == 0
	switch {
	case zEven && nEven:
		pairing = aP / math.Sqrt(af)
	case !zEven && !nEven:
		pairing = -aP / math.Sqrt(af)
	default:
		pairing = 0
	}

	return volume + surface + coulomb + asymmetry + pairing, nil
}

func BindingEnergyPerNucleon(z, a int) (float64, error) {
	be, err := BetheWeizsackerBindingEnergy(z, a)
	if err != nil {
		return 0, err
	}
	return be / float64(a), nil
}

func ShellCorrectionEstimate(z, a int) (float64, error) {
	if z <= 0 || a <= 0 {
		return 0, fmt.Errorf("invalid nuclide")
	}
	n := a - z
	magicZ := []int{114, 126}
	magicN := []int{162, 172, 178, 184}

	shellEnergy := 0.0
	for _, mz := range magicZ {
		dz := math.Abs(float64(z - mz))
		if dz < 10 {
			shellEnergy -= 3.0 * math.Exp(-0.5 * dz * dz / 25.0)
		}
	}
	for _, mn := range magicN {
		dn := math.Abs(float64(n - mn))
		if dn < 10 {
			shellEnergy -= 3.0 * math.Exp(-0.5 * dn * dn / 25.0)
		}
	}
	return shellEnergy, nil
}
