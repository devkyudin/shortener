package model

type Alphabet struct {
	Chars       []rune
	alphabetMap map[rune]int
}

func NewAlphabet(chars []rune) *Alphabet {
	alphabetMap := createMap(chars)
	return &Alphabet{
		Chars:       chars,
		alphabetMap: alphabetMap,
	}
}

func createMap(alphabet []rune) map[rune]int {
	result := make(map[rune]int, len(alphabet))
	for i := 0; i < len(alphabet); i++ {
		result[alphabet[i]] = i
	}

	return result
}
