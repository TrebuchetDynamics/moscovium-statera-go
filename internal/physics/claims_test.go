package physics

import (
	"strings"
	"testing"
	"time"
)

func TestValidateClaimRejectsLongLivedKnownShortHalfLife(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
	}

	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "288Mc",
		MinimumHalfLife: time.Hour,
	}, catalog)

	if result.Status != ClaimStatusStabilityIncongruent {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusStabilityIncongruent)
	}
	if !strings.Contains(result.Reason, "known half-life 170ms is shorter than claimed minimum 1h0m0s") {
		t.Fatalf("reason = %q, want half-life comparison", result.Reason)
	}
	if len(result.Evidence) != 1 || result.Evidence[0] == "" {
		t.Fatalf("evidence = %v, want citation link", result.Evidence)
	}
}

func TestValidateClaimSupportsMinimumHalfLifeWithinKnownData(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 288, HalfLife: 170 * time.Millisecond, CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
	}

	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "288Mc",
		MinimumHalfLife: 100 * time.Millisecond,
	}, catalog)

	if result.Status != ClaimStatusSupportedByTrackA {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusSupportedByTrackA)
	}
}

func TestValidateClaimReturnsInsufficientDataForUnknownIsotope(t *testing.T) {
	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "299Mc",
		MinimumHalfLife: time.Second,
	}, Catalog{})

	if result.Status != ClaimStatusInsufficientData {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusInsufficientData)
	}
}

func TestValidateClaimRejectsMismatchedCatalogIdentity(t *testing.T) {
	catalog := Catalog{
		"288Mc": {Symbol: "Mc", Z: 115, A: 290, HalfLife: 170 * time.Millisecond, CitationLink: "https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf"},
	}

	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "288Mc",
		MinimumHalfLife: 100 * time.Millisecond,
	}, catalog)

	if result.Status != ClaimStatusInvalidClaim {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusInvalidClaim)
	}
	if !strings.Contains(result.Reason, "catalog key 288Mc does not match isotope ID 290Mc") {
		t.Fatalf("reason = %q, want catalog identity mismatch", result.Reason)
	}
}

func TestValidateClaimRejectsUnsupportedMechanism(t *testing.T) {
	result := ValidateClaim(Claim{
		Kind:      ClaimKindMechanism,
		Mechanism: "antigravity propulsion",
	}, Catalog{})

	if result.Status != ClaimStatusOutsideSupportedModel {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusOutsideSupportedModel)
	}
	if !strings.Contains(result.Reason, "outside Standard Model-compatible nuclear physics") {
		t.Fatalf("reason = %q, want model-boundary explanation", result.Reason)
	}
}

func TestValidateClaimRejectsMalformedMinimumHalfLifeClaim(t *testing.T) {
	result := ValidateClaim(Claim{
		Kind:            ClaimKindMinimumHalfLife,
		IsotopeID:       "",
		MinimumHalfLife: time.Second,
	}, Catalog{})

	if result.Status != ClaimStatusInvalidClaim {
		t.Fatalf("status = %q, want %q", result.Status, ClaimStatusInvalidClaim)
	}
}
