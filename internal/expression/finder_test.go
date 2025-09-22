package expression

import (
	"reflect"
	"testing"
)

func TestFinder_SingleSimpleExpression(t *testing.T) {
	text := "Цена равно 2 + 2 доллара."
	// text2 := "22 hello world 123 bla bla bla 12213.23 / 234 _ как дела + *     (  2 + 2) + 2.1 * 2.2 - 2.31  новое выражение (2 * (3+4)) * 2"
	expected := []FoundExpression{
		{
			"2 + 2",
			11,
			16,
		},
	}

	actual := Finder(text)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Ожидали %+v, но получили %+v", expected, actual)
	}
}
