package physics

import (
	"math"
	"testing"
)

func TestBetheWeizsackerBindingEnergy_288Mc(t *testing.T) {
	be, err := BetheWeizsackerBindingEnergy(115, 288)
	if err != nil {
		t.Fatalf("binding energy: %v", err)
	}
	perNucleon := be / 288.0
	if perNucleon < 6.0 || perNucleon > 8.5 {
		t.Errorf("unexpected B/A = %.2f MeV for 288Mc", perNucleon)
	}
}

func TestBindingEnergyPerNucleon_SuperheavyTrend(t *testing.T) {
	for _, id := range []string{"288Mc", "290Mc", "284Nh", "280Rg"} {
		var z, a int
		switch id {
		case "288Mc":
			z, a = 115, 288
		case "290Mc":
			z, a = 115, 290
		case "284Nh":
			z, a = 113, 284
		case "280Rg":
			z, a = 111, 280
		}
		ba, err := BindingEnergyPerNucleon(z, a)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if ba < 6.5 || ba > 8.0 {
			t.Errorf("%s B/A = %.2f MeV, expected 7.0±0.5", id, ba)
		}
	}
}

func TestShellCorrectionEstimate_MagicNumbers(t *testing.T) {
	sc, err := ShellCorrectionEstimate(114, 298)
	if err != nil {
		t.Fatal(err)
	}
	if sc >= 0 {
		t.Errorf("expected negative shell correction near Z=114, N=184, got %.2f", sc)
	}

	scFar, _ := ShellCorrectionEstimate(100, 250)
	if math.Abs(scFar) > math.Abs(sc) {
		t.Errorf("expected smaller correction far from magic numbers, got %.2f vs %.2f", scFar, sc)
	}
}

func TestBetheWeizsacker_InvalidInputs(t *testing.T) {
	_, err := BetheWeizsackerBindingEnergy(0, 288)
	if err == nil {
		t.Error("expected error for Z=0")
	}
	_, err = BetheWeizsackerBindingEnergy(115, 100)
	if err == nil {
		t.Error("expected error for A < Z")
	}
}
