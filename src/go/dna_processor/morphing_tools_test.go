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
	if err != Finish {
		t.Errorf("Expected nat(\"\") to throw 'Finish' exception")
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
	cases := []struct {
		name           string
		dna            DnaStorage
		expected       string
		expectedDnaLen int
	}{
		{
			"quote test 00",
			NewSimpleDnaStorage("CFPICIIF"),
			"ICFP",
			3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			received, err := consts(tc.dna)
			if err != nil {
				t.Fatalf("Error occurred: %v", err)
			}

			if received != tc.expected || tc.dna.Len() != tc.expectedDnaLen {
				t.Errorf("Fail: expected %s unquoted received %s",
					tc.expected, received)
			}
		})
	}
}

func TestPattern(t *testing.T) {
	cases := []struct {
		name             string
		dna              DnaStorage
		expected_pattern string
		expected_dna_len int
	}{
		{
			name:             "SingleToken_00",
			dna:              NewSimpleDnaStorage("CIIC"),
			expected_pattern: "I",
			expected_dna_len: 0,
		},
		{
			name:             "SingleToken_01",
			dna:              NewSimpleDnaStorage("FIIC"),
			expected_pattern: "C",
			expected_dna_len: 0,
		},
		{
			name:             "SingleToken_02",
			dna:              NewSimpleDnaStorage("PIIC"),
			expected_pattern: "F",
			expected_dna_len: 0,
		},
		{
			name:             "SingleToken_03",
			dna:              NewSimpleDnaStorage("ICIIC"),
			expected_pattern: "P",
			expected_dna_len: 0,
		},
		{
			name:             "MultipleTokens_00",
			dna:              NewSimpleDnaStorage("IIPIPICPIICICIIF"),
			expected_pattern: "(!(2))P",
			expected_dna_len: 0,
		},
		{
			name:             "MultipleTokens_01",
			dna:              NewSimpleDnaStorage("IPCICPIIF"),
			expected_pattern: "!(5)",
			expected_dna_len: 0,
		},
		{
			name:             "MultipleTokens_02",
			dna:              NewSimpleDnaStorage("IFICIIC"),
			expected_pattern: "?(I)",
			expected_dna_len: 0,
		},
		{
			name:             "MultipleTokens_03",
			dna:              NewSimpleDnaStorage("IFIFIIC"),
			expected_pattern: "?(C)",
			expected_dna_len: 0,
		},
		{
			name:             "MultipleTokens_04",
			dna:              NewSimpleDnaStorage("IFIPIIC"),
			expected_pattern: "?(F)",
			expected_dna_len: 0,
		},
		{
			name:             "MultipleTokens_05",
			dna:              NewSimpleDnaStorage("IFIICIIC"),
			expected_pattern: "?(P)",
			expected_dna_len: 0,
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

			if sp != tc.expected_pattern || tc.dna.Len() != tc.expected_dna_len {
				t.Errorf("Fail: expected %s pattern with length %d received %s with length %d",
					tc.expected_pattern, tc.expected_dna_len, sp, tc.dna.Len())
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
			env, err := match(tc.dna, tc.pattern, false)

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
			err := Step(dna, 0, false)
			if err != nil {
				t.Errorf("Fail: error %v", err)
			}

			if dna.String() != tc.expected {
				t.Errorf("Fail: expected %s received %s", tc.expected, dna.String())
			}
		})
	}
}
