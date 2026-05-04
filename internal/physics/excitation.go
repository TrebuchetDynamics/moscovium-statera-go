package physics

import (
	"fmt"
	"math"
)

func ExcitationFunctionPeakXS(targetA, targetZ, projectileA, projectileZ, neutronsEvaporated int) (float64, error) {
	if targetA <= 0 || projectileA <= 0 || neutronsEvaporated < 0 {
		return 0, fmt.Errorf("invalid input: targetA=%d projectileA=%d neutronsEvaporated=%d", targetA, projectileA, neutronsEvaporated)
	}

	compoundA := targetA + projectileA
	compoundZ := targetZ + projectileZ

	r0 := 1.2
	rc := r0 * (math.Pow(float64(targetA), 1.0/3.0) + math.Pow(float64(projectileA), 1.0/3.0))
	coulombBarrier := 1.44 * float64(targetZ*projectileZ) / rc

	zProd := float64(targetZ * projectileZ)
	fusionProb := math.Exp(-0.4 * zProd / math.Sqrt(coulombBarrier))

	survivalProb := 1.0
	for range neutronsEvaporated {
		survivalProb *= 0.1
	}

	geometricXS := math.Pi * rc * rc * 1e-26

	// Ignore unused variables for fusion probability scaling
	_ = compoundA
	_ = compoundZ

	xs := geometricXS * fusionProb * survivalProb
	return xs, nil
}
