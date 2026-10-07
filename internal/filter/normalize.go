package filter

// normalize maps a raw char to 'a'..'z', or 0 if it isn't a letter-like char.
// Handles case (S -> s) and leetspeak (1 -> i, $ -> s, @ -> a ...).
func normalize(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		c += 'a' - 'A'
	}
	if c >= 'a' && c <= 'z' {
		return c
	}
	switch c {
	case '0':
		return 'o'
	case '1':
		return 'i'
	case '3':
		return 'e'
	case '4':
		return 'a'
	case '5':
		return 's'
	case '7':
		return 't'
	case '@':
		return 'a'
	case '$':
		return 's'
	}
	return 0
}

// isSeparator reports chars an abuser puts *inside* a word to dodge filters:
// s.h.i.t, s-h-i-t, s h i t
func isSeparator(c byte) bool {
	return c == '.' || c == '-' || c == '_' || c == '*' || c == ' '
}

func isWordChar(c byte) bool {
	return normalize(c) != 0
}

// atBoundary reports whether position i is the edge of a word
// (start of text / end of text / non-word char).
func atBoundary(text string, i int) bool {
	return i < 0 || i >= len(text) || !isWordChar(text[i])
}
