package core

import "testing"

func TestGetMapValue(t *testing.T) {
	t.Run("ExistingKey", func(t *testing.T) {
		m := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		}

		result := GetMapValue(m, "key1", "default")
		if result != "value1" {
			t.Errorf("Expected 'value1', got %v", result)
		}
	})

	t.Run("NonExistingKey", func(t *testing.T) {
		m := map[string]interface{}{
			"key1": "value1",
		}

		result := GetMapValue(m, "nonexistent", "default")
		if result != "default" {
			t.Errorf("Expected 'default', got %v", result)
		}
	})

	t.Run("EmptyMap", func(t *testing.T) {
		m := make(map[string]interface{})

		result := GetMapValue(m, "anykey", "default")
		if result != "default" {
			t.Errorf("Expected 'default', got %v", result)
		}
	})

	t.Run("NilMap", func(t *testing.T) {
		var m map[string]interface{}

		result := GetMapValue(m, "anykey", "default")
		if result != "default" {
			t.Errorf("Expected 'default', got %v", result)
		}
	})

	t.Run("DifferentTypes", func(t *testing.T) {
		m := map[string]interface{}{
			"string": "text",
			"int":    123,
			"bool":   false,
			"nil":    nil,
		}

		tests := []struct {
			key      string
			expected interface{}
		}{
			{"string", "text"},
			{"int", 123},
			{"bool", false},
			{"nil", nil},
		}

		for _, test := range tests {
			result := GetMapValue(m, test.key, "default")
			if result != test.expected {
				t.Errorf("Key %s: expected %v, got %v", test.key, test.expected, result)
			}
		}
	})
}
