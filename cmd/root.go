package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
        Short: "Export content from source streaming",
    }

func Execute(){
    if err := RootCmd.Execute(); err != nil {
	fmt.Println(err)
	os.Exit(1)
    }
}
