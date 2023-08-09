package telegraph

import (
	"testing"
)

func TestContentFormat(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		_, err := ContentFormat(42)
		if ErrInvalidDataType != err {
			t.Errorf("error must be: %v", ErrInvalidDataType)
		}
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			validContentDOM, err := ContentFormat(`<p>Hello, World!</p>`)
			if err != nil {
				t.Error(err)
			}
			if validContentDOM == nil {
				t.Error("content is invalid")
			}
		})
		t.Run("bytes", func(t *testing.T) {
			validContentDOM, err := ContentFormat([]byte(`<p>Hello, World!</p>`))
			if err != nil {
				t.Error(err)
			}
			if validContentDOM == nil {
				t.Error("content is invalid")
			}
		})
	})
}
