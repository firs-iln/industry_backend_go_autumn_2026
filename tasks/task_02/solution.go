package main

func rotateRunes(s string, shift int) string {
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return ""
	}
	// shift % n never overflows: n > 0, so the remainder fits in int.
	k := shift % n
	if k < 0 {
		k += n
	}
	// Re-encode even without a shift: invalid UTF-8 becomes utf8.RuneError.
	if k == 0 {
		return string(runes)
	}
	return string(runes[k:]) + string(runes[:k])
}
