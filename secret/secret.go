package secret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Store struct {
	Filepath string
}

const environKey = "ENCRYPTION_KEY"

func (s Store) Get(name string) (string, error) {
	pairs, err := s.getPairs()
	if err != nil {
		return "", err
	}

	if _, ok := pairs[name]; !ok {
		return "", fmt.Errorf("key %q not found", name)
	}

	return pairs[name], nil
}

func (s Store) Set(name, val string) error {
	pairs, err := s.getPairs()
	if err != nil {
		return err
	}

	pairs[name] = val

	return s.savePairs(pairs)
}

func (s Store) Delete(name string) error {
	pairs, err := s.getPairs()
	if err != nil {
		return err
	}

	if _, ok := pairs[name]; !ok {
		return fmt.Errorf("key %q does not exist", name)
	}

	delete(pairs, name)

	return s.savePairs(pairs)
}

func (s Store) getPairs() (map[string]string, error) {
	file, err := os.OpenFile(s.Filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	decrypted, err := decrypt(data)
	if err != nil {
		return nil, err
	}

	return decodePairs(decrypted)
}

func (s Store) savePairs(pairs map[string]string) error {
	encoded, err := encodePairs(pairs)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(encoded)
	if err != nil {
		return err
	}

	file, err := os.Create(s.Filepath)
	if err != nil {
		return err
	}

	_, err = file.Write(encrypted)
	return err
}

func decodePairs(data []byte) (map[string]string, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	pairs := make(map[string]string, len(records))
	for _, record := range records {
		if len(record) != 2 {
			return nil, fmt.Errorf("corrupted file")
		}
		pairs[record[0]] = record[1]
	}

	return pairs, nil
}

func encodePairs(pairs map[string]string) ([]byte, error) {
	records := make([][]string, 0, len(pairs))

	for k, v := range pairs {
		records = append(records, []string{k, v})
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func encrypt(data []byte) ([]byte, error) {
	key, found := os.LookupEnv(environKey)
	if !found {
		return nil, fmt.Errorf("envrionment variable %q not found", environKey)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err // Error if the key length is invalid (16, 24, 32 bytes for AES)
	}

	// Step 2: Create a GCM mode instance for the block cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err // Error if the block cipher doesn't support GCM
	}

	// Step 3: Generate a random nonce (unique per encryption)
	nonce := make([]byte, gcm.NonceSize()) // GCM requires a specific nonce size (12 bytes recommended)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err // Error if randomness fails
	}

	// Step 4: Encrypt the plaintext and append an authentication tag
	// `Seal` encrypts and adds the authentication tag. No additional data (AAD) is used here.
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Step 5: Combine the nonce and ciphertext
	return append(nonce, ciphertext...), nil
}

func decrypt(data []byte) ([]byte, error) {
	key, found := os.LookupEnv(environKey)
	if !found {
		return nil, fmt.Errorf("envrionment variable %q not found", environKey)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, err // Ciphertext too short
	}

	nonce, encrypted := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, encrypted, nil)
}
