package expression

import "testing"

func TestCalculator(t *testing.T) {
	testCases := []struct {
		name        string
		expression  string
		expectedVal float64
		expectError bool
	}{
		{
			name:        "Базовый пример с приоритетом",
			expression:  "3 + 2 * 2",
			expectError: false,
			expectedVal: 7.0,
		},
		{
			name:        "Сложное выражение с унарными минусами и скобками",
			expression:  "5 * -(-2 - 2)",
			expectError: false,
			expectedVal: 20.0, // 5 * -(-4) = 5 * 4 = 20
		},
		{
			name:        "Полный хардкор-тест",
			expression:  "- ( - ( - ( 3 + 2) + 5 )) * - ( 1 + 9) + 6",
			expectError: false,
			expectedVal: 6.0,
		},
		{
			name:        "Деление на ноль",
			expression:  "100 / (5 - 5)",
			expectError: true,
		},
		{
			name:        "Лексическая ошибка (недопустимый символ)",
			expression:  "2 + & 3",
			expectError: true,
		},
		{
			name:        "Синтаксическая ошибка (два оператора)",
			expression:  "2 ++ 3",
			expectError: true,
		},
		{
			name:        "Синтаксическая ошибка (пропущен оператор)",
			expression:  " (2) (3)",
			expectError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Calculate(test.expression)

			if !test.expectError && err != nil {
				t.Errorf("not expected error (%s) in Calculate(%s) = %f, expected %f", err, test.expression, actual, test.expectedVal)
			} else if test.expectError && err == nil {
				t.Errorf("expected error, but not get it in Calculate(%s) = %f, expected %f", test.expression, actual, test.expectedVal)
			}

			if actual != test.expectedVal {
				t.Errorf("Calculate(%s) = %f, expected %f", test.expression, actual, test.expectedVal)
			}
		})
	}
}
