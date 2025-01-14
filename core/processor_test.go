package core

import (
	"errors"
	"testing"
	"time"

	"gitub.com/sriramr98/go_kvdb/store"
)

func TestGETProcessor(t *testing.T) {
	mockStore := store.NewInMemoryStore()
	cmd_processor := CommandProcessor{Store: mockStore}

	mockStore.Set("test1", store.NewDataObj([]byte("test1"), -1))
	mockStore.Set("test_expired", store.NewDataObj([]byte("test2"), int(time.Now().Unix())))

	testCases := map[string]struct {
		key         string
		expected    string
		shouldError bool
		expectedErr error
	}{
		"test GET without a key present": {
			key:         "test",
			shouldError: true,
			expectedErr: store.ErrKeyNotFound,
		},

		"test GET with a key present": {
			key:         "test1",
			expected:    "test1",
			shouldError: false,
		},
		"test GET with expired key": {
			key:         "test_expired",
			shouldError: true,
			expectedErr: store.ErrKeyNotFound,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			req := Request{Command: CMDGet, Params: []string{testCase.key}}
			res, err := cmd_processor.Process(req)

			if testCase.shouldError {
				if err == nil {
					t.Errorf("Expected error %v got nil", testCase.expectedErr)
				}
				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("Expected error %v got %v", testCase.expected, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected value but got error %v", err)
				}

				if !res.Success {
					t.Error("Expected successful response, but for Success: false")
				}

				if string(res.Value) != testCase.expected {
					t.Errorf("Expected value %s but got %s", testCase.expected, string(res.Value))
				}
			}
		})
	}
}

func TestSETProcessor(t *testing.T) {
	mockStore := store.NewInMemoryStore()
	cmd_processor := CommandProcessor{Store: mockStore}

	testCases := map[string]struct {
		key         string
		value       string
		expirySec   string
		shouldError bool
		expectedErr error
	}{
		"test SET without a key present": {
			key:         "test1",
			value:       "test1",
			shouldError: false,
		},

		"test SET with a key present": {
			key:         "test1",
			value:       "test2",
			shouldError: false,
		},

		"test SET with valid expiry": {
			key:         "test_expiry",
			value:       "exp_value",
			expirySec:   "10",
			shouldError: false,
		},

		"test SET with invalid expiry (expiry seconds < 0)": {
			key:         "test_exp",
			value:       "exp_val",
			expirySec:   "-10",
			shouldError: true,
			expectedErr: ErrExpiryLessThanZero,
		},

		"test SET with invalid expiry (expiry not a number)": {
			key:         "test_exp",
			value:       "exp_val",
			expirySec:   "abc",
			shouldError: true,
			expectedErr: ErrExpiryNotANumber,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			params := []string{testCase.key, testCase.value}
			if testCase.expirySec != "" {
				params = append(params, testCase.expirySec)
			}

			req := Request{Command: CMDSet, Params: params}
			res, err := cmd_processor.Process(req)

			if testCase.shouldError {
				if err == nil {
					t.Errorf("Expected error %v got nil", testCase.expectedErr)
				}
				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("Expected error %v got %v", testCase.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected value but got error %v", err)
				}

				if !res.Success {
					t.Error("Expected successful response, but for Success: false")
				}

				storeVal, err := mockStore.Get(testCase.key)
				if err != nil {
					t.Errorf("Expected key %s to be present at store but got error %v", testCase.key, err)
				}

				if string(storeVal.Value()) != testCase.value {
					t.Errorf("Expected value %s but got %s", testCase.value, string(storeVal.Value()))
				}

				if testCase.expirySec == "" && storeVal.ExpiresAt() != -1 {
					t.Error("Expected expires at to be -1 for SET without expiry but expiry set")
				}

				if testCase.expirySec != "" && storeVal.ExpiresAt() == -1 {
					t.Errorf("Expected expiresAt to be present at value for key %s but expiry not set", testCase.key)
				}
			}
		})
	}
}

func TestDELETEProcessor(t *testing.T) {
	mockStore := store.NewInMemoryStore()
	cmd_processor := CommandProcessor{Store: mockStore}

	mockStore.Set("test1", store.NewDataObj([]byte("test1"), 0))

	testCases := map[string]struct {
		key string
	}{
		"test DELETE without a key present": {
			key: "test2",
		},

		"test DELETE with a key present": {
			key: "test1",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			req := Request{Command: CMDDel, Params: []string{testCase.key}}

			res, err := cmd_processor.Process(req)
			if err != nil {
				t.Errorf("Expected no error but got error %v", err)
			}

			if !res.Success {
				t.Error("Expected success but got Success: false")
			}
		})
	}
}
