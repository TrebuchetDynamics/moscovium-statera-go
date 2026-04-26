package physics

import "testing"

func TestMeVToJoulesUsesExactElectronVoltConstant(t *testing.T) {
	joules, err := MeVToJoules("1")
	if err != nil {
		t.Fatalf("MeVToJoules returned error: %v", err)
	}

	got := joules.Text('e', 9)
	want := "1.602176634e-13"
	if got != want {
		t.Fatalf("1 MeV = %s J, want %s J", got, want)
	}
}

func TestMeVToJoulesRejectsInvalidDecimal(t *testing.T) {
	if _, err := MeVToJoules("not-a-number"); err == nil {
		t.Fatal("MeVToJoules returned nil error for invalid decimal")
	}
}
