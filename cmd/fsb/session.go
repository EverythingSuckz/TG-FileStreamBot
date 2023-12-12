package main

import (
	"fmt"

	"EverythingSuckz/fsb/pkg/qrlogin"

	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:                "session",
	Short:              "Generate a string session.",
	DisableSuggestions: false,
	Run:                generateSession,
}

func init() {
	sessionCmd.Flags().StringP("login-type", "T", "qr", "The login type to use. Can be either 'qr' or 'phone'")
	sessionCmd.Flags().Int32P("api-id", "I", 0, "The API ID to use for the session (required).")
	sessionCmd.Flags().StringP("api-hash", "H", "", "The API hash to use for the session (required).")
	sessionCmd.MarkFlagRequired("api-id")
	sessionCmd.MarkFlagRequired("api-hash")
}

func generateSession(cmd *cobra.Command, args []string) {
	loginType, _ := cmd.Flags().GetString("login-type")
	apiId, _ := cmd.Flags().GetInt32("api-id")
	apiHash, _ := cmd.Flags().GetString("api-hash")
	if loginType == "qr" {
		qrlogin.GenerateQRSession(int(apiId), apiHash)
	} else if loginType == "phone" {
		generatePhoneSession()
	} else {
		fmt.Println("Invalid login type. Please use either 'qr' or 'phone'")
	}
}

func generatePhoneSession() {
	fmt.Println("Phone session is not implemented yet.")
}
