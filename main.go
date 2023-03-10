package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "ai gateway for develop team",
	}
	rootCmd.AddCommand(initCmd, runCmd)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
