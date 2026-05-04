package physics

import (
	"fmt"
	"math"
)

func predictVSSHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
	if err := validateAlphaPredictInput(z, a, qAlphaMeV); err != nil {
		return 0, err
	}
	zf := float64(z)
	const (
		aVSS = 1.66175
		bVSS = -8.5166
		cVSS = -0.20228
		dVSS = -33.9069
	)
	sqrtQ := math.Sqrt(qAlphaMeV)
	logT := (aVSS*zf+bVSS)/sqrtQ + cVSS*zf + dVSS
	par := classifyParity(z, a)
	switch par {
	case ParityOddZ, ParityOddN:
		logT += 0.6
	case ParityOddOdd:
		logT += 1.2
	}
	return logT, nil
}

func predictUNIVHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
	if err := validateAlphaPredictInput(z, a, qAlphaMeV); err != nil {
		return 0, err
	}
	zf := float64(z)
	af := float64(a)
	const (
		aUNIV = 0.3940
		bUNIV = 1.5617
		cUNIV = -48.7633
	)
	chiPrime := zf / math.Sqrt(qAlphaMeV)
	rhoPrime := math.Sqrt(af * zf)
	return aUNIV*rhoPrime + bUNIV*chiPrime + cUNIV, nil
}

func predictDenisovHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
	if err := validateAlphaPredictInput(z, a, qAlphaMeV); err != nil {
		return 0, err
	}
	zf := float64(z)
	af := float64(a)
	const (
		aD = -28.942
		bD = -0.049
		cD = 0.589
		dD = 0.968
	)
	return aD + bD*af + cD*math.Pow(zf, 2)/math.Sqrt(qAlphaMeV) + dD*zf/math.Sqrt(qAlphaMeV), nil
}

func validateAlphaPredictInput(z, a int, qAlphaMeV float64) error {
	if z <= 0 || a <= 0 || a < z {
		return errInvalidNuclide(z, a)
	}
	if qAlphaMeV <= 0 || math.IsNaN(qAlphaMeV) || math.IsInf(qAlphaMeV, 0) {
		return errInvalidQAlpha(qAlphaMeV)
	}
	return nil
}

func errInvalidNuclide(z, a int) error { return &alphaError{msg: fmt.Sprintf("invalid nuclide: Z=%d A=%d", z, a)} }
func errInvalidQAlpha(q float64) error { return &alphaError{msg: fmt.Sprintf("Q_alpha must be positive finite, got %v", q)} }

type alphaError struct{ msg string }

func (e *alphaError) Error() string { return e.msg }

// Public wrappers for consuming packages.
func PredictVSSHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
	return predictVSSHalfLife(z, a, qAlphaMeV)
}
func PredictUNIVHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
	return predictUNIVHalfLife(z, a, qAlphaMeV)
}
func PredictDenisovHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
	return predictDenisovHalfLife(z, a, qAlphaMeV)
}
func PredictWKBHalfLife(z, a int, qAlphaMeV float64) (float64, error) {
	return predictWKBHalfLife(z, a, qAlphaMeV)
}
