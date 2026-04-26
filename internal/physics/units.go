package physics

import (
	"fmt"
	"math/big"
)

const (
	precisionBits = 256

	// Exact SI value from the 2019 SI redefinition.
	electronVoltJoule = "1.602176634e-19"
)

func MeVToJoules(mev string) (*big.Float, error) {
	value, ok := new(big.Float).SetPrec(precisionBits).SetMode(big.ToNearestEven).SetString(mev)
	if !ok {
		return nil, fmt.Errorf("invalid MeV value %q", mev)
	}

	ev, ok := new(big.Float).SetPrec(precisionBits).SetMode(big.ToNearestEven).SetString(electronVoltJoule)
	if !ok {
		return nil, fmt.Errorf("invalid electron volt constant %q", electronVoltJoule)
	}

	million := new(big.Float).SetPrec(precisionBits).SetFloat64(1_000_000)
	scale := new(big.Float).SetPrec(precisionBits).Mul(ev, million)
	return new(big.Float).SetPrec(precisionBits).Mul(value, scale), nil
}
