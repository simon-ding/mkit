/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"mkit/dji"

	"github.com/spf13/cobra"
)

// smsCmd represents the sms command
var smsCmd = &cobra.Command{
	Use:   "sms",
	Short: "Get all sms",
	Long: `Get all sms`,
	Run: func(cmd *cobra.Command, args []string) {
		dji, err := dji.NewDjiModem()
		if err != nil {
			log.Printf("Failed to connect to DJI modem: %v", err)
			return
		}
		defer dji.Close()
		sms, err := dji.GetSms()
		if err != nil {
			log.Printf("fail to get sms: %v", err)
			return
		}
		fmt.Println(sms)
	},
}

func init() {
	rootCmd.AddCommand(smsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// smsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// smsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
