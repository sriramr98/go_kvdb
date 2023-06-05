package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

		p := ClientProtocol{}

		t.Run("should return CMDGet", func(t *testing.T) {
			cmd, err := p.extractCommand("GET")
			assert.Equal(t, CMDGet, cmd)
			assert.NoError(t, err)
		})

		t.Run("should return CMDSet", func(t *testing.T) {
			cmd, err := p.extractCommand("SET")
			assert.Equal(t, CMDSet, cmd)
			assert.NoError(t, err)
		})

		t.Run("should return CMDDel", func(t *testing.T) {
			cmd, err := p.extractCommand("DEL")
			assert.Equal(t, CMDDel, cmd)
			assert.NoError(t, err)
		})

		t.Run("should return CMDPing", func(t *testing.T) {
			cmd, err := p.extractCommand("PING")
			assert.Equal(t, CMDPing, cmd)
			assert.NoError(t, err)
		})

		t.Run("should return empty command", func(t *testing.T) {
			cmd, err := p.extractCommand("INVALID")
			assert.Equal(t, Command{}, cmd)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidCommand)
		})

	})

	t.Run("Parse", func(t *testing.T) {

		p := ClientProtocol{}

		t.Run("should return error", func(t *testing.T) {
			_, err := p.Parse("")

			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidCommand)
		})

		t.Run("parses valid command", func(t *testing.T) {
			req, err := p.Parse("GET key")
			assert.NoError(t, err)
			assert.Equal(t, CMDGet, req.Command)
			assert.Equal(t, []string{"key"}, req.Params)
			assert.Equal(t, len(req.Params), 1)
		})

		t.Run("parses valid SET command with optional TTL", func(t *testing.T) {
			req, err := p.Parse("SET key value 100")
			assert.NoError(t, err)
			assert.Equal(t, CMDSet, req.Command)
			assert.Equal(t, []string{"key", "value", "100"}, req.Params)
			assert.Equal(t, len(req.Params), 3)

		})

	})

}
