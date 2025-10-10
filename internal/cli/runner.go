package cli

import (
	"arithmetic-calc/internal/processor"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func Execute() error {
	cmd := newRootCommand()
	cmd.AddCommand(newShowFileTypesCommand())
	cmd.AddCommand(newShowPipelineTypesCommand())
	return cmd.Execute()
}

func newRootCommand() *cobra.Command {
	config := AppConfig{}

	cmd := &cobra.Command{
		Use:   "npcalc",
		Short: "Обработчик математических выражений в файлах (без паттернов)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApplication(&config)
		},
	}

	cmd.Flags().StringVarP(&config.InputFile, "input", "i", "", "Путь к входному файлу для обработки")
	cmd.Flags().StringVarP(&config.OutputFile, "output", "o", "", "Путь к выходному файлу для обработки")
	cmd.Flags().StringVar(&config.InputType, "i-type", "", "Тип входного файла (пример: txt, json, xml ... можно узнать командой f-types)")
	cmd.Flags().StringVar(&config.OutputType, "o-type", "", "Тип выходного файла (пример: txt, json, xml ... можно узнать командой f-types)")
	cmd.Flags().StringVar(&config.InputPipeline, "i-pipeline", "", "Входной конвейер обработки файла (пример: aes,zip... можно узнать командой p-types)")
	cmd.Flags().StringVar(&config.OutputPipeline, "o-pipeline", "", "Выходной конвейробработки файла (пример: zip,aes... можно узнать командой p-types)")

	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("output")

	return cmd
}

type AppConfig struct {
	InputFile      string
	OutputFile     string
	InputType      string
	OutputType     string
	InputPipeline  string
	OutputPipeline string
}

func runApplication(config *AppConfig) error {
	// получаем сырые данные
	inputData, err := os.ReadFile(config.InputFile)
	if err != nil {
		return err
	}

	// переменная для обработанных данных
	var outputData []byte

	if config.InputType == "" {
		config.InputType = strings.TrimPrefix(filepath.Ext(config.InputFile), ".")
	}

	switch config.InputType {
	case "json":
		var processedJson []byte
		processedJson, err = processor.ProcessJSON(inputData)
		outputData = processedJson
	case "yaml", "yml":
		var processedYAML []byte
		processedYAML, err = processor.ProcessYAML(inputData)
		outputData = processedYAML
	default:
		fmt.Printf("Не нашли нужный тип для обработки: %s. Используем текстовую обработку.", config.InputType)
		var processedString string
		// string() и []byte() - это не функции, это ПРЕОБРАЗОВАНИЕ ТИПОВ (type conversions)
		// инструкция для компилятора он подставляет код для
		// узнать длину, выделить память, скопировать данные,создать заголовок строки
		// указывающий на этот кусок памяти.
		// Выполняется ДИНАМИЧЕСКИ во время выполнения программы используются иструкции выше
		// В момент КОМПИЛЯЦИИ проверка что возможно вообще преобразовать типы
		// Например int(myBytes) нельзя будет ошибка компиляции.

		// Надежны НЕ ВЫЗЫВАЮТ ошибки, НО НУЖНО проверять данные так как могут быть некорректны
		// Например он преобразует строку некорректную, получим ошибку при обработке, так как данные
		// были не корректные.

		// Работа:

		// Преобразуем входные данные в строку из байтов
		processedString, err = processor.ProcessTxt(string(inputData))
		// преобразуем обработанную строку в байты, чтобы записать в файл
		outputData = []byte(processedString)
	}

	if err != nil {
		return fmt.Errorf("Не смогли обработать файл %s: %w", config.InputFile, err)
	}

	// возвращаем сразу результат с ошибкой, тоесть ошибку саму
	// 0 перед числом говорит что оно в 8-чной системе исчисления
	// 0644 - получается 644, 4- чтение, 2- запись, 1 - выполнение
	// первое владельца 6(4+2, чтение и запись)
	// второе для группы 4(только чтение)
	// третье для остальных 4(только чтение)
	// на windows ACL - Access Control List игнорирует perm
	return os.WriteFile(config.OutputFile, outputData, 0644)
}

func newShowFileTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "f-types",
		Short:   "Показывает типы файлов для флагов --input --output",
		Aliases: []string{"ft", "file-types"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Доступные типы файлов:\n")
			for _, t := range SupportedFileTypes {
				fmt.Printf(" - %s\n", t)
			}
			return nil
		},
	}

	return cmd
}

func newShowPipelineTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "p-types",
		Short:   "Показывает типы конвейра для команд --i-pipeline --o-pipeline",
		Aliases: []string{"pt", "pipeline-types"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Доступные операции в конвейере:\n")
			for _, t := range SupportedPipelineTypes {
				fmt.Printf(" - %s\n", t)
			}
			return nil
		},
	}
	return cmd
}
