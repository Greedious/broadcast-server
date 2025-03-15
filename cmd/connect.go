/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "A command to connect to a broadcasting server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateConnectionFlags(cmd)
		if err != nil {
			return err
		}
		return runConnection(cmd)
	},
}

func validateConnectionFlags(cmd *cobra.Command) error {
	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return errors.New("failed to read port flag")
	}

	// Validate port: must be a valid integer between 1 and 65535
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port: %s. Must be a number between 1 and 65535", port)
	}
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return errors.New("cannot read flag name")
	}
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func runConnection(cmd *cobra.Command) error {
	port, _ := cmd.Flags().GetString("port")
	name, _ := cmd.Flags().GetString("name")

	// WebSocket URL with name as a query parameter
	serverURL := fmt.Sprintf("ws://localhost:%s/ws?name=%s", port, url.QueryEscape(name))

	u, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %s", err)
	}

	fmt.Printf("Connecting as '%s' to server ws://localhost:%s...\n", name, port)

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket server: %s", err)
	}
	defer conn.Close()

	fmt.Println("Connected successfully! Type messages and press ENTER to send. Type 'exit' to quit.")

	// Message handling
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Connection closed by server.")
				return
			}
			fmt.Printf("%s\n", message) // Server should send messages with names
		}
	}()

	// Sending messages
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		// Remove the newline character
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Disconnecting...")
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(input))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringP("port", "p", "8080", "Port to connect to the broadcasting server.")
	connectCmd.Flags().StringP("name", "n", "", "Your display name.")
}
