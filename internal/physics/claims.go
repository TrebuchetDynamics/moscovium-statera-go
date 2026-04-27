package physics

import (
	"fmt"
	"time"
)

type ClaimKind string

const (
	ClaimKindMinimumHalfLife ClaimKind = "minimum-half-life"
	ClaimKindMechanism       ClaimKind = "mechanism"
)

type ClaimStatus string

const (
	ClaimStatusSupportedByTrackA     ClaimStatus = "supported-by-track-a"
	ClaimStatusStabilityIncongruent  ClaimStatus = "stability-incongruent"
	ClaimStatusOutsideSupportedModel ClaimStatus = "outside-supported-model"
	ClaimStatusInsufficientData      ClaimStatus = "insufficient-data"
	ClaimStatusInvalidClaim          ClaimStatus = "invalid-claim"
)

type Claim struct {
	Kind            ClaimKind
	IsotopeID       string
	MinimumHalfLife time.Duration
	Mechanism       string
}

type ClaimResult struct {
	Status   ClaimStatus
	Reason   string
	Evidence []string
}

func ValidateClaim(claim Claim, catalog Catalog) ClaimResult {
	if result, ok := validateClaimShape(claim); !ok {
		return result
	}

	switch claim.Kind {
	case ClaimKindMinimumHalfLife:
		return validateMinimumHalfLifeClaim(claim, catalog)
	case ClaimKindMechanism:
		return validateMechanismClaim(claim)
	default:
		return ClaimResult{
			Status: ClaimStatusInvalidClaim,
			Reason: fmt.Sprintf("unknown claim kind %q", claim.Kind),
		}
	}
}

func validateClaimShape(claim Claim) (ClaimResult, bool) {
	switch claim.Kind {
	case ClaimKindMinimumHalfLife:
		if claim.Mechanism != "" {
			return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: "claim mixes multiple assertions"}, false
		}
	case ClaimKindMechanism:
		if claim.IsotopeID != "" || claim.MinimumHalfLife != 0 {
			return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: "claim mixes multiple assertions"}, false
		}
	}
	return ClaimResult{}, true
}

func validateMinimumHalfLifeClaim(claim Claim, catalog Catalog) ClaimResult {
	if claim.IsotopeID == "" {
		return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: "minimum half-life claim requires isotope ID"}
	}
	if claim.MinimumHalfLife <= 0 {
		return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: "minimum half-life claim requires positive duration"}
	}

	isotope, ok := catalog[claim.IsotopeID]
	if !ok {
		return ClaimResult{
			Status: ClaimStatusInsufficientData,
			Reason: fmt.Sprintf("isotope %s is not present in Track A catalog", claim.IsotopeID),
		}
	}
	if isotope.ID() != claim.IsotopeID {
		return ClaimResult{
			Status: ClaimStatusInvalidClaim,
			Reason: fmt.Sprintf("catalog key %s does not match isotope ID %s", claim.IsotopeID, isotope.ID()),
		}
	}
	if err := isotope.Validate(); err != nil {
		return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: err.Error()}
	}
	if isotope.HalfLife <= 0 {
		return ClaimResult{
			Status:   ClaimStatusInsufficientData,
			Reason:   fmt.Sprintf("isotope %s has no evaluated half-life in Track A catalog", claim.IsotopeID),
			Evidence: []string{isotope.CitationLink},
		}
	}

	if isotope.HalfLife < claim.MinimumHalfLife {
		return ClaimResult{
			Status:   ClaimStatusStabilityIncongruent,
			Reason:   fmt.Sprintf("%s known half-life %s is shorter than claimed minimum %s", claim.IsotopeID, isotope.HalfLife, claim.MinimumHalfLife),
			Evidence: []string{isotope.CitationLink},
		}
	}

	return ClaimResult{
		Status:   ClaimStatusSupportedByTrackA,
		Reason:   fmt.Sprintf("%s known half-life %s meets claimed minimum %s", claim.IsotopeID, isotope.HalfLife, claim.MinimumHalfLife),
		Evidence: []string{isotope.CitationLink},
	}
}

func validateMechanismClaim(claim Claim) ClaimResult {
	if claim.Mechanism == "" {
		return ClaimResult{Status: ClaimStatusInvalidClaim, Reason: "mechanism claim requires mechanism text"}
	}
	return ClaimResult{
		Status: ClaimStatusOutsideSupportedModel,
		Reason: fmt.Sprintf("%q is outside Standard Model-compatible nuclear physics implemented by this project", claim.Mechanism),
	}
}
