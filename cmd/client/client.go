package client

import (
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Days     int
	Output   string
	FilePath string
	Insecure bool
}

type Client struct {
	Config Config
	Conn   *websocket.Conn
	// MessageID represents the sequential ID of each message after the initial auth.
	// These must be incremented with each subsequent request, otherwise the API will
	// return an error.
	MessageID int
}

// APIResponse represents the structure of the response received from the Home Assistant API.
type APIResponse struct {
	ID      int    `json:"id"`      // ID is the unique identifier of the response.
	Type    string `json:"type"`    // Type is the type of the response.
	Success bool   `json:"success"` // Success indicates whether the response was successful or not.
	Result  map[string][]struct {
		Change float64 `json:"change"`
		End    int64   `json:"end"`
		Start  int64   `json:"start"`
	} `json:"result"` // Result contains the data returned by the API.
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

const hoursInADay = 24

// defaultHandshakeTimeout is the timeout for the websocket handshake.
// Increased from 10s to 15s to avoid timeouts on slower home networks.
const defaultHandshakeTimeout = 15 * time.Second

func New(cfg Config) *Client {
	return &Client{
		Config: cfg,
	}
}

func (c *Client) Connect() error {
	c.MessageID = 1

	// Set up the websocket dialer
	dialer := websocket.Dialer{
		HandshakeTimeout: defaultHandshakeTimeout,
	}

	// Work out the URL to dial
	if viper.GetString("url") == "" {
		return fmt.Errorf("url is required")
	}
	dialURL, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	if dialURL.Scheme == "http" {
		dialURL.Scheme = "ws"
	} else if dialURL.Scheme == "https" {
		dialURL.Scheme = "wss"
	}
	dialURL.Path = "/api/websocket"

	// Skip TLS verification if insecure flag is set
	if c.Config.Insecure {
		dialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// Dial the websocket
	log.Info().Msgf("connecting to %s", dialURL.String())
	conn, _, err := dialer.Dial(dialURL.String(), nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	log.Info().Msg("connected")

	// Read the initial message
	var initMsg map[string]any
	if err := conn.ReadJSON(&initMsg); err != nil {
		return fmt.Errorf("initial message: %w", err)
	}

	// Send the authentication message
	if err := conn.WriteJSON(map[string]string{
		"type":         "auth",
		"access_token": viper.GetString("api_key"),
	}); err != nil {
		return fmt.Errorf("auth message: %w", err)
	}

	// Read the authentication response
	var authResp map[string]any
	if err := conn.ReadJSON(&authResp); err != nil {
		return fmt.Errorf("auth response: %w", err)
	}
	if authResp["type"] != "auth_ok" {
		return fmt.Errorf("authentication failed: %v", authResp["message"])
	}
	log.Info().Msg("authenticated")

	c.Conn = conn
	return nil
}

// computePowerStats computes the power statistics for a given number of days a
