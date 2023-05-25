package protocol

import (
	"errors"
	"testing"
)

func TestProtocol(t *testing.T) {

	t.Run("isParamsValid", func(t *testing.T) {
		// It's fine not testing all commands since they all use the same method
		t.Run("should return true", func(t *testing.T) {
			cmd := CMDGet
			params := []string{"key"}
			if !cmd.isParamsValid(params) {
				t.Errorf("Expected true, got false")
			}
		})

		t.Run("should return false", func(t *testing.T) {
			cmd := CMDGet
			params := []string{}
			if cmd.isParamsValid(params) {
				t.Errorf("Expected false, got true")
			}
		})

	})

	t.Run("parseCommand", func(t *testing.T) {

		t.Run("should return CMDGet", func(t *testing.T) {
			cmd, err := parseCommand("GET")
			if cmd != CMDGet {
				t.Errorf("Expected %s, got %s", CMDGet.Op, cmd.Op)
			}
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
		})

		t.Run("should return CMDSet", func(t *testing.T) {
			cmd, err := parseCommand("SET")
			if cmd != CMDSet {
				t.Errorf("Expected %s, got %s", CMDSet.Op, cmd.Op)
			}
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
		})

		t.Run("should return CMDDel", func(t *testing.T) {
			cmd, err := parseCommand("DEL")
			if cmd != CMDDel {
				t.Errorf("Expected %s, got %s", CMDDel.Op, cmd.Op)
			}
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
		})

		t.Run("should return CMDPing", func(t *testing.T) {
			cmd, err := parseCommand("PING")
			if cmd != CMDPing {
				t.Errorf("Expected %s, got %s", CMDPing.Op, cmd.Op)
			}
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
		})

		t.Run("should return empty command", func(t *testing.T) {
			cmd, err := parseCommand("INVALID")
			if cmd != (Command{}) {
				t.Errorf("Expected empty command, got %s", cmd.Op)
			}
			if !errors.Is(err, ErrInvalidCommand) {
				t.Errorf("Expected %s, got %s", ErrInvalidCommand, err)
			}
		})

	})

	t.Run("Parse", func(t *testing.T) {

		p := ClientProtocol{}

		t.Run("should return error", func(t *testing.T) {
			_, err := p.Parse("")
			if err == nil {
				t.Errorf("Expected error, got nil")
			}

			if !errors.Is(err, ErrInvalidCommand) {
				t.Errorf("Expected %s, got %s", ErrInvalidCommand, err)
			}
		})

		t.Run("parses valid command", func(t *testing.T) {
			req, err := p.Parse("GET key")
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
			if req.Command != CMDGet {
				t.Errorf("Expected %s, got %s", CMDGet.Op, req.Command.Op)
			}
			if len(req.Params) != 1 {
				t.Errorf("Expected 1, got %d", len(req.Params))
			}
			if req.Params[0] != "key" {
				t.Errorf("Expected key, got %s", req.Params[0])
			}
		})

		t.Run("parses valid SET command with optional TTL", func(t *testing.T) {
			req, err := p.Parse("SET key value 100")
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
			if req.Command != CMDSet {
				t.Errorf("Expected %s, got %s", CMDSet.Op, req.Command.Op)
			}
			if len(req.Params) != 3 {
				t.Errorf("Expected 3, got %d", len(req.Params))
			}
			if req.Params[0] != "key" {
				t.Errorf("Expected key, got %s", req.Params[0])
			}
			if req.Params[1] != "value" {
				t.Errorf("Expected value, got %s", req.Params[1])
			}
			if req.Params[2] != "100" {
				t.Errorf("Expected 100, got %s", req.Params[2])
			}
		})

	})

}
