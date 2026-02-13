package runes

func Next(ch rune) rune {
	return rune(ch+1-'A')%('Z'-'A'+1) + 'A'
}
