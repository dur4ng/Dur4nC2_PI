package cli

import (
	"Dur4nC2/server/console"
	"Dur4nC2/server/console/command"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Setup Configurations, can load jsons conf files and run console shell
var rootCmd = &cobra.Command{
	Use:   "dur4nc2-server",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic:\n%s", debug.Stack())
				fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
				os.Exit(99)
			}
		}()

		err := StartClientConsole()
		if err != nil {
			fmt.Printf("[!] %s\n", err)
		}
	},
}

func StartClientConsole() error {
	return console.Start(command.BindCommands)
}

// Execute - Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
