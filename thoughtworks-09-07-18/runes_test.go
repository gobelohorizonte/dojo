package main

import (

		"testing"
		"strings"
		"os"
		)

func Example() {
	main()
	// Output:
	// Please provide one or more words to search.
}

func TestParseLine(t *testing.T) {
  line := "0059;LATIN CAPITAL LETTER Y;Lu;0;L;;;;;N;;;;0079;"
  gotChar, gotName := parseLine(line)
  wantChar := 'Y'
  wantName := "LATIN CAPITAL LETTER Y"
  if gotChar != wantChar {
  	t.Errorf("Got: %q want: %q", gotChar, wantChar)
  }

  if gotName != wantName {
  	t.Errorf("Got name: %q want: %q", gotChar, wantName)
  }
}

const sample = `002A;ASTERISK;Po;0;ON;;;;;N;;;;;
002B;PLUS SIGN;Sm;0;ES;;;;;N;;;;;
002C;COMMA;Po;0;CS;;;;;N;;;;;
002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;
002E;FULL STOP;Po;0;CS;;;;;N;PERIOD;;;;
002F;SOLIDUS;Po;0;CS;;;;;N;SLASH;;;;
0030;DIGIT ZERO;Nd;0;EN;;0;0;0;N;;;;;
0031;DIGIT ONE;Nd;0;EN;;1;1;1;N;;;;;`

func TestFilter(t *testing.T) {
	query := "solidus"
	got := Filter(strings.NewReader(sample), query)
	if len(got) != 1 {
		t.Errorf("expected slice of len 1, got: %v", got)
	}
}

func ExampleSingleResult() {
	oldArgs := os.Args
	defer func() {os.Args = oldArgs}()
	os.Args = []string{"", "registered", "sign"}
	main()
	// Output: 
	// U+00AE	®	REGISTERED SIGN
	
}

func TestFormat(t *testing.T) {
	charName := CharName{'A', "LATIN CAPITAL LETTER A"}
	got := format(charName)
	want := "U+0041\tA\tLATIN CAPITAL LETTER A"
	if got != want {
		t.Errorf("Got: %q want: %q", got, want)
	}
}














