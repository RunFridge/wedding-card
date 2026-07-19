package moderation

import (
	"testing"
)

func TestIsFlaggedByOpenAI(t *testing.T) {
	w := &Worker{}
	result := &ModerationResult{Flagged: true}

	if !w.isFlagged(result) {
		t.Error("expected flagged=true when result.Flagged is true")
	}
}

func TestIsFlaggedNotFlagged(t *testing.T) {
	w := &Worker{}
	result := &ModerationResult{Flagged: false}

	if w.isFlagged(result) {
		t.Error("expected flagged=false when result.Flagged is false and no thresholds")
	}
}

func TestIsFlaggedByThreshold(t *testing.T) {
	w := &Worker{
		ThresholdGetter: func() map[string]float64 {
			return map[string]float64{
				"violence": 0.5,
				"sexual":   0.8,
			}
		},
	}

	result := &ModerationResult{
		Flagged: false,
		CategoryScores: map[string]float64{
			"violence": 0.6,
			"sexual":   0.1,
		},
	}

	if !w.isFlagged(result) {
		t.Error("expected flagged=true when violence score exceeds threshold")
	}
}

func TestIsFlaggedByThresholdExactMatch(t *testing.T) {
	w := &Worker{
		ThresholdGetter: func() map[string]float64 {
			return map[string]float64{
				"violence": 0.5,
			}
		},
	}

	result := &ModerationResult{
		Flagged: false,
		CategoryScores: map[string]float64{
			"violence": 0.5,
		},
	}

	if !w.isFlagged(result) {
		t.Error("expected flagged=true when score equals threshold (>=)")
	}
}

func TestIsFlaggedBelowThreshold(t *testing.T) {
	w := &Worker{
		ThresholdGetter: func() map[string]float64 {
			return map[string]float64{
				"violence": 0.5,
				"sexual":   0.8,
			}
		},
	}

	result := &ModerationResult{
		Flagged: false,
		CategoryScores: map[string]float64{
			"violence": 0.4,
			"sexual":   0.7,
		},
	}

	if w.isFlagged(result) {
		t.Error("expected flagged=false when all scores below thresholds")
	}
}

func TestIsFlaggedNilThresholdGetter(t *testing.T) {
	w := &Worker{ThresholdGetter: nil}
	result := &ModerationResult{
		Flagged: false,
		CategoryScores: map[string]float64{
			"violence": 0.9,
		},
	}

	if w.isFlagged(result) {
		t.Error("expected flagged=false when ThresholdGetter is nil and Flagged=false")
	}
}

func TestIsFlaggedUnknownCategory(t *testing.T) {
	w := &Worker{
		ThresholdGetter: func() map[string]float64 {
			return map[string]float64{
				"violence": 0.5,
			}
		},
	}

	result := &ModerationResult{
		Flagged: false,
		CategoryScores: map[string]float64{
			"unknown_category": 0.9,
		},
	}

	if w.isFlagged(result) {
		t.Error("expected flagged=false when score category has no configured threshold")
	}
}
