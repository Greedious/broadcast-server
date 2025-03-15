/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"broadcast-server/cmd/utils"
	"broadcast-server/server"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A command to start a broadcasting server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateServerFlags(cmd); err != nil {
			return err
		}
		return runBroadcastServer(cmd)
	},
}

func validateServerFlags(cmd *cobra.Command) error {
	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return errors.New("failed to read port flag")
	}

	if isPortTaken := utils.IsPortTaken(port); isPortTaken {
		return fmt.Errorf("error: port %s is already in use", port)
	}

	// Validate port: must be a valid integer between 1 and 65535
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port: %s. Must be a number between 1 and 65535", port)
	}
	return nil
}

func runBroadcastServer(cmd *cobra.Command) error {
	port, _ := cmd.Flags().GetString("port")

	address := ":" + port

	fmt.Printf("Starting server on port %s: \n", port)

	server := server.NewServer()
	http.HandleFunc("/ws", server.HandleConnections)

	go server.HandleMessages()

	return http.ListenAndServe(address, nil)
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP("port", "p", "8080", "choose port to run the broadcasting server on.")

}
