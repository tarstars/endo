package dna_processor

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

// Finish Custom exception for signaling the end of sequence
var Finish = errors.New("Finish")

// Function to interpret a sequence of letters 'I', 'C', 'F', 'P'
func nat(s DnaStorage) (int, error) {
	num := 0
	k := 1

	if s.IsEmpty() {
		return 0, Finish
	}

	for !s.IsEmpty() {
		c := s.GetChar()
		if c == 'C' {
			num += k
		} else if c == 'P' {
			return num, nil
		} else if !(c == 'I' || c == 'F') {
			return 0, errors.New("invalid letter")
		}
		k *= 2
	}

	return 0, errors.New("missing end of sequence")
}

func consts(s DnaStorage) (string, error) {
	buf := bytes.Buffer{}
	for {
		if s.IsEmpty() {
			return buf.String(), nil
		}
		c := s.GetChar()
		switch c {
		case 'C':
			buf.WriteByte('I')
		case 'F':
			buf.WriteByte('C')
		case 'P':
			buf.WriteByte('F')
		case 'I':
			if s.IsEmpty() {
				return buf.String(), nil
			}
			cc := s.GetChar()
			if cc == 'C' {
				buf.WriteByte('P')
			} else {
				s.UndoGet()
				s.UndoGet()
				return buf.String(), nil
			}
		default:
			return "", errors.New("invalid letter")
		}
	}
}

type PatternToken interface {
	getType() string
}

type TemplateToken interface {
	getType() string
}

type ConstToken struct {
	c byte
}

func (t ConstToken) getType() string {
	return "ConstToken"
}

func NewConstToken(c byte) *ConstToken {
	return &ConstToken{c}
}

type SkipToken struct {
	n int
}

func (t SkipToken) getType() string {
	return "SkipToken"
}

func NewSkipToken(n int) *SkipToken {
	return &SkipToken{n}
}

type SearchToken struct {
	s string
}

func NewSearchToken(s string) *SearchToken {
	return &SearchToken{s}
}

func (t SearchToken) getType() string {
	return "SearchToken"
}

type BraToken struct{}

func (t BraToken) getType() string {
	return "BraToken"
}

func NewBraToken() *BraToken {
	return &BraToken{}
}

type KetToken struct{}

func (t KetToken) getType() string {
	return "KetToken"
}

func NewKetToken() *KetToken {
	return &KetToken{}
}

type ReferenceToken struct {
	n, l int
}

func (t ReferenceToken) getType() string {
	return "ReferenceToken"
}

func NewReferenceToken(n, l int) *ReferenceToken {
	return &ReferenceToken{n, l}
}

type LenToken struct {
	n int
}

func (t LenToken) getType() string {
	return "LenToken"
}

func NewLenToken(n int) *LenToken {
	return &LenToken{n}
}

func pattern(s DnaStorage) ([]PatternToken, error) {
	tokens := make([]PatternToken, 0)
	lvl := 0

	for {
		if s.IsEmpty() {
			return nil, Finish
		}
		c := s.GetChar()
		switch c {
		case 'C':
			tokens = append(tokens, NewConstToken('I'))
		case 'F':
			tokens = append(tokens, NewConstToken('C'))
		case 'P':
			tokens = append(tokens, NewConstToken('F'))
		case 'I':
			if s.IsEmpty() {
				return nil, Finish
			}
			cc := s.GetChar()
			switch cc {
			case 'C':
				tokens = append(tokens, NewConstToken('P'))
			case 'P':
				if n, err := nat(s); err != nil {
					return nil, err
				} else {
					tokens = append(tokens, NewSkipToken(n))
				}
			case 'F':
				if s.IsEmpty() {
					return nil, Finish
				}
				s.GetChar()
				if substring, err := consts(s); err == nil {
					tokens = append(tokens, NewSearchToken(substring))
				}
			case 'I':
				if s.IsEmpty() {
					return nil, Finish
				}
				ccc := s.GetChar()
				switch ccc {
				case 'P':
					tokens = append(tokens, NewBraToken())
					lvl += 1
				case 'C', 'F':
					if lvl == 0 {
						return tokens, nil
					} else {
						lvl -= 1
						tokens = append(tokens, NewKetToken())
					}
				case 'I':
					for p := 3; p < 10; p += 1 {
						fmt.Print(string(s.GetChar()))
					}
				}
			}
		}
	}
}

func patternToString(pattern []PatternToken) string {
	var buf bytes.Buffer

	for _, token := range pattern {
		switch t := token.(type) {
		case *ConstToken:
			buf.WriteString(string(t.c))
		case *SkipToken:
			buf.WriteString(fmt.Sprintf("!(%d)", t.n))
		case *SearchToken:
			buf.WriteString(fmt.Sprintf("?(%s)", t.s))
		case *BraToken:
			buf.WriteString("(")
		case *KetToken:
			buf.WriteString(")")
		default:
			buf.WriteString("UnknownToken") // Consider how to handle this case
		}
	}

	return buf.String()
}

func template(s DnaStorage) ([]TemplateToken, error) {
	tokens := make([]TemplateToken, 0)

	for {
		if s.IsEmpty() {
			return nil, Finish
		}
		c := s.GetChar()
		switch c {
		case 'C':
			tokens = append(tokens, NewConstToken('I'))
		case 'F':
			tokens = append(tokens, NewConstToken('C'))
		case 'P':
			tokens = append(tokens, NewConstToken('F'))
		case 'I':
			if s.IsEmpty() {
				return nil, Finish
			}
			cc := s.GetChar()
			switch cc {
			case 'C':
				tokens = append(tokens, NewConstToken('P'))
			case 'F', 'P':
				l, err := nat(s)
				if err != nil {
					return nil, err
				}
				n, err := nat(s)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, NewReferenceToken(n, l))
			case 'I':
				if s.IsEmpty() {
					return nil, Finish
				}
				ccc := s.GetChar()
				switch ccc {
				case 'P':
					n, err := nat(s)
					if err != nil {
						return nil, err
					}
					tokens = append(tokens, NewLenToken(n))
				case 'C', 'F':
					return tokens, nil
				case 'I':
					for p := 3; p < 10; p += 1 {
						fmt.Print(string(s.GetChar()))
					}
				}
			}
		}
	}
}

func templateToString(tokens []TemplateToken) string {
	var buf bytes.Buffer

	for _, token := range tokens {
		switch t := token.(type) {
		case *ConstToken:
			buf.WriteString(string(t.c))
		case *ReferenceToken:
			buf.WriteString(fmt.Sprintf("<%d,%d>", t.n, t.l))
		case *LenToken:
			buf.WriteString(fmt.Sprintf("|%d|", t.n))
		default:
			buf.WriteString("UnknownToken") // Consider how to handle this case
		}
	}

	return buf.String()
}

func envToString(elements []string) string {
	var buf bytes.Buffer

	buf.WriteByte('[')
	if len(elements) != 0 {
		buf.WriteString(elements[0])

		for _, element := range elements[1:] {
			buf.WriteByte(',')
			buf.WriteString(element)
		}
	}
	buf.WriteByte(']')

	return buf.String()
}

type Environment []string

func (e Environment) NotEqual(other Environment) bool {
	if len(e) != len(other) {
		return true
	}
	for i, v := range e {
		if v != other[i] {
			return true
		}
	}
	return false
}

func match(dna DnaStorage, pattern []PatternToken, debug bool) (Environment, error) {
	env := make(Environment, 0)
	stringStack := make([]bytes.Buffer, 0)
	dna.SaveOffset()
	simpleDna, _ := dna.(*SimpleDnaStorage)
	if debug {
		fmt.Fprintln(os.Stderr, "match: saved offset = ", simpleDna.savedOffset, " offset = ", simpleDna.offset)
	}

	for _, token := range pattern {
		switch t := token.(type) {
		case *ConstToken:
			if dna.IsEmpty() {
				return nil, errors.New("Not enough DNA")
			}
			if dna.GetChar() != t.c {
				dna.RestoreOffset()
				if debug {
					fmt.Fprintln(os.Stderr, "pattern mismatched")
				}
				return nil, errors.New("Mismatch")
			}
			for q := range stringStack {
				stringStack[q].WriteByte(t.c)
			}
		case *SkipToken:
			if len(stringStack) == 0 {
				dna.Skip(t.n)
				if dna.IsEmpty() {
					dna.RestoreOffset()
					return nil, errors.New("Not enough DNA")
				}
			} else {
				simpleDna, _ := dna.(*SimpleDnaStorage)
				if debug {
					fmt.Fprintln(os.Stderr, "advance from ", simpleDna.offset, " to ", simpleDna.offset+t.n)
				}
				s, err := dna.GetString(t.n)
				if err != nil {
					dna.RestoreOffset()
					return nil, err
				}
				for q := range stringStack {
					stringStack[q].WriteString(s)
				}
			}
		case *SearchToken:
			index := dna.Index(t.s)
			if index == -1 {
				dna.RestoreOffset()
				fmt.Fprintln(os.Stderr, "pattern not found")
				return nil, errors.New("Pattern not found")
			}
			if debug {
				fmt.Fprintln(os.Stderr, "successfull advance to pattern on ", index, " positions")
			}
			simpleDna, _ := dna.(*SimpleDnaStorage)
			if debug {
				fmt.Fprintln(os.Stderr, "advance from ", simpleDna.offset, " to ", simpleDna.offset+index)
			}

			if len(stringStack) != 0 {
				for p := 0; p < index+len(t.s); p++ {
					if dna.IsEmpty() {
						return nil, errors.New("Not enough DNA")
					}
					c := dna.GetChar()
					for q := range stringStack {
						stringStack[q].WriteByte(c)
					}
				}
			} else {
				dna.Skip(index + len(t.s))
			}
		case *BraToken:
			stringStack = append(stringStack, bytes.Buffer{})
			if debug {
				fmt.Fprintln(os.Stderr, "bra: new buffer")
			}
		case *KetToken:
			env = append(env, stringStack[len(stringStack)-1].String())
			stringStack = stringStack[:len(stringStack)-1]
			if debug {
				fmt.Fprintln(os.Stderr, "ket: add string with length ", len(env[len(env)-1]))
			}
		}
	}

	if debug {
		fmt.Fprintln(os.Stderr, "successfull match")
	}

	return env, nil
}

func formPrefix(template []TemplateToken, env Environment) (string, error) {
	var buf bytes.Buffer

	for _, token := range template {
		switch t := token.(type) {
		case *ConstToken:
			buf.WriteByte(t.c)
		case *ReferenceToken:
			subStr := ""
			if t.n < len(env) {
				subStr = protect(t.l, env[t.n])
			}
			buf.WriteString(subStr)
		case *LenToken:
			clippedN := 0
			if t.n < len(env) {
				clippedN = len(env[t.n])
			}
			buf.WriteString(AsNat(clippedN))
		}
	}

	return buf.String(), nil
}

func AsNat(n int) string {
	var buf bytes.Buffer
	for n != 0 {
		if (n % 2) == 1 {
			buf.WriteByte('C')
		} else {
			buf.WriteByte('I')
		}
		n = n / 2
	}
	buf.WriteByte('P')
	return buf.String()
}

func quote(s string) string {
	var buf bytes.Buffer

	for _, c := range s {
		switch c {
		case 'I':
			buf.WriteByte('C')
		case 'C':
			buf.WriteByte('F')
		case 'F':
			buf.WriteByte('P')
		case 'P':
			buf.WriteString("IC")
		}
	}

	return buf.String()
}

func protect(l int, s string) string {
	for ; l > 0; l-- {
		s = quote(s)
	}
	return s
}

func Step(dna DnaStorage, meter int, debug bool) error {
	simpleDna, _ := dna.(*SimpleDnaStorage)
	if debug {
		fmt.Fprintln(os.Stderr, "initial offset ", simpleDna.offset)
	}
	currentPattern, err := pattern(dna)
	if err != nil {
		return err
	}

	if debug {
		fmt.Fprintln(os.Stderr, "iteration ", meter, " offset after pattern ", simpleDna.offset)
	}

	if debug {
		patternString := patternToString(currentPattern)
		if len(patternString) < 1000 {
			fmt.Fprintln(os.Stderr, "pattern: ", patternString)
		} else {
			fmt.Fprintln(os.Stderr, "pattern too long, len=", len(patternString))
		}
	}

	currentTemplate, err := template(dna)
	if err != nil {
		return err
	}

	if debug {
		if debug {
			fmt.Fprintln(os.Stderr, "offset after template ", simpleDna.offset)
		}

		templateString := templateToString(currentTemplate)

		if len(templateString) < 1000 {
			fmt.Fprintln(os.Stderr, "template: ", templateString)
		} else {
			fmt.Fprintln(os.Stderr, "template string is too long, len =", len(templateString))
		}
	}

	currentEnv, err := match(dna, currentPattern, debug)
	if err != nil {
		return err
	}

	if debug {
		envString := envToString(currentEnv)
		if len(envString) < 1000 {
			fmt.Fprintln(os.Stderr, "env: ", currentEnv)
		} else {
			fmt.Fprintln(os.Stderr, "env is too long, len = ", len(envString))
			fmt.Fprintln(os.Stderr, "env len = ", len(currentEnv))
			for p := range currentEnv {
				fmt.Fprintln(os.Stderr, "\tlen(env[", p, "])=", len(currentEnv[p]))
			}
		}
	}

	currentPrefix, err := formPrefix(currentTemplate, currentEnv)
	if err != nil {
		return err
	}

	if debug {
		if len(currentPrefix) < 1000 {
			fmt.Fprintln(os.Stderr, "prefix:", currentPrefix)
		} else {
			fmt.Fprintln(os.Stderr, "prefix is too long, len = ", len(currentPrefix))
		}
	}

	dna.PrependPrefix(currentPrefix)
	return nil
}

// IIP IP ICP IIC IC IIF  IC C IF P P IIC CFPC
//  (   !  2   )   P  pe  P  I  <0,0>  te

// IIP IP IICP IIC IIC C IIC FCFC
//  (   !  4    )   pe I  te
