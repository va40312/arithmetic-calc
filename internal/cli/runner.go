package cli

import (
	"fmt"
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
	return nil
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
