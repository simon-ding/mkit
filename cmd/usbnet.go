/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"mkit/dji"
	"strconv"

	"github.com/spf13/cobra"
)

// usbnetCmd represents the usbnet command
var usbnetCmd = &cobra.Command{
	Use:   "usbnet",
	Short: "get or set usbnet mode",
	Long: `get or set usbnet mode:
	0：RNDIS
	1：ECM
	2：MBIM
	`,
	Run: func(cmd *cobra.Command, args []string) {
		dji, err := dji.NewDjiModem()
		if err != nil {
			log.Printf("Failed to connect to DJI modem: %v", err)
			return
		}
		defer dji.Close()
		if len(args) == 0 {
			mode, err := dji.GetUsbnetMode()
			if err != nil {
				log.Printf("fail to get usbnet: %v", err)
				return
			}
			fmt.Println("Usbnet:", mode)
		} else {
			mode, err := strconv.Atoi(args[0])
			if err != nil {
				log.Printf("%v", err)
				return
			}
			if err := dji.SetUsbnetMode(mode); err != nil {
				log.Printf("set mode: %v", err)
			} else {
				fmt.Println("OK")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(usbnetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usbnetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usbnetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
