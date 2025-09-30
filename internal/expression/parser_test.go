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
			name:        "Тестирование не валидного оператора `.`",
			expr:        "2 /. 2",
			expectError: true,
			expectedRPN: nil,
		},
		//{
		//	name:        "Тестирование синтакстической ошибки c повторами и более одного оператора",
		//	expr:        "3 +- 3 ++ 4",
		//	expectError: true,
		//	expectedRPN: nil,
		//},
		//{
		//	name:        "Тестирование синтакстической ошибки с порядком символа",
		//	expr:        "- 3 3",
		//	expectError: true,
		//	expectedRPN: nil,
		//},
		//{
		//	name:        "Тестирование синтакстической ошибки с пропуском символа",
		//	expr:        "(  2  (3+4)) * 2",
		//	expectError: true,
		//	expectedRPN: nil,
		//},
		//{
		//	name:        "Тестирование синтакстической ошибки с незакрытими скобками",
		//	expr:        "2 + (3 + 2 + 3",
		//	expectError: true,
		//	expectedRPN: nil,
		//},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
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
