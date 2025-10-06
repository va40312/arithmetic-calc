package expression

import (
	"math"
	"testing"
)

func TestEvaluator(t *testing.T) {
	testCases := []struct {
		name        string
		rpn         []Token
		expectError bool
		expectedVal float64
	}{
		{
			name: "Проверяем когда недостаточно операторов",
			rpn: []Token{
				{TokenType: Number, Value: "2"},
				{TokenType: Operator, Value: "+"},
			},
			expectError: true,
			expectedVal: 0,
		},
		{
			name: "Проверяем ошибку деления на ноль",
			rpn: []Token{
				{TokenType: Number, Value: "2"},
				{TokenType: Number, Value: "0"},
				{TokenType: Operator, Value: "/"},
			},
			expectError: true,
			expectedVal: 0,
		},
		{
			name: "Проверяем невалидный rpn",
			rpn: []Token{
				{TokenType: Number, Value: "2"},
				{TokenType: Operator, Value: "+"},
				{TokenType: Operator, Value: "+"},
			},
			expectError: true,
			expectedVal: 0,
		},
		{
			name: "Тест простого сложения",
			rpn: []Token{
				{TokenType: Number, Value: "2"},
				{TokenType: Number, Value: "3"},
				{TokenType: Operator, Value: "+"},
			},
			expectedVal: 5.0,
		},
		{
			name: "Тест унарного минуса",
			rpn: []Token{
				{TokenType: Number, Value: "-2"},
				{TokenType: UnaryOperator, Value: "#"},
			},
			expectedVal: 2.0,
		},
		{
			name: "Тест простого умножения",
			rpn: []Token{
				{TokenType: Number, Value: "2"},
				{TokenType: Number, Value: "3"},
				{TokenType: Operator, Value: "*"},
			},
			expectedVal: 6.0,
		},
		{
			name: "Тест простого минуса",
			rpn: []Token{
				{TokenType: Number, Value: "1"},
				{TokenType: Number, Value: "10"},
				{TokenType: Operator, Value: "-"},
			},
			expectedVal: -9.0,
		},
		{
			name: "Тест простого деления",
			rpn: []Token{
				{TokenType: Number, Value: "1"},
				{TokenType: Number, Value: "2"},
				{TokenType: Operator, Value: "/"},
			},
			expectedVal: 0.5,
		},
		{
			name: "Тест умножения на унарный минус от числа с минусом 5 * -(-2)",
			rpn: []Token{
				{TokenType: Number, Value: "5"},
				{TokenType: Number, Value: "-2"},
				{TokenType: UnaryOperator, Value: "#"},
				{TokenType: Operator, Value: "*"},
			},
			expectedVal: 10.0,
		},
		{
			name: "Тест очень сложного выражения со всеми операторами - ( - ( - ( 3 + 2) + 5 )) * - ( 1 + 9) + 6",
			rpn: []Token{
				{TokenType: Number, Value: "3"},
				{TokenType: Number, Value: "2"},
				{TokenType: Operator, Value: "+"},
				{TokenType: UnaryOperator, Value: "#"},
				{TokenType: Number, Value: "5"},
				{TokenType: Operator, Value: "+"},
				{TokenType: UnaryOperator, Value: "#"},
				{TokenType: UnaryOperator, Value: "#"},
				{TokenType: Number, Value: "1"},
				{TokenType: Number, Value: "9"},
				{TokenType: Operator, Value: "+"},
				{TokenType: UnaryOperator, Value: "#"},
				{TokenType: Operator, Value: "*"},
				{TokenType: Number, Value: "6"},
				{TokenType: Operator, Value: "+"},
			},
			expectedVal: 6.0,
		},
		{
			name: "Тест работы с 1 числом",
			rpn: []Token{
				{TokenType: Number, Value: "42"},
			},
			expectedVal: 42.0,
		},
		{
			name: "Проверка унарного оператора",
			rpn: []Token{
				{TokenType: Number, Value: "42"},
				{TokenType: UnaryOperator, Value: "#"},
			},
			expectedVal: -42.0,
		},
		{
			name: "Проверка лишних операндов в конце",
			rpn: []Token{
				{TokenType: Number, Value: "2"},
				{TokenType: Number, Value: "3"},
				{TokenType: Number, Value: "4"},
				{TokenType: UnaryOperator, Value: "+"},
			},
			expectError: true,
			expectedVal: 0.0,
		},
		{
			name: "Проверка невалидно rpn с оператором в начале",
			rpn: []Token{
				{TokenType: Operator, Value: "/"},
				{TokenType: Number, Value: "3"},
				{TokenType: Number, Value: "4"},
			},
			expectError: true,
			expectedVal: 0.0,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Evaluate(test.rpn)

			if err != nil && !test.expectError {
				t.Errorf("Unexpected error: %s", err)
			} else if err == nil && test.expectError {
				t.Errorf("Expect error, but not get it")
			}

			// нельзя просто сравнить float в програмировании
			// так как компьютер записывает 0.1 и другие числа как
			// бесконечную переодическую дробь и вынужден в какой-то момент округлить
			// казалось бы, простые операции, накапливаются микроскопические ошибки округления.
			// берем по модулю потому что нам пофиг на знак главное разница
			// далее мы сравниваем разницу с порогом эпсилон 1e-9 = 0.000000001 (допустимая погрешность)
			// если разница меньше чем указанная величина, то считаем что они равны

			if math.Abs(actual-test.expectedVal) > 1e-9 {
				t.Errorf("expected = %f, actual = %f", test.expectedVal, actual)
			}
		})
	}
}
