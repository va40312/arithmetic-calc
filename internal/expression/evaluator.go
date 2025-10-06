package expression

import (
	"fmt"
	"strconv"
)

// Evaluate - принимает rpn последовательности и правильно ее применяет,
// чтобы получить результат в виде числа. Используя для этого стек.
// если rpn не валидный: число не правильные, не существующий оператор в значении или
// последовательность 2 / 0 все верно по синтаксису, но не правильно математически

func Evaluate(rpnTokens []Token) (float64, error) {
	// стек для сложения чисел
	valueStack := make([]float64, 0)

	for _, token := range rpnTokens {
		switch token.TokenType {
		case Number:
			// конвертируем число из строки в численую переменную
			// 64 - указываем размер float, в нашем случае float64
			value, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, fmt.Errorf("[syntax error]: Не удалось преобразовать %s в число: %w", token.Value, err)
			}
			valueStack = append(valueStack, value)
		case Operator:
			// проверяем бинарные операторы
			if len(valueStack) < 2 {
				return 0, fmt.Errorf("[syntax error]: Недостаточно опредандов для оператора %s", token.Value)
			}
			// правый операнд
			b := valueStack[len(valueStack)-1]
			// левый операнд
			a := valueStack[len(valueStack)-2]
			// взяли операнды и удаляем их из вершины стека
			valueStack = valueStack[:len(valueStack)-2]
			var result float64
			switch token.Value {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("[syntax error]: Деление на ноль запрещено")
				}
				result = a / b
			}
			valueStack = append(valueStack, result)
		case UnaryOperator:
			if len(valueStack) < 1 {
				return 0, fmt.Errorf("[syntax error]: Недостаточно операндов для унарного минуса")
			}
			a := valueStack[len(valueStack)-1]
			valueStack = valueStack[:len(valueStack)-1]

			if token.Value != "#" {
				return 0, fmt.Errorf("[syntax error]: Невалидный символ для унарго минуса %s, должен быть символ `#`", token.Value)
			}
			result := -a
			valueStack = append(valueStack, result)
		}
	}

	if len(valueStack) != 1 {
		return 0, fmt.Errorf("[syntax error]: синтаксическая ошибка в выражении, возьможно лишние операторы")
	}

	return valueStack[0], nil
}
