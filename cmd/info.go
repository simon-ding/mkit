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

var signals = [4]string{"●○○○", "●●○○", "●●●○", "●●●●"}
func mapSignal(s int) string { //+CSQ: (0-31,99),(0-7,99)
	if s < 10 {
		return signals[0]
	} else if s < 20 {
		return signals[1]
	} else if s < 30 {
		return signals[2]
	} else if s <= 31 {
		return signals[3]
	}
	return "-"
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "get sim card info",
	Long:  `get sim card info`,
	Run: func(cmd *cobra.Command, args []string) {
		dji, err := dji.NewDjiModem()
		if err != nil {
			log.Printf("Failed to connect to DJI modem: %v", err)
			return
		}
		defer dji.Close()
		info, err := dji.SimInfo()
		if err != nil {
			log.Printf("get info: %v", err)
			return
		}
		fmt.Println("Number  :", info.PhoneNumber)
		fmt.Println("Status  :", info.Status)
		fmt.Println("Operator:", info.Operator)
		fmt.Println("ICCID   :", info.ICCID)
		fmt.Println("Signal  :", mapSignal(info.SignalStrength))
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
