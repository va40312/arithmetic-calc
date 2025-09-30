package expression

type operatorInfo struct {
	priority int
}

var operators = map[string]operatorInfo{
	"+": {priority: 1},
	"-": {priority: 1},
	"*": {priority: 2},
	"/": {priority: 2},
}

// ----- функции для работы со стеком
// взять первый элемент не удаляя его
func peek(stack []Token) Token {
	if len(stack) == 0 {
		return Token{}
	}
	return stack[len(stack)-1]
}

// взять и удалить 1 элемент
func pop(stack []Token) (Token, []Token) {
	if len(stack) == 0 {
		return Token{}, stack
	}
	token := peek(stack)
	stack = stack[:len(stack)-1]
	return token, stack
}

func Parser(tokens []Token) ([]Token, error) {
	return nil, nil
}
