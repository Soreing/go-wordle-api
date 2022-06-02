package wordle

import(
	"strings"
	"bytes"
)

// Creates a string that contains the letters of the guess
// and the colors compared to a solution in wordle
func ColorWord(guess string, solution string) string{
	guessUpper := strings.ToUpper(guess)
	solutionUpper := strings.ToUpper(solution)
	
	source := make([]byte, 5)	// Contains letters from the solution
	letters := make([]byte, 5)	// Contains letters of the guess
	colors := make([]byte, 5)	// Contains colors of the guess letters

	// Determine green blocks
	for i:=0; i<5; i++ {
		letters[i] = guessUpper[i]
		if guessUpper[i] == solutionUpper[i] {
			source[i] = ' '
			colors[i] = 'g'
		} else {
			source[i] = solutionUpper[i]
			colors[i] = 'b'
		}
	}

	// Determine yellow blocks
	for i:=0; i<5; i++ {
		if colors[i] == 'b' {
			if idx := bytes.IndexByte(source, guessUpper[i]); idx >= 0 {
				colors[i] = 'y'
				source[i] = ' '
			}
		}	
	}

	// Create Result and return it
	// Couples of letters and colors ("yA" for yellow A)
	result := make([]byte, 15)
	for i:=0; i<5; i++ {
		result[i*3] = colors[i]
		result[i*3+1] = letters[i]
		result[i*3+2] = ' '
	}

	return string(result[:])
}