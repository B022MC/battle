package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

var (
	// ErrGameServerConnection is returned when connection to game server fails
	ErrGameServerConnection = errors.New("failed to connect to game server")

	// ErrInvalidCredentials is returned when game account credentials are invalid
	ErrInvalidCredentials = errors.New("invalid game account credentials")

	// ErrGameServerTimeout is returned when game server request times out
	ErrGameServerTimeout = errors.New("game server request timeout")
)

// GameAccountValidator validates game accounts with the game server
type GameAccountValidator struct {
	gameServerAddr string
	timeout        time.Duration
}

// NewGameAccountValidator creates a new game account validator
func NewGameAccountValidator(gameServerAddr string) *GameAccountValidator {
	if gameServerAddr == "" {
		gameServerAddr = "androidsc.foxuc.com:8200" // Default game server
	}
	return &GameAccountValidator{
		gameServerAddr: gameServerAddr,
		timeout:        30 * time.Second,
	}
}

// GameAccountInfo represents game account information from game server
type GameAccountInfo struct {
	GameUserID string `json:"game_user_id"` // Game server user ID
	GameID     string `json:"game_id"`      // Game ID
	Nickname   string `json:"nickname"`     // Game nickname
	Account    string `json:"account"`      // Account (phone number)
	Success    bool   `json:"success"`      // Whether validation succeeded
	ErrorMsg   string `json:"error_msg"`    // Error message if failed
}

// ValidateByMobile validates game account using mobile phone number
// This is ported from passing-dragonfly/waiter/plaza/mobilelogon.go
func (v *GameAccountValidator) ValidateByMobile(ctx context.Context, mobile, password string) (*GameAccountInfo, error) {
	return v.validateAccount(ctx, mobile, password, true)
}

// ValidateByAccount validates game account using account name
func (v *GameAccountValidator) ValidateByAccount(ctx context.Context, account, password string) (*GameAccountInfo, error) {
	return v.validateAccount(ctx, account, password, false)
}

// validateAccount performs the actual validation with game server
func (v *GameAccountValidator) validateAccount(ctx context.Context, identifier, password string, isMobile bool) (*GameAccountInfo, error) {
	// Create connection to game server with timeout
	conn, err := net.DialTimeout("tcp", v.gameServerAddr, v.timeout)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGameServerConnection, err)
	}
	defer conn.Close()

	// Set read/write deadlines
	deadline := time.Now().Add(v.timeout)
	if err := conn.SetDeadline(deadline); err != nil {
		return nil, fmt.Errorf("failed to set deadline: %w", err)
	}

	// Enable TCP keepalive
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetKeepAlive(true)
		_ = tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	// Build login command
	cmd := v.buildLoginCommand(identifier, password, isMobile)

	// Encrypt and send command
	encoder := NewProtocolEncoder()
	encrypted := encoder.Encrypt(cmd)

	_, err = conn.Write(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to send login command: %w", err)
	}

	// Read response
	response := make([]byte, 4096)
	n, err := conn.Read(response)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, ErrGameServerTimeout
		}
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	info, err := v.parseLoginResponse(response[:n])
	if err != nil {
		return nil, err
	}

	if !info.Success {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCredentials, info.ErrorMsg)
	}

	return info, nil
}

// buildLoginCommand builds the login command packet
func (v *GameAccountValidator) buildLoginCommand(identifier, password string, isMobile bool) []byte {
	// This is a simplified version - actual implementation should match
	// the protocol used in passing-dragonfly/waiter/plaza/session.go

	loginMode := 2 // Account mode
	if isMobile {
		loginMode = 1 // Mobile mode
	}

	// Hash password with MD5
	pwdMD5 := v.hashPassword(password)

	// Build command structure (simplified)
	cmd := map[string]interface{}{
		"cmd":        "account_logon",
		"login_mode": loginMode,
		"identifier": identifier,
		"password":   pwdMD5,
		"timestamp":  time.Now().Unix(),
	}

	data, _ := json.Marshal(cmd)
	return data
}

// parseLoginResponse parses the login response from game server
func (v *GameAccountValidator) parseLoginResponse(data []byte) (*GameAccountInfo, error) {
	// This is a simplified version - actual implementation should match
	// the protocol used in passing-dragonfly/waiter/plaza/session.go

	var response struct {
		Success    bool   `json:"success"`
		ErrorMsg   string `json:"error_msg"`
		GameUserID string `json:"game_user_id"`
		GameID     string `json:"game_id"`
		Nickname   string `json:"nickname"`
		Account    string `json:"account"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &GameAccountInfo{
		GameUserID: response.GameUserID,
		GameID:     response.GameID,
		Nickname:   response.Nickname,
		Account:    response.Account,
		Success:    response.Success,
		ErrorMsg:   response.ErrorMsg,
	}, nil
}

// hashPassword hashes password with MD5
func (v *GameAccountValidator) hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// ProtocolEncoder handles protocol encoding/encryption
type ProtocolEncoder struct {
	// Add encryption keys and methods as needed
}

// NewProtocolEncoder creates a new protocol encoder
func NewProtocolEncoder() *ProtocolEncoder {
	return &ProtocolEncoder{}
}

// Encrypt encrypts the command data
// This should match the encryption used in passing-dragonfly/waiter/plaza/session.go
func (e *ProtocolEncoder) Encrypt(data []byte) []byte {
	// Simplified version - actual implementation should include proper encryption
	// For now, just return the data as-is
	// TODO: Implement actual encryption algorithm from legacy code
	return data
}

// Decrypt decrypts the response data
func (e *ProtocolEncoder) Decrypt(data []byte) []byte {
	// Simplified version - actual implementation should include proper decryption
	// TODO: Implement actual decryption algorithm from legacy code
	return data
}

// GameAccountService provides high-level game account operations
type GameAccountService struct {
	validator *GameAccountValidator
}

// NewGameAccountService creates a new game account service
func NewGameAccountService(gameServerAddr string) *GameAccountService {
	return &GameAccountService{
		validator: NewGameAccountValidator(gameServerAddr),
	}
}

// VerifyAndGetInfo verifies game account and retrieves account information
func (s *GameAccountService) VerifyAndGetInfo(ctx context.Context, account, password string) (*GameAccountInfo, error) {
	// Try mobile validation first (most common)
	info, err := s.validator.ValidateByMobile(ctx, account, password)
	if err == nil {
		return info, nil
	}

	// If mobile validation fails, try account validation
	info, err = s.validator.ValidateByAccount(ctx, account, password)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// GetNickname retrieves nickname for a game account
func (s *GameAccountService) GetNickname(ctx context.Context, account, password string) (string, error) {
	info, err := s.VerifyAndGetInfo(ctx, account, password)
	if err != nil {
		return "", err
	}

	return info.Nickname, nil
}

// ValidateCredentials validates game account credentials
func (s *GameAccountService) ValidateCredentials(ctx context.Context, account, password string) (bool, error) {
	info, err := s.VerifyAndGetInfo(ctx, account, password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return false, nil
		}
		return false, err
	}

	return info.Success, nil
}

// NOTE: This is a simplified implementation for demonstration purposes.
// The actual implementation should:
// 1. Port the complete protocol handling from passing-dragonfly/waiter/plaza/
// 2. Implement proper encryption/decryption algorithms
// 3. Handle all protocol message types
// 4. Implement proper error handling and retry logic
// 5. Add connection pooling for better performance
// 6. Add proper logging and monitoring
//
// Reference files from legacy project:
// - passing-dragonfly/waiter/plaza/session.go
// - passing-dragonfly/waiter/plaza/mobilelogon.go
// - passing-dragonfly/waiter/plaza/encoder.go
// - passing-dragonfly/waiter/plaza/const.go
