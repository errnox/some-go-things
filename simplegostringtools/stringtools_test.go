package simplegostringtools

import (
	"testing"
)

func TestHline(t *testing.T) {
	tests := []struct {
		Char string
		N    int
		Out  string
	}{
		{"=", 10, "==========\n"},
		{".", 5, ".....\n"},
		{"*", 75, "***************************************************************************\n"},
		{"a", 20, "aaaaaaaaaaaaaaaaaaaa\n"},
		{"b", 21, "bbbbbbbbbbbbbbbbbbbbb\n"},
		{"ab", 10, "abababababababababab\n"},
	}
	for _, test := range tests {
		want := test.Out
		got := Hline(test.Char, test.N)
		if got != want {
			t.Errorf("Hline(\"%v\", %v) = \"%v\", want \"%v\"",
				test.Char, test.N, want)
		}
	}

}

func TestSep(t *testing.T) {
	want := "---------------------------------------------------------------------------\n"
	got := Sep()
	if got != want {
		t.Errorf("Sep() = \"%v\", want \"%v\"", want, got)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		S   string
		Out string
	}{
		{"Hello, world!", "!dlrow ,olleH"},
		{"!dlrow ,olleH", "Hello, world!"},
		{"123", "321"},
		{"abcd - 123", "321 - dcba"},
	}
	for _, test := range tests {
		want := test.Out
		got := Reverse(test.S)
		if got != want {
			t.Errorf("Reverse(\"%v\") = \"%v\", want \"%v\"",
				test.S, test.Out, want)
		}
	}
}

func TestFrame(t *testing.T) {
	tests := []struct {
		Str  string
		Char string
		Out  string
	}{
		{"This is a string.", "#", "#####################\n#                   #\n# This is a string. #\n#                   #\n#####################"},
		{"This is another string.\n\nIs it framed?", "*", "***************************\n*                         *\n* This is another string. *\n*                         *\n* Is it framed?           *\n*                         *\n***************************"},
	}
	for _, test := range tests {
		want := test.Out
		got := Frame(test.Str, test.Char)
		if got != want {
			t.Errorf("Frame(\"%s\", \"%s\") = \"%s\", want \"%s\"",
				test.Str, test.Char, got, want)
		}
	}
}

func TestStudly(t *testing.T) {
	tests := []struct {
		In, Out string
	}{
		{"This is a test.", "ThIs iS A TeSt."},
		{"Hello", "HeLlO"},
		{"StudlyCaps", "StUdLyCaPs"},
	}
	for _, test := range tests {
		want := test.Out
		got := Studly(test.In)
		if got != want {
			t.Errorf("Studly(\"%s\") = \"%s\", want \"%s\"",
				test.In, got, want)
		}
	}

}
