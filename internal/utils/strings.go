package utils

import "strings"

func normalizeStringForLooseComparison(s string) string {
	lowerS := strings.ToLower(s)
	return strings.ReplaceAll(lowerS, " ", "")
}

// eg: LooselyCompareTwoStrings("James", "jam") => true
func LooselyCompareTwoStrings(masterString, candidateString string) bool {
	mS := normalizeStringForLooseComparison(masterString)
	cS := normalizeStringForLooseComparison(candidateString)
	if len(mS) == 0 || len(cS) == 0 {
		return false
	}
	return strings.Contains(mS, cS)
}
