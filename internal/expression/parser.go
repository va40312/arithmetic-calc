package expression

import "fmt"

type operatorInfo struct {
	priority    int
	isLeftAssoc bool
}

var operators = map[string]operatorInfo{
	"+": {priority: 1, isLeftAssoc: true},
	"-": {priority: 1, isLeftAssoc: true},
	"*": {priority: 2, isLeftAssoc: true},
	"/": {priority: 2, isLeftAssoc: true},
	// унарный минус
	"#": {priority: 3, isLeftAssoc: false},
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

// функция принимает срез токенов и возвращает склееные токены
// с унарным минусом
func preprocessUnaryMinus(tokens []Token) ([]Token, error) {
	var outputTokens []Token
	for index := 0; index < len(tokens); index++ {
		token := tokens[index]
		switch token.TokenType {
		case Operator:
			if token.Value != "-" {
				outputTokens = append(outputTokens, token)
				continue
			}
			if index+1 > len(tokens)-1 {
				outputTokens = append(outputTokens, token)
				continue
			}
			// сработает проверка i == 0 ниже которая
			//if index-1 < 0 {
			//	outputTokens = append(outputTokens, token)
			//	continue
			//}

			nextToken := tokens[index+1]

			if nextToken.TokenType != Number {
				if nextToken.TokenType == LeftParen {
					outputTokens = append(outputTokens, Token{
						TokenType: UnaryOperator,
						Value:     "#",
					})
				} else {
					outputTokens = append(outputTokens, token)
				}
				continue
			}

			if index == 0 {
				outputTokens = append(outputTokens, Token{
					TokenType: Number,
					Value:     "-" + nextToken.Value,
				})
				index++ // увеличиваем счетчик на 1 чтобы пропустить число
				continue
			}

			previousToken := tokens[index-1]
			if previousToken.TokenType != Number && nextToken.TokenType == LeftParen {
				outputTokens = append(outputTokens,
					Token{
						TokenType: UnaryOperator,
						Value:     "#",
					},
				)
				continue
			}

			if previousToken.TokenType == LeftParen {
				outputTokens = append(outputTokens, Token{
					TokenType: Number,
					Value:     "-" + nextToken.Value,
				})
				index++ // увеличиваем счетчик на 1 чтобы пропустить число
				continue
			}

			if previousToken.TokenType == Operator {
				outputTokens = append(outputTokens, Token{
					TokenType: Number,
					Value:     "-" + nextToken.Value,
				})
				index++ // увеличиваем счетчик на 1 чтобы пропустить число
				continue
			}
			outputTokens = append(outputTokens, token)
		default:
			outputTokens = append(outputTokens, token)
			//if token.TokenType == RightParen {
			//	outputTokens = append(outputTokens, Token{
			//		TokenType: RightParen,
			//		Value:     ")",
			//	})
			//}
		}
	}
	// обhаботать синтаксическую ошибку если есть унарный минус
	// но нет числа 2 - (- * 2) после минуса нет числа, это ошибка
	// !!!!!!!
	//
	// -1 - 2
	// 1 + (-1 + 2)
	// - -234 стоит еще минус перед выражением
	// обобщим чтобы минус был перед любым оператором
	// * -2, / -2, + -2,

	return outputTokens, nil
}

func validSyntax(tokens []Token) error {
	for index := 0; index < len(tokens); index++ {
		token := tokens[index]
		switch token.TokenType {
		case Number:
			if index+1 > len(tokens)-1 {
				continue
			}
			nextToken := tokens[index+1]
			if nextToken.TokenType == Number || nextToken.TokenType == LeftParen {
				return fmt.Errorf("[syntax error]: отсутствует оператор между числом %s и %s", token.Value, nextToken.Value)
			}

			if index == 0 {
				continue
			}
			previousToken := tokens[index-1]
			if previousToken.TokenType == Number || previousToken.TokenType == RightParen {
				return fmt.Errorf("[syntax error]: отсутствует оператор между числом `%s` и `%s`", token.Value, previousToken.Value)
			}
		case Operator:
			if index+1 > len(tokens)-1 {
				return fmt.Errorf("[syntax error]: бинарный оператор `%s` не может быть без правого операнда", token.Value)
			}

			nextToken := tokens[index+1]
			isUnaryMinus := token.Value == "-" && nextToken.TokenType == LeftParen

			if isUnaryMinus {
				// вставить фиктивный ноль
				continue
			} else if index == 0 {
				return fmt.Errorf("[syntax error]: бинарный оператор `%s` не может быть без левого операнда", token.Value)
			}

			previousToken := tokens[index-1]

			if previousToken.TokenType != Number && previousToken.TokenType != RightParen {
				return fmt.Errorf("[syntax error]: перед бинарным оператор `%s` должно быть число", token.Value)
			} else if nextToken.TokenType != Number && nextToken.TokenType != LeftParen && nextToken.TokenType != UnaryOperator {
				return fmt.Errorf("[syntax error]: после бинарного оператора `%s` должно быть число", token.Value)
			}
		}
	}
	return nil
}

func Parser(tokens []Token) ([]Token, error) {
	var err error
	// преобразовать унарный минус
	tokens, err = preprocessUnaryMinus(tokens)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("preprocessUnaryMinus = %+v\n", tokens)

	// проверка синтаксиса на ошибки перед применением алгоритма
	if err = validSyntax(tokens); err != nil {
		return nil, err
	}

	// начало алгоритма сортировочной станции
	var output []Token
	var operatorStack []Token
	var topOperator Token

	for _, token := range tokens {
		switch token.TokenType {
		case Number:
			output = append(output, token)

		case Operator, UnaryOperator:
			// проверяем что стек не пустой, что-то в нем есть на вершине
			// 2 + 2 * 2
			// 2 * 2 + 2
			// 10 * 20 / 5 + 3
			// A + B * C / D - E
			for len(operatorStack) > 0 && (peek(operatorStack).TokenType == Operator || peek(operatorStack).TokenType == UnaryOperator) {
				topOperator = peek(operatorStack)
				topOperatorInfo := operators[topOperator.Value]
				tokenOperatorInfo := operators[token.Value]

				if (tokenOperatorInfo.isLeftAssoc &&
					tokenOperatorInfo.priority <= topOperatorInfo.priority) ||
					(!tokenOperatorInfo.isLeftAssoc && tokenOperatorInfo.priority < topOperatorInfo.priority) {
					// убираем верхний элемент в output
					topOperator, operatorStack = pop(operatorStack)
					output = append(output, topOperator)
				} else {
					// прошли все вершины стека и достигли того, что вытолкнули
					// все в правильном порядке и приоритете
					break
				}
			}
			// добавляем в stack оператор новый
			operatorStack = append(operatorStack, token)

		case LeftParen:
			operatorStack = append(operatorStack, token)

		case RightParen:
			// удаляем все пока не встретим открывающую скобку
			// если нет левой скобочки ошибка
			for len(operatorStack) > 0 && peek(operatorStack).TokenType != LeftParen {
				topOperator, operatorStack = pop(operatorStack)
				output = append(output, topOperator)
			}
			// на этом этапе должна остаться LeftParen иначе ошибка синтаксиса
			if len(operatorStack) == 0 {
				return nil, fmt.Errorf("[syntax error]: нет открывающей скобки")
			}

			_, operatorStack = pop(operatorStack)
		}
	}
	// все остатки operatorStack нужно добавить в output после цикла
	// которые не прошли выталкивание из за порядка они последние
	// или приоритета
	for len(operatorStack) > 0 {
		topOperator, operatorStack = pop(operatorStack)
		if topOperator.TokenType == LeftParen {
			return nil, fmt.Errorf("[syntax error]: не найдено закрывающей скобочки для (")
		}
		output = append(output, topOperator)
	}
	// конец алгоритма сортировочной станции
	return output, nil
}
