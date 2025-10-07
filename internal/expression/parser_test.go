package expression

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		name        string
		expr        string
		expectError bool
		expectedRPN []Token
	}{
		{
			name:        "Тестирование простого примера для приоритета",
			expr:        "3 + 2 * 2",
			expectError: false,
			expectedRPN: []Token{
				{
					TokenType: Number,
					Value:     "3",
				},
				{
					TokenType: Number,
					Value:     "2",
				},
				{
					TokenType: Number,
					Value:     "2",
				},
				{
					TokenType: Operator,
					Value:     "*",
				},
				{
					TokenType: Operator,
					Value:     "+",
				},
			},
		},
		{
			name:        "Тестирование простого примера наоборот для приоритета",
			expr:        "2 * 3 - 1",
			expectError: false,
			expectedRPN: []Token{
				{
					TokenType: Number,
					Value:     "2",
				},
				{
					TokenType: Number,
					Value:     "3",
				},
				{
					TokenType: Operator,
					Value:     "*",
				},
				{
					TokenType: Number,
					Value:     "1",
				},
				{
					TokenType: Operator,
					Value:     "-",
				},
			},
		},
		{
			name:        "Тестирование выражение из нескольких операторов с разным приоритетами",
			expr:        "10 * 20 / 5 - 3",
			expectError: false,
			expectedRPN: []Token{
				{
					TokenType: Number,
					Value:     "10",
				},
				{
					TokenType: Number,
					Value:     "20",
				},
				{
					TokenType: Operator,
					Value:     "*",
				},
				{
					TokenType: Number,
					Value:     "5",
				},
				{
					TokenType: Operator,
					Value:     "/",
				},
				{
					TokenType: Number,
					Value:     "3",
				},
				{
					TokenType: Operator,
					Value:     "-",
				},
			},
		},
		{
			name:        "Тестирование синтаксической ошибки c повторами и более одного оператора",
			expr:        "3 +- 3 ++ 4",
			expectError: true,
			expectedRPN: nil,
		},
		{
			name:        "Тестирование синтаксической ошибки с порядком символа",
			expr:        "- 3 3",
			expectError: true,
			expectedRPN: nil,
		},
		{
			name:        "Тестирование синтаксической ошибки с порядком символа",
			expr:        "* 3 3",
			expectError: true,
			expectedRPN: nil,
		},
		{
			name:        "Тестирование синтаксической ошибки с пропуском символа",
			expr:        "(  2  (3+4)) * 2",
			expectError: true,
			expectedRPN: nil,
		},
		{
			name:        "Тестирование синтаксической ошибки с незакрытыми скобками",
			expr:        "2 + (3 + 2 + 3",
			expectError: true,
			expectedRPN: nil,
		},
		{
			name:        "Тестирование синтаксической ошибки с незакрытыми скобками",
			expr:        "-2 + 3 - 2 - 3)",
			expectError: true,
			expectedRPN: nil,
		},
		{
			name:        "Тестирование синтаксической ошибки с незакрытыми скобками",
			expr:        "(-2 -2 ))",
			expectError: true,
			expectedRPN: nil,
		},
		{
			name:        "Тестирование синтаксической ошибки с незакрытыми скобками",
			expr:        "((-2 -2 )",
			expectError: true,
			expectedRPN: nil,
		},
		{
			name:        "Тестирование минуса с унарным минусом",
			expr:        "2 - -5",
			expectError: false,
			expectedRPN: []Token{
				{TokenType: Number, Value: "2"},
				{TokenType: Number, Value: "-5"},
				{TokenType: Operator, Value: "-"},
			},
		},
		{
			name:        "Тестирование синтаксической ошибки с 1 оператором",
			expr:        "-",
			expectError: true,
			expectedRPN: nil,
		},
		// число первое в скобочках
		// 1) -(
		// 2) (-( - перед минусом скобочка
		// 3) 5 * -(
		// 3 2 + # 5 + # # 1 9 + # * 6 +
		{
			name:        "Тестирование сложного вложенного минуса со скобками",
			expr:        "- ( - ( - ( 3 + 2) + 5 )) * - ( 1 + 9) + 6",
			expectError: false,
			expectedRPN: []Token{
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
		},
		{
			name:        "Тестирование унарного минуса со скобками",
			expr:        "-(-2 -2 )",
			expectError: false,
			expectedRPN: []Token{
				{TokenType: Number, Value: "-2"},
				{TokenType: Number, Value: "2"},
				{TokenType: Operator, Value: "-"},      // Это внутрений минус
				{TokenType: UnaryOperator, Value: "#"}, // Это внешний минус
			},
		},
		{
			name:        "Унарный минус в скобках",
			expr:        "5 * (-2)",
			expectError: false,
			expectedRPN: []Token{
				{TokenType: Number, Value: "5"},
				{TokenType: Number, Value: "-2"},
				{TokenType: Operator, Value: "*"},
			},
		},
		{
			name:        "Унарный минус в скобках",
			expr:        "5 * -2",
			expectError: false,
			expectedRPN: []Token{
				{TokenType: Number, Value: "5"},
				{TokenType: Number, Value: "-2"},
				{TokenType: Operator, Value: "*"},
			},
		},
		{
			name:        "Унарный минус в скобках с минусом перед оператором",
			expr:        "5 * -(-2)",
			expectError: false,
			expectedRPN: []Token{
				{TokenType: Number, Value: "5"},
				{TokenType: Number, Value: "-2"},
				{TokenType: UnaryOperator, Value: "#"},
				{TokenType: Operator, Value: "*"},
			},
		},
		{
			name:        "Проверка обычных чисел со скобками",
			expr:        "(5) + (-2)",
			expectError: false,
			expectedRPN: []Token{
				{TokenType: Number, Value: "5"},
				{TokenType: Number, Value: "-2"},
				{TokenType: Operator, Value: "+"},
			},
		},
		{
			name:        "Тестирование унарного минуса со скобками и другим оператором перед ним",
			expr:        "3 * -(-2 -2 )",
			expectError: false,
			expectedRPN: []Token{
				{TokenType: Number, Value: "3"},
				{TokenType: Number, Value: "-2"},
				{TokenType: Number, Value: "2"},
				{TokenType: Operator, Value: "-"},
				{TokenType: UnaryOperator, Value: "#"},
				{TokenType: Operator, Value: "*"},
			},
		},
		{
			name:        "Тестирование синтаксической ошибки с незакрытыми скобками",
			expr:        "-(-3.41 -2.432 ) - 1.003",
			expectError: false,
			expectedRPN: []Token{
				{TokenType: Number, Value: "-3.41"},
				{TokenType: Number, Value: "2.432"},
				{TokenType: Operator, Value: "-"},
				{TokenType: UnaryOperator, Value: "#"},
				{TokenType: Number, Value: "1.003"},
				{TokenType: Operator, Value: "-"},
			},
		},
		{
			name:        "Тестирование синтаксической ошибки с незакрытыми скобками",
			expr:        "(-2 -2 )-",
			expectError: true,
			expectedRPN: nil,
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tokens, err := Lexer(test.expr)
			if err != nil {
				t.Fatalf("Lexer вернул неожиданую ошибку %+v", err)
			}

			actual, err := Parser(tokens)

			if err != nil && !test.expectError {
				t.Errorf("unexpected error(%s). Parser(%+v) = %+v, expect %+v", err, tokens, actual, test.expectedRPN)
			} else if err == nil && test.expectError {
				t.Errorf("expected error, but not get it")
			}

			if !reflect.DeepEqual(actual, test.expectedRPN) {
				t.Errorf("RPN not expected. Parser(%+v) = %+v, expect %+v", tokens, actual, test.expectedRPN)
			}
		})
	}
}
