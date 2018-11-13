package parser_test

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/internal/parser"
	"github.com/influxdata/flux/internal/token"
)

var CompareOptions = []cmp.Option{
	cmp.Transformer("", func(re *regexp.Regexp) string {
		return re.String()
	}),
}

type Token struct {
	Pos   token.Pos
	Token token.Token
	Lit   string
}

type Scanner struct {
	Tokens   []Token
	i        int
	buffered bool
}

func (s *Scanner) Scan() (token.Pos, token.Token, string) {
	pos, tok, lit := s.ScanWithRegex()
	if tok == token.REGEX {
		// The parser was asking for a non regex token and our static
		// scanner gave it one. This indicates a bug in the parser since
		// the parser should have invoked Scan instead.
		s.Unread()
		return 0, token.ILLEGAL, ""
	}
	return pos, tok, lit
}

func (s *Scanner) ScanWithRegex() (token.Pos, token.Token, string) {
	if s.i >= len(s.Tokens) {
		return 0, token.EOF, ""
	}
	tok := s.Tokens[s.i]
	s.i++
	s.buffered = false
	return tok.Pos, tok.Token, tok.Lit
}

func (s *Scanner) Unread() {
	// Buffered indicates that the value is "buffered". Since we keep everything
	// in memory, we use it to prevent unread from going backwards more than once
	// to prevent accidentally using a lookahead of 2 when testing the parser.
	if !s.buffered {
		s.buffered = true
		s.i--
	}
}

func TestParser(t *testing.T) {
	for _, tt := range []struct {
		name   string
		tokens []Token
		want   *ast.Program
	}{
		{
			name: "declare variable as an int",
			tokens: []Token{
				{Token: token.IDENT, Lit: `howdy`},
				{Token: token.ASSIGN, Lit: `=`},
				{Token: token.INT, Lit: `1`},
			},
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.VariableDeclaration{
						Declarations: []*ast.VariableDeclarator{{
							ID:   &ast.Identifier{Name: "howdy"},
							Init: &ast.IntegerLiteral{Value: 1},
						}},
					},
				},
			},
		},
		{
			name: "declare variable as a float",
			tokens: []Token{
				{Token: token.IDENT, Lit: `howdy`},
				{Token: token.ASSIGN, Lit: `=`},
				{Token: token.FLOAT, Lit: `1.1`},
			},
			want: &ast.Program{
				Body: []ast.Statement{
					&ast.VariableDeclaration{
						Declarations: []*ast.VariableDeclarator{{
							ID:   &ast.Identifier{Name: "howdy"},
							Init: &ast.FloatLiteral{Value: 1.1},
						}},
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			scanner := &Scanner{Tokens: tt.tokens}
			result, err := parser.NewAST(scanner)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got, want := result, tt.want; !cmp.Equal(want, got, CompareOptions...) {
				t.Fatalf("unexpected statement -want/+got\n%s", cmp.Diff(want, got, CompareOptions...))
			}
		})
	}
}
