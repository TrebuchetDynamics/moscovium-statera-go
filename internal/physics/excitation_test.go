package physics

import "testing"

func TestExcitationFunctionPeakXS_Am243_Ca48_3n(t *testing.T) {
	xs, err := ExcitationFunctionPeakXS(243, 95, 48, 20, 3)
	if err != nil {
		t.Fatalf("excitation function: %v", err)
	}
	if xs <= 0 || xs > 1e-34 {
		t.Errorf("unexpected cross-section: %.1e cm^2", xs)
	}
}

func TestExcitationFunctionPeakXS_InvalidInputs(t *testing.T) {
	_, err := ExcitationFunctionPeakXS(0, 95, 48, 20, 3)
	if err == nil {
		t.Error("expected error for zero targetA")
	}
	_, err = ExcitationFunctionPeakXS(243, 95, 48, 20, -1)
	if err == nil {
		t.Error("expected error for negative neutrons")
	}
}

func TestExcitationFunctionPeakXS_DifferentChannels(t *testing.T) {
	for _, xn := range []int{2, 3, 4} {
		xs, err := ExcitationFunctionPeakXS(243, 95, 48, 20, xn)
		if err != nil {
			t.Errorf("%d neutrons: %v", xn, err)
			continue
		}
		if xs <= 0 {
			t.Errorf("%d neutrons: non-positive cross-section", xn)
		}
	}
}
