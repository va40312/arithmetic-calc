package expression

import (
	"reflect"
	"strings"
	"testing"
)

func TestFinder_SingleSimpleExpression(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Простой случай с 1 выражением",
			text:     "Цена равна 2 + 2 доллара.",
			expected: []string{"2 + 2"},
		},
		{
			name:     "Строка без выражения",
			text:     "Цена равна доллара.",
			expected: nil,
		},
		{
			name: "Строка из 2-х выражений",
			text: "Цена без ндс (7 * 10) - 15 + 12 / (20 - 15) доллара, а с ндс ((57 - 19) * 2 / 15) + 12",
			expected: []string{
				"(7 * 10) - 15 + 12 / (20 - 15)",
				"((57 - 19) * 2 / 15) + 12",
			},
		},
		{
			name: "Строка из 2-х выражений с числами с запятой",
			text: "Цена без ндс (7.1231 * 10) - 15.12 + 12 / (20 - 15.123) доллара, а с ндс ((57.1 - 19) * 2.2 / 15) + 12.432",
			expected: []string{
				"(7.1231 * 10) - 15.12 + 12 / (20 - 15.123)",
				"((57.1 - 19) * 2.2 / 15) + 12.432",
			},
		},
		{
			name: "Сложная строка из англиских русских букв и цифр в перемешку с выражениями",
			text: "22 hello world 123 bla bla bla 12213.23 / 234 _ как дела + *     (  2 + 2) + 2.1 * 2.2 - 2.31  новое выражение (2 * (3+4)) * 2",
			expected: []string{
				"12213.23 / 234",
				"(  2 + 2) + 2.1 * 2.2 - 2.31",
				"(2 * (3+4)) * 2",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			var expected []FoundExpression
			for _, expression := range test.expected {
				expected = append(expected, FoundExpression{
					expression,
					strings.Index(test.text, expression),
					strings.Index(test.text, expression) + len(expression),
				})
			}

			actual := Finder(test.text)
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Ожидали %+v, но получили %+v", expected, actual)
			}
		})
	}
}
