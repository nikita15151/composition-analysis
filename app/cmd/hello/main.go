package main

import (
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

func main() {
	var showVersion bool

	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	root := &cobra.Command{
		Use:   "hello",
		Short: "Учебное приложение для композиционного анализа ПО",
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("hello version %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
				return
			}

			logger.Info().Msg("Запуск учебного приложения")
			fmt.Println("Hello, World!")
			fmt.Printf("Version: %s\n", version)
			fmt.Printf("Run ID: %s\n", uuid.NewString())
		},
	}

	root.Flags().BoolVarP(&showVersion, "version", "v", false, "показать версию")

	if err := root.Execute(); err != nil {
		logger.Error().Err(err).Msg("ошибка выполнения")
	}
}
