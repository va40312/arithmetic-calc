package expression

import "fmt"

// принимает строку с выражением и возвращает результат или ошибку

func Calculate(text string) (float64, error) {
	tokens, err := Lexer(text)
	if err != nil {
		return 0, fmt.Errorf("Ошибка лексера: %w", err)
	}
	//fmt.Printf("tokens: %+v\n", tokens)
	rpn, err := Parser(tokens)
	if err != nil {
		return 0, fmt.Errorf("Ошибка парсера: %w\n", err)
	}
	//fmt.Printf("rpn: %+v\n", rpn)
	result, err := Evaluate(rpn)
	if err != nil {
		return 0, fmt.Errorf("Ошибка вычисления: %w", err)
	}
	//fmt.Printf("result: %+v\n", result)
	return result, nil
}
