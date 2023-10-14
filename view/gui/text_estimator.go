package gui

var LetterWidth map[rune]float64 = map[rune]float64{
	'a':  0.5,
	'b':  0.5,
	'c':  0.5,
	'd':  0.5,
	'e':  0.6,
	'f':  0.3,
	'g':  0.5,
	'h':  0.5,
	'i':  0.2,
	'j':  0.4,
	'k':  0.4,
	'l':  0.2,
	'm':  0.6,
	'n':  0.5,
	'o':  0.6,
	'p':  0.5,
	'q':  0.4,
	'r':  0.3,
	's':  0.5,
	't':  0.3,
	'u':  0.5,
	'v':  0.4,
	'w':  0.7,
	'x':  0.4,
	'y':  0.4,
	'z':  0.4,
	'A':  0.8,
	'B':  0.8,
	'C':  0.8,
	'D':  0.8,
	'E':  0.9,
	'F':  0.8,
	'G':  0.8,
	'H':  0.8,
	'I':  0.4,
	'J':  0.8,
	'K':  0.8,
	'L':  0.8,
	'M':  1.0,
	'N':  0.8,
	'O':  0.8,
	'P':  0.9,
	'Q':  0.8,
	'R':  0.8,
	'S':  0.8,
	'T':  0.8,
	'U':  0.8,
	'V':  0.9,
	'W':  1.1,
	'X':  0.8,
	'Y':  0.8,
	'Z':  0.8,
	'\'': 0.2,
	' ':  0.2,
}

func EstimateWidth(text string) float64 {
	var result = 0.0
	for _, c := range text {
		if w, ok := LetterWidth[c]; ok {
			result += w
		} else {
			result += 0.5
		}
	}
	return result
}
