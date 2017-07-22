package simplegostringtools

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Hline takes a character and a length n and returns a string repeating
// the character n times. The string returned has a trailing newline.
func Hline(char string, n int) string {
	line := ""
	for i := 0; i < n; i++ {
		line = line + char
	}
	return line + "\n"
}

// Sep returns a string repeating the character "-" 75 times which can be
// used as a visual separator in other strings/output. The returned string
// has a trailing newline.
func Sep() string {
	return Hline("-", 75)
}

// Reverse takes a string and and returnes a string wich is the rune-exact
// reverse representation of the input string.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Frame takes a string and a character and returns a string which is the
// input string framed by the character.
func Frame(s string, char string) string {
	lines := strings.Split(s, "\n")
	max := 0
	length := 0
	framed := ""
	hline := ""
	margin := 1
	filler := " "
	for _, line := range lines {
		length = len(line)
		if length > max {
			max = length
		}
	}
	// Make the horizontal line for the top and the bottom.
	for i := 0; i < max+2*margin+2*len(char); i++ {
		hline = fmt.Sprintf("%s%s", hline, char)
	}
	// Make the inner part of the framed string without top and bottom.
	hfill := ""
	for _, line := range lines {
		hfill = ""
		length = len(line)
		// Fill up the right side of each line, if necessary.
		if length < max {
			for i := 0; i < max-length; i++ {
				hfill = fmt.Sprintf("%s%s", hfill, filler)
			}
		}
		line = fmt.Sprintf("%s%s", line, hfill)
		// Compose the line.
		framed = fmt.Sprintf("%s\n%s%s%s%s%s",
			framed, char, filler, line, filler, char)
	}
	// Add the top and bottom margin.
	hmargin := ""
	for i := 0; i < max+2*margin; i++ {
		hmargin = fmt.Sprintf("%s%s", hmargin, filler)
	}
	hmargin = fmt.Sprintf("\n%s%s%s", char, hmargin, char)
	framed = fmt.Sprintf("%s%s%s", hmargin, framed, hmargin)
	// Add the top and bottom horizontal lines.
	framed = fmt.Sprintf("%s%s\n%s", hline, framed, hline)
	return framed
}

// Studly takes a string and returns it in studly case format. The function
// stricly alternates case on every character, so even whitespace counts
// as character when alternating between upper case and lower case.
func Studly(s string) string {
	studly := ""
	for i, c := range s {
		if i%2 == 0 {
			studly = fmt.Sprintf("%s%s", studly,
				strings.ToUpper(string(c)))
		} else {
			studly = fmt.Sprintf("%s%s", studly, string(c))
		}
	}
	return studly
}

// TODO: Document.
// TODO: Test.
func RandomStrings(n int) string {
	s := ""
	word := ""
	ascii := []int{
		97,  // a
		98,  // b
		99,  // c
		100, // d
		101, // e
		102, // f
		103, // g
		104, // h
		105, // i
		106, // j
		107, // k
		108, // l
		109, // m
		110, // n
		111, // o
		112, // p
		113, // q
		114, // r
		115, // s
		116, // t
		117, // u
		118, // v
		119, // w
		120, // x
		121, // y
		122, // z
	}
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < n; i++ {
		for j := 0; j < rand.Intn(10); j++ {
			word = fmt.Sprintf("%s%s", word,
				string(ascii[rand.Intn(len(ascii))]))
		}
		s = fmt.Sprintf("%s %s", s, word)
		word = ""
	}
	return s
}
