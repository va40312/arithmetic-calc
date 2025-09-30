package expression

import (
	"fmt"
	"regexp"
	"strings"
)

// создаем собственный тип который ведет себя как int.
// Это нужно чтобы код был безопаснее, чтобы нельзя
// было передать другой int случайно, а только TokenType

type TokenType int

// Мы создаем 4 переменные типа TokenType по факту int
// и указываем iota который нумерует переменные от нуля до n переменных
const (
	Number TokenType = iota
	Operator
	LeftParen
	RightParen
)

type Token struct {
	TokenType TokenType
	Value     string
}

// Первый этап вычисления выражения разбитие его на токены
// Этим занимается Lexer так же его могут называть токенизатор, сканер
// Лексер берет сырую строку и разбивает ее на структуры чтобы проще работать с ними в других функциях
// чтобы не проверять по стораз символ мы сразу работает с структой которую проще проверить
// он обрасывает мусор например пробелы или другие и готов к следующей обработке, знаем с чем работаем

var tokenizerRegex = regexp.MustCompile(`\d+(?:\.\d+)?|[\(\)\-\+\*\/]`)

func Lexer(text string) ([]Token, error) {
	rawTokens := tokenizerRegex.FindAllString(text, -1)
	if rawTokens == nil {
		return nil, fmt.Errorf("[lexer error]: выражение не содержит валидных токенов")
	}

	joinedTokens := strings.Join(rawTokens, "")
	exprWithoutSpaces := strings.ReplaceAll(text, " ", "")
	tokensWithoutSpaces := strings.ReplaceAll(joinedTokens, " ", "")

	if exprWithoutSpaces != tokensWithoutSpaces {
		return nil, fmt.Errorf("[lexer error]: обнаружены недопустимые символы в выражении")
	}

	var tokens []Token
	for _, rawToken := range rawTokens {
		var token Token

		switch rawToken {
		case "(":
			token.TokenType = LeftParen
		case ")":
			token.TokenType = RightParen
		case "+", "-", "/", "*":
			token.TokenType = Operator
		default:
			token.TokenType = Number
		}

		token.Value = rawToken
		tokens = append(tokens, token)
	}

	return tokens, nil
}
