package research

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseENSDFNuclideHeader extracts nuclide ID from an ENSDF header line.
// Input format: "288MC    ADOPTED LEVELS, GAMMAS"
// Returns the normalized nuclide ID like "288Mc".
func ParseENSDFNuclideHeader(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return "", fmt.Errorf("empty ENSDF header line")
	}
	raw := strings.ToUpper(fields[0])

	i := 0
	for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
		i++
	}
	if i == 0 || i == len(raw) {
		return "", fmt.Errorf("invalid ENSDF nuclide header: %q", line)
	}

	mass := raw[:i]
	symbol := strings.Title(strings.ToLower(raw[i:]))
	return mass + symbol, nil
}

// ParseENSDFHalfLifeField extracts half-life in seconds from an ENSDF T field.
// Input format: "  T  1.700E-01 S"
// Unit suffixes: S (seconds), MS, US, NS, PS, M (minutes), H (hours), D (days), Y (years).
func ParseENSDFHalfLifeField(line string) (float64, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 || fields[0] != "T" {
		return 0, fmt.Errorf("invalid ENSDF half-life line: %q", line)
	}
	val, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, fmt.Errorf("parse half-life value %q: %w", fields[1], err)
	}
	switch strings.ToUpper(fields[2]) {
	case "S":
		return val, nil
	case "MS":
		return val * 1e-3, nil
	case "US":
		return val * 1e-6, nil
	case "NS":
		return val * 1e-9, nil
	case "PS":
		return val * 1e-12, nil
	case "M":
		return val * 60, nil
	case "H":
		return val * 3600, nil
	case "D":
		return val * 86400, nil
	case "Y":
		return val * 31556952, nil
	default:
		return 0, fmt.Errorf("unknown ENSDF time unit %q", fields[2])
	}
}

// ParseENSDFQValueField extracts Q-value in MeV from an ENSDF Q field.
// Input format: " Q  1.075E+04 KEV" or " Q  10.75 MEV"
func ParseENSDFQValueField(line string) (float64, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 || fields[0] != "Q" {
		return 0, fmt.Errorf("invalid ENSDF Q-value line: %q", line)
	}
	val, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, fmt.Errorf("parse Q-value %q: %w", fields[1], err)
	}
	switch strings.ToUpper(fields[2]) {
	case "KEV":
		return val * 1e-3, nil
	case "MEV":
		return val, nil
	case "EV":
		return val * 1e-6, nil
	default:
		return val * 1e-3, nil // assume keV
	}
}

// ParseENSDFDaughterField extracts daughter nuclide ID from ENSDF DA field.
// Input format: " DA  284NH"
func ParseENSDFDaughterField(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 || fields[0] != "DA" {
		return "", fmt.Errorf("invalid ENSDF daughter line: %q", line)
	}
	raw := strings.ToUpper(fields[1])

	i := 0
	for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
		i++
	}
	if i == 0 || i == len(raw) {
		return "", fmt.Errorf("invalid ENSDF daughter entry: %q", line)
	}

	mass := raw[:i]
	symbol := strings.Title(strings.ToLower(raw[i:]))
	return mass + symbol, nil
}
