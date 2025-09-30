package expression

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		expected    []Token
		expectError bool
	}{
		{
			name:        "Проверяем самый простой случай разбивания на токены",
			text:        "2 + 2",
			expectError: false,
			expected: []Token{
				{
					TokenType: Number,
					Value:     "2",
				},
				{
					TokenType: Operator,
					Value:     "+",
				},
				{
					TokenType: Number,
					Value:     "2",
				},
			},
		},
		{
			name:        "Проверяем когда пустая строка",
			text:        "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Проверяем когда случайный текст",
			text:        "adsfasdfasdf",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Проверяем когда текст и выражения",
			text:        "2 + 2 будет 4",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Проверяем со скобками",
			text:        "(3.52 + 4.58)",
			expectError: false,
			expected: []Token{
				{
					TokenType: LeftParen,
					Value:     "(",
				},
				{
					TokenType: Number,
					Value:     "3.52",
				},
				{
					TokenType: Operator,
					Value:     "+",
				},
				{
					TokenType: Number,
					Value:     "4.58",
				},
				{
					TokenType: RightParen,
					Value:     ")",
				},
			},
		},
		{
			name:        "Проверяем со скобками и другими знаками",
			text:        "(3.52 + 4.58) * 2.563",
			expectError: false,
			expected: []Token{
				{
					TokenType: LeftParen,
					Value:     "(",
				},
				{
					TokenType: Number,
					Value:     "3.52",
				},
				{
					TokenType: Operator,
					Value:     "+",
				},
				{
					TokenType: Number,
					Value:     "4.58",
				},
				{
					TokenType: RightParen,
					Value:     ")",
				},
				{
					TokenType: Operator,
					Value:     "*",
				},
				{
					TokenType: Number,
					Value:     "2.563",
				},
			},
		},
		{
			name:        "Проверяем работу с минусом",
			text:        "-3.52 - 4.58 - (-42.12 - 15)",
			expectError: false,
			expected: []Token{
				{
					TokenType: Operator,
					Value:     "-",
				},
				{
					TokenType: Number,
					Value:     "3.52",
				},
				{
					TokenType: Operator,
					Value:     "-",
				},
				{
					TokenType: Number,
					Value:     "4.58",
				},
				{
					TokenType: Operator,
					Value:     "-",
				},
				{
					TokenType: LeftParen,
					Value:     "(",
				},
				{
					TokenType: Operator,
					Value:     "-",
				},
				{
					TokenType: Number,
					Value:     "42.12",
				},
				{
					TokenType: Operator,
					Value:     "-",
				},
				{
					TokenType: Number,
					Value:     "15",
				},
				{
					TokenType: RightParen,
					Value:     ")",
				},
			},
		},
		{
			name:        "Проверяем работу с невалидными данными",
			text:        "2 /. 2",
			expectError: true,
			expected:    nil,
		},
		{
			name:        "Проверяем работу с невалидными данными",
			text:        "2 . 2",
			expectError: true,
			expected:    nil,
		},
		{
			name:        "Проверяем работу с невалидными данными",
			text:        "2 +-./ 2",
			expectError: true,
			expected:    nil,
		},
		{
			name:        "Проверяем работу с двойными операторами или разными",
			text:        "2 +- 2 // 4",
			expectError: true,
			expected: []Token{
				{TokenType: Number, Value: "2"},
				{TokenType: Operator, Value: "+"},
				{TokenType: Operator, Value: "-"},
				{TokenType: Number, Value: "2"},
				{TokenType: Operator, Value: "/"},
				{TokenType: Operator, Value: "/"},
				{TokenType: Number, Value: "4"},
			},
		},
		{
			name:        "Проверяем работу с неправильным порядком операторов",
			text:        "- 2 2",
			expectError: true,
			expected: []Token{
				{TokenType: Operator, Value: "-"},
				{TokenType: Number, Value: "2"},
				{TokenType: Number, Value: "2"},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Lexer(test.text)

			if err != nil && !test.expectError {
				t.Errorf("not expected error (%s)", err)
			} else if err == nil && test.expectError {
				t.Errorf("expected error, but not get it")
			}

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Lexer(%s) = %+v, exprect %+v", test.text, actual, test.expected)
			}
		})
	}
}
