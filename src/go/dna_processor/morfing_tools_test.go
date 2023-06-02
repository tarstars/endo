package dna_processor

import (
	"testing"
)

func TestNat(t *testing.T) {
	// Test case 1: "ICFP"
	result, err := nat(NewSimpleDnaStorage("ICFP"))
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	expected := 2
	if result != expected {
		t.Errorf("Expected nat(\"ICFP\") to be %d, but got %d", expected, result)
	}

	// Test case 2: "ICFFFFP"
	result, err = nat(NewSimpleDnaStorage("ICFFFFP"))
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	expected = 2
	if result != expected {
		t.Errorf("Expected nat(\"ICFFFFP\") to be %d, but got %d", expected, result)
	}

	// Test case 3: "P"
	result, err = nat(NewSimpleDnaStorage("P"))
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	expected = 0
	if result != expected {
		t.Errorf("Expected nat(\"P\") to be %d, but got %d", expected, result)
	}

	// Test case 4: Empty sequence
	_, err = nat(NewSimpleDnaStorage(""))
	if err != finish {
		t.Errorf("Expected nat(\"\") to throw 'finish' exception")
	}

	// Test case 5: Invalid letter
	_, err = nat(NewSimpleDnaStorage("X"))
	if err == nil || err.Error() != "invalid letter" {
		t.Errorf("Expected nat(\"X\") to throw 'invalid letter' error")
	}

	// Test case 6: Missing end of sequence
	_, err = nat(NewSimpleDnaStorage("ICF"))
	if err == nil || err.Error() != "missing end of sequence" {
		t.Errorf("Expected nat(\"ICF\") to throw 'missing end of sequence' error")
	}
}

func TestConsts(t *testing.T) {
	// Test case 1: "CFPIF"
	result, err := consts(NewSimpleDnaStorage("CFPIF"))
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	expected := "ICF"
	if result != expected {
		t.Errorf("Expected consts(\"CFPIF\") to be %s, but got %s", expected, result)
	}

	// Test case 2: "ICFPC"
	result, err = consts(NewSimpleDnaStorage("ICFPC"))
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	expected = "PCFI"
	if result != expected {
		t.Errorf("Expected consts(\"ICFPC\") to be %s, but got %s", expected, result)
	}

	// Test case 3: "ICF"
	result, err = consts(NewSimpleDnaStorage("ICF"))
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	expected = "PC"
	if result != expected {
		t.Errorf("Expected consts(\"ICF\") to be %s, but got %s", expected, result)
	}

	// Test case 4: "X"
	_, err = consts(NewSimpleDnaStorage("X"))
	if err == nil || err.Error() != "invalid letter" {
		t.Errorf("Expected consts(\"X\") to throw 'invalid letter' error")
	}
}

func TestPattern(t *testing.T) {
	cases := []struct {
		name     string
		dna      DnaStorage
		expected string
	}{
		{
			name:     "SingleToken",
			dna:      NewSimpleDnaStorage("CIIC"),
			expected: "I",
		},
		{
			name:     "MultipleTokens",
			dna:      NewSimpleDnaStorage("IIPIPICPIICICIIF"),
			expected: "(!(2))P",
		},
		// Add more test cases here as needed.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := pattern(tc.dna)
			if err != nil {
				t.Fatalf("Error occurred: %v", err)
			}

			sp := patternToString(p)

			if sp != tc.expected {
				t.Errorf("Fail: expected %s received %s", tc.expected, sp)
			}
		})
	}
}

func TestPatternToString(t *testing.T) {
	cases := []struct {
		name     string
		pattern  []PatternToken
		expected string
	}{
		{
			name:     "EmptyPattern",
			pattern:  []PatternToken{},
			expected: "",
		},
		{
			name: "PatternWithSingleConstToken",
			pattern: []PatternToken{
				NewConstToken('C'),
			},
			expected: "C",
		},
		{
			name: "PatternWithSkipAndSearchTokens",
			pattern: []PatternToken{
				NewSkipToken(4),
				NewSearchToken("IFPC"),
			},
			expected: "!(4)?(IFPC)",
		},
		{
			name: "PatternWithBraAndKetTokens",
			pattern: []PatternToken{
				NewBraToken(),
				NewConstToken('F'),
				NewKetToken(),
			},
			expected: "(F)",
		},
		{
			name: "PatternWithBraAndKetTokens",
			pattern: []PatternToken{
				NewBraToken(),
				NewSkipToken(2),
				NewKetToken(),
				NewConstToken('P'),
			},
			expected: "(!(2))P",
		},
		// Add more test cases as needed
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := patternToString(tc.pattern)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestTemplate(t *testing.T) {
	cases := []struct {
		name     string
		dna      DnaStorage
		expected string
	}{
		{
			name:     "Sub_00",
			dna:      NewSimpleDnaStorage("ICCIFPPIIC"),
			expected: "PI<0,0>",
		},
		// Add more test cases here as needed.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			template, err := template(tc.dna)
			if err != nil {
				t.Fatalf("Error occurred: %v", err)
			}

			sp := templateToString(template)

			if sp != tc.expected {
				t.Errorf("Fail: expected %s received %s", tc.expected, sp)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	cases := []struct {
		name     string
		dna      DnaStorage
		pattern  []PatternToken
		expected Environment
	}{
		{
			name: "Match_00",
			dna:  NewSimpleDnaStorage("CFPC"),
			pattern: []PatternToken{
				NewBraToken(),
				NewSkipToken(2),
				NewKetToken(),
				NewConstToken('P'),
			},
			expected: Environment{"CF"},
		},
		// Add more test cases here as needed.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env, err := match(tc.dna, tc.pattern)

			if err != nil {
				t.Fatalf("Errof: %v", err)
			}

			if env.NotEqual(tc.expected) {
				t.Errorf("Fail: expected %s received %s", tc.expected, env)
			}
		})
	}
}

func TestApplyTemplate(t *testing.T) {
	cases := []struct {
		name        string
		environment Environment
		template    []TemplateToken
		expected    string
	}{
		{
			name:        "Match_00",
			environment: []string{"CF"},
			template: []TemplateToken{
				NewConstToken('P'),
				NewConstToken('I'),
				NewReferenceToken(0, 0),
			},
			expected: "ICF",
		},
		// Add more test cases here as needed.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newPrefix, err := formPrefix(tc.template, tc.environment)

			if err != nil {
				t.Fatalf("Errof: %v", err)
			}

			if newPrefix == tc.expected {
				t.Errorf("Fail: expected %s received %s", tc.expected, newPrefix)
			}
		})
	}
}

func TestAsNat(t *testing.T) {
	cases := []struct {
		name      string
		num       int
		stringNum string
	}{
		{
			name:      "one",
			num:       1,
			stringNum: "CP",
		},
		{
			name:      "zero",
			num:       0,
			stringNum: "P",
		},
		{
			name:      "eight",
			num:       8,
			stringNum: "IIICP",
		},
		// Add more test cases here as needed.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			strNum := AsNat(tc.num)

			if strNum != tc.stringNum {
				t.Errorf("Fail: expected %s received %s", tc.stringNum, strNum)
			}
		})
	}
}

func TestProtect(t *testing.T) {
	cases := []struct {
		name      string
		level     int
		toProtect string
		expected  string
	}{
		{
			name:      "case 00",
			level:     1,
			toProtect: "ICFP",
			expected:  "CFPIC",
		},
		{
			name:      "case 01",
			level:     2,
			toProtect: "ICFP",
			expected:  "FPICCF",
		},
		// Add more test cases here as needed.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			received := protect(tc.level, tc.toProtect)

			if received != tc.expected {
				t.Errorf("Fail: expected %s received %s", tc.expected, received)
			}
		})
	}
}

func TestStep(t *testing.T) {
	cases := []struct {
		name      string
		sourceDna string
		expected  string
	}{
		{
			name:      "case 00",
			sourceDna: "IIPIPICPIICICIIFICCIFPPIICCFPC",
			expected:  "PICFC",
		},
		{
			name:      "case 01",
			sourceDna: "IIPIPICPIICICIIFICCIFCCCPPIICCFPC",
			expected:  "PIICCFCFFPC",
		},
		{
			name:      "case 02",
			sourceDna: "IIPIPIICPIICIICCIICFCFC",
			expected:  "I",
		},
		// Add more test cases here as needed.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dna := NewSimpleDnaStorage(tc.sourceDna)
			err := Step(dna)
			if err != nil {
				t.Errorf("Fail: error %v", err)
			}

			if dna.String() != tc.expected {
				t.Errorf("Fail: expected %s received %s", tc.expected, dna.String())
			}
		})
	}
}
