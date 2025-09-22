package main

import (
	"arithmetic-calc/internal/cli"
	"fmt"
	"os"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
