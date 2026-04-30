package physics

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type EvidenceClass string

const (
	// EvidenceClassEvaluated is reserved for future Isotope.EvidenceClass labeling once the evaluated-data path is formalized; do not silently delete as unused.
	EvidenceClassEvaluated         EvidenceClass = "evaluated"
	EvidenceClassPeerReviewedModel EvidenceClass = "peer-reviewed-model"
)

type ParityClass string

const (
	ParityEvenEven ParityClass = "even-even"
	ParityOddZ     ParityClass = "odd-Z"
	ParityOddN     ParityClass = "odd-N"
	ParityOddOdd   ParityClass = "odd-odd"
)

// Coefficients of the Royer log10(T_1/2 [s]) formula:
//
//	log10(T_1/2) = A + B * Aiso^(1/6) * sqrt(Z) + C * Z / sqrt(Q_alpha)
//
// where Aiso is the parent mass number, Z is the parent atomic number, and
// Q_alpha is in MeV. Coefficients vary by parent parity class.
type Coefficients struct {
	A float64
	B float64
	C float64
}

type Model struct {
	Name          string
	Reference     string
	EvidenceClass EvidenceClass
	Coefficients  map[ParityClass]Coefficients
	Notes         string
}

type Prediction struct {
	ParityClass        ParityClass
	LogHalfLifeSeconds float64
	HalfLife           time.Duration
	Coefficients       Coefficients
	EvidenceClass      EvidenceClass
}

type AlphaResidualRecord struct {
	IsotopeID         string
	Z                 int
	A                 int
	ParityClass       ParityClass
	QAlphaMeV         float64
	EvaluatedHalfLife time.Duration
	PredictedHalfLife time.Duration
	LogResidual       float64
	FactorError       float64
	ModelName         string
	ModelReference    string
	EvidenceClass     EvidenceClass
	Skipped           bool
	SkipReason        string
}

type AlphaResidualReport struct {
	Records                 []AlphaResidualRecord
	RecordCount             int
	SkippedCount            int
	MeanAbsoluteLogResidual float64
	MaxAbsoluteLogResidual  float64
}

// RoyerModel returns the Royer-family analytic alpha-decay half-life model.
// Coefficients are the parity-class fits from the Royer 2000 / Royer & Zhang 2008
// analytic-formula family (DOI 10.1103/PhysRevC.77.037602). Output is a
// peer-reviewed-model estimate, never evaluated data.
func RoyerModel() Model {
	return Model{
		Name:          "Royer 2008 analytic alpha-decay formula",
		Reference:     "DOI 10.1103/PhysRevC.77.037602",
		EvidenceClass: EvidenceClassPeerReviewedModel,
		Notes:         "log10(T_1/2 [s]) = A + B*Aiso^(1/6)*sqrt(Z) + C*Z/sqrt(Q_alpha). Parity-class coefficients from the Royer analytic-formula family; cross-check against the Royer & Zhang 2008 PDF before relying on absolute predictions.",
		Coefficients: map[ParityClass]Coefficients{
			ParityEvenEven: {A: -25.31, B: -1.1629, C: 1.5864},
			ParityOddZ:     {A: -26.65, B: -1.0859, C: 1.5848},
			ParityOddN:     {A: -25.68, B: -1.1423, C: 1.5920},
			ParityOddOdd:   {A: -29.48, B: -1.1130, C: 1.6971},
		},
	}
}

func (m Model) Predict(z, a int, qAlphaMeV float64) (Prediction, error) {
	if z <= 0 {
		return Prediction{}, fmt.Errorf("Predict: Z must be positive, got %d", z)
	}
	if a <= 0 {
		return Prediction{}, fmt.Errorf("Predict: A must be positive, got %d", a)
	}
	if a < z {
		return Prediction{}, fmt.Errorf("Predict: A=%d must be >= Z=%d", a, z)
	}
	if math.IsNaN(qAlphaMeV) || math.IsInf(qAlphaMeV, 0) {
		return Prediction{}, errors.New("Predict: Q_alpha must be finite")
	}
	if qAlphaMeV <= 0 {
		return Prediction{}, fmt.Errorf("Predict: Q_alpha must be positive MeV, got %v", qAlphaMeV)
	}

	parity := classifyParity(z, a)
	coef, ok := m.Coefficients[parity]
	if !ok {
		return Prediction{}, fmt.Errorf("Predict: model has no coefficients for parity %q", parity)
	}

	logT := coef.A +
		coef.B*math.Pow(float64(a), 1.0/6.0)*math.Sqrt(float64(z)) +
		coef.C*float64(z)/math.Sqrt(qAlphaMeV)

	seconds := math.Pow(10, logT)
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return Prediction{}, fmt.Errorf("Predict: predicted half-life is not finite (log10=%v)", logT)
	}
	half := durationFromSeconds(seconds)

	return Prediction{
		ParityClass:        parity,
		LogHalfLifeSeconds: logT,
		HalfLife:           half,
		Coefficients:       coef,
		EvidenceClass:      m.EvidenceClass,
	}, nil
}

func AlphaResidualSummary(catalog Catalog, isotopeIDs []string, model Model) (AlphaResidualReport, error) {
	report := AlphaResidualReport{Records: make([]AlphaResidualRecord, 0, len(isotopeIDs))}
	var absoluteLogResidualSum float64

	for _, id := range isotopeIDs {
		record := AlphaResidualRecord{
			IsotopeID:      id,
			LogResidual:    math.NaN(),
			FactorError:    math.NaN(),
			ModelName:      model.Name,
			ModelReference: model.Reference,
			EvidenceClass:  model.EvidenceClass,
		}

		isotope, ok := catalog[id]
		if !ok {
			record.Skipped = true
			record.SkipReason = "isotope not in catalog"
			report.SkippedCount++
			report.Records = append(report.Records, record)
			continue
		}

		record.Z = isotope.Z
		record.A = isotope.A
		record.QAlphaMeV = isotope.QAlphaMeV
		record.EvaluatedHalfLife = isotope.HalfLife
		if isotope.QAlphaMeV <= 0 {
			record.Skipped = true
			record.SkipReason = "no positive Q_alpha in catalog"
			report.SkippedCount++
			report.Records = append(report.Records, record)
			continue
		}
		if isotope.HalfLife <= 0 {
			record.Skipped = true
			record.SkipReason = "no positive evaluated half-life in catalog"
			report.SkippedCount++
			report.Records = append(report.Records, record)
			continue
		}

		prediction, err := model.Predict(isotope.Z, isotope.A, isotope.QAlphaMeV)
		if err != nil {
			return AlphaResidualReport{}, fmt.Errorf("%s residual prediction: %w", id, err)
		}
		evaluatedSeconds := isotope.HalfLife.Seconds()
		logResidual := prediction.LogHalfLifeSeconds - math.Log10(evaluatedSeconds)
		factorError := math.Pow(10, logResidual)
		if math.IsNaN(logResidual) || math.IsInf(logResidual, 0) || math.IsNaN(factorError) || math.IsInf(factorError, 0) {
			return AlphaResidualReport{}, fmt.Errorf("%s residual is not finite", id)
		}

		record.ParityClass = prediction.ParityClass
		record.PredictedHalfLife = prediction.HalfLife
		record.LogResidual = logResidual
		record.FactorError = factorError
		report.RecordCount++
		abs := math.Abs(logResidual)
		absoluteLogResidualSum += abs
		if abs > report.MaxAbsoluteLogResidual {
			report.MaxAbsoluteLogResidual = abs
		}
		report.Records = append(report.Records, record)
	}

	if report.RecordCount > 0 {
		report.MeanAbsoluteLogResidual = absoluteLogResidualSum / float64(report.RecordCount)
	}
	return report, nil
}

func classifyParity(z, a int) ParityClass {
	zOdd := z%2 == 1
	nOdd := (a-z)%2 == 1
	switch {
	case !zOdd && !nOdd:
		return ParityEvenEven
	case zOdd && !nOdd:
		return ParityOddZ
	case !zOdd && nOdd:
		return ParityOddN
	default:
		return ParityOddOdd
	}
}

// durationFromSeconds clamps the predicted half-life into time.Duration without
// overflowing for very long-lived nuclei. Predictions above ~292 years saturate
// at math.MaxInt64 nanoseconds; the LogHalfLifeSeconds field stays exact.
func durationFromSeconds(seconds float64) time.Duration {
	const maxSeconds = float64(math.MaxInt64) / float64(time.Second)
	if seconds >= maxSeconds {
		return time.Duration(math.MaxInt64)
	}
	if seconds <= 0 {
		return 0
	}
	return time.Duration(seconds * float64(time.Second))
}
