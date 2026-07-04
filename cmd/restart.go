/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"mkit/dji"

	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart dji moodem",
	Long:  `restart dji moodem`,
	Run: func(cmd *cobra.Command, args []string) {
		dji, err := dji.NewDjiModem()
		if err != nil {
			log.Printf("Failed to connect to DJI modem: %v", err)
			return
		}
		defer dji.Close()
		if err := dji.Restart(); err != nil {
			log.Printf("fail to restart: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
