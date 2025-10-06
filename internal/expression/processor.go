package expression

import (
	"fmt"
	"strings"
)

// Вычисляет и заменяет все мат выражения в строке

func ProcessString(text string) (string, error) {
	foundExpressions := Finder(text)

	if len(foundExpressions) == 0 {
		//fmt.Println("не нашли не одно выражение.")
		return text, nil
	}

	// создаем string builder для эффективной конкатенации
	var builder strings.Builder
	lastIndex := 0

	for _, found := range foundExpressions {
		result, err := Calculate(found.Expression)
		if err != nil {
			return "", fmt.Errorf("обнаружено невалидное выражение '%s': %w", found.Expression, err)
		}
		// замена выражения в строке
		// можем делать так как как итерация по найдем выражениям, а не строке
		// Так делать плохо так как,
		// нарезка resultString[:found.StartPos] - быстрая операция
		// она хранит ссылку на оригинальную строку
		// но вот для опреций  + ... + ....
		// для каждого плюса будет создана отдельная строка
		// и будет скопирован первый операнд, а затем второй
		// далее уже для второго плюса будет опять выделена новая строка
		// и скопирован результат сложения прошлых строк так как это отдельная строка
		// и скопирован второй операнд.
		// В итоге получиться 2 выделенные строки и 2 раза копирование
		// resultString = resultString[:found.StartPos] + resultStr + resultString[found.EndPos:]

		// для решения проблемы используем strings.Builder
		// он создает строку 1 раз в виде байтов
		// и просто дописывает в конец пока не запросим строку
		// если строка больше чем он выделил память он выделяет память в 2 раза
		// больше и копирует старую строку.

		builder.WriteString(text[lastIndex:found.StartPos])
		builder.WriteString(fmt.Sprintf("%g", result))

		// записываем последний индекс выражения
		// потом запишем число в виде строки, а остальное просто пропустим
		lastIndex = found.EndPos
	}
	// добавляем остаток строки без выражения
	builder.WriteString(text[lastIndex:])
	return builder.String(), nil
}
