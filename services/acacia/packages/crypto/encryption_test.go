package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEncryptionService(t *testing.T) {
	t.Run("valid 32-byte key", func(t *testing.T) {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err, "failed to generate random key")

		service, err := NewEncryptionService(key)
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("invalid key size - too short", func(t *testing.T) {
		key := make([]byte, 16)
		service, err := NewEncryptionService(key)
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "32 bytes")
	})

	t.Run("invalid key size - too long", func(t *testing.T) {
		key := make([]byte, 64)
		service, err := NewEncryptionService(key)
		assert.Error(t, err)
		assert.Nil(t, service)
	})
}

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err, "failed to generate random key")

	service, err := NewEncryptionService(key)
	require.NoError(t, err, "failed to create service")

	testCases := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple API key",
			plaintext: "sk-ant-api03-1234567890abcdef",
		},
		{
			name:      "long API key",
			plaintext: "sk-ant-api03-1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "special characters",
			plaintext: "sk-ant-api03-!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:      "unicode characters",
			plaintext: "sk-ant-api03-测试-тест-🔐",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := service.Encrypt(tc.plaintext)
			require.NoError(t, err, "encryption failed")

			// Verify encryption happened (encrypted != plaintext)
			if tc.plaintext != "" {
				assert.NotEqual(t, tc.plaintext, encrypted, "encrypted text should be different from plaintext")
			}

			// Decrypt
			decrypted, err := service.Decrypt(encrypted)
			require.NoError(t, err, "decryption failed")

			// Verify roundtrip
			assert.Equal(t, tc.plaintext, decrypted, "decrypted text should match original")
		})
	}
}

func TestEncryptionProducesDifferentCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err, "failed to generate random key")

	service, err := NewEncryptionService(key)
	require.NoError(t, err, "failed to create service")

	plaintext := "sk-ant-api03-test-key"

	// Encrypt same plaintext twice
	encrypted1, err := service.Encrypt(plaintext)
	require.NoError(t, err, "first encryption failed")

	encrypted2, err := service.Encrypt(plaintext)
	require.NoError(t, err, "second encryption failed")

	// Ciphertexts should be different (due to random nonce)
	assert.NotEqual(t, encrypted1, encrypted2, "encrypting same plaintext twice should produce different ciphertexts")

	// But both should decrypt to the same plaintext
	decrypted1, err := service.Decrypt(encrypted1)
	require.NoError(t, err, "first decryption failed")

	decrypted2, err := service.Decrypt(encrypted2)
	require.NoError(t, err, "second decryption failed")

	assert.Equal(t, plaintext, decrypted1, "first ciphertext should decrypt to original plaintext")
	assert.Equal(t, plaintext, decrypted2, "second ciphertext should decrypt to original plaintext")
}

func TestDecryptInvalidData(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err, "failed to generate random key")

	service, err := NewEncryptionService(key)
	require.NoError(t, err, "failed to create service")

	t.Run("invalid base64", func(t *testing.T) {
		_, err := service.Decrypt("not-valid-base64!!!")
		assert.Error(t, err, "should error for invalid base64")
	})

	t.Run("ciphertext too short", func(t *testing.T) {
		_, err := service.Decrypt("YWJj") // "abc" in base64
		assert.Error(t, err, "should error for short ciphertext")
		assert.Contains(t, err.Error(), "too short")
	})

	t.Run("corrupted ciphertext", func(t *testing.T) {
		plaintext := "test-api-key"
		encrypted, err := service.Encrypt(plaintext)
		require.NoError(t, err, "encryption failed")

		// Corrupt the encrypted data
		corrupted := encrypted[:len(encrypted)-5] + "XXXXX"

		_, err = service.Decrypt(corrupted)
		assert.Error(t, err, "should error for corrupted ciphertext")
	})
}

func TestDecryptWithWrongKey(t *testing.T) {
	// Create two different keys
	key1 := make([]byte, 32)
	_, err := rand.Read(key1)
	require.NoError(t, err, "failed to generate first key")

	key2 := make([]byte, 32)
	_, err = rand.Read(key2)
	require.NoError(t, err, "failed to generate second key")

	service1, err := NewEncryptionService(key1)
	require.NoError(t, err, "failed to create first service")

	service2, err := NewEncryptionService(key2)
	require.NoError(t, err, "failed to create second service")

	plaintext := "sk-ant-api03-secret-key"

	// Encrypt with first key
	encrypted, err := service1.Encrypt(plaintext)
	require.NoError(t, err, "encryption failed")

	// Try to decrypt with second key (should fail)
	_, err = service2.Decrypt(encrypted)
	assert.Error(t, err, "should error when decrypting with wrong key")
}
