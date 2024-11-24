package main

import (
    "fmt"
    "github.com/spf13/cobra"
)

func Export() *cobra.Command{
    return &cobra.Command{
        Use: "export from source streaming",
        Args: cobra.MinimumNArgs(2),
        Run: func(cmd *cobra.Command, args []string){
            fmt.Printf("Solicitacao: exportar: [%s], origem: [%s]\n", args[1], args[0])
        },
    }
}
