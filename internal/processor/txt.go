package processor

import "arithmetic-calc/internal/expression"

func ProcessTxt(text string) (string, error) {
	return expression.ProcessString(text)
}
