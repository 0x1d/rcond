package user

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// AddAuthorizedKey verifies and adds an SSH public key to /home/<user>/.ssh/authorized_keys
// if it doesn't already exist. Returns the key's fingerprint.
func AddAuthorizedKey(user string, pubKey string) (string, error) {
	// Verify the public key format and get fingerprint
	parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKey))
	if err != nil {
		return "", fmt.Errorf("invalid SSH public key: %v", err)
	}
	fingerprint := ssh.FingerprintSHA256(parsed)

	// Ensure .ssh directory exists
	sshDir := fmt.Sprintf("/home/%s/.ssh", user)
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create .ssh directory: %v", err)
	}

	// Check if key already exists
	keyFile := fmt.Sprintf("%s/authorized_keys", sshDir)
	if _, err := os.Stat(keyFile); err == nil {
		existingKeys, err := os.ReadFile(keyFile)
		if err != nil {
			return "", fmt.Errorf("failed to read authorized_keys: %v", err)
		}
		if string(existingKeys) != "" {
			for _, line := range strings.Split(string(existingKeys), "\n") {
				if line == "" {
					continue
				}
				parsed, err := ssh.ParsePublicKey([]byte(line))
				if err != nil {
					continue
				}
				if ssh.FingerprintSHA256(parsed) == fingerprint {
					// Key already exists, return fingerprint
					return fingerprint, nil
				}
			}
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check authorized_keys: %v", err)
	}

	// Open authorized_keys file in append mode
	f, err := os.OpenFile(keyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to open authorized_keys: %v", err)
	}
	defer f.Close()

	// Write the public key
	if _, err := f.WriteString(pubKey + "\n"); err != nil {
		return "", fmt.Errorf("failed to write public key: %v", err)
	}

	return fingerprint, nil
}

// RemoveAuthorizedKey removes an authorized SSH key from /home/<user>/.ssh/authorized_keys
// using the key's fingerprint to identify which key to remove
func RemoveAuthorizedKey(user string, fingerprint string) error {
	// Check if authorized_keys file exists
	keyFile := fmt.Sprintf("/home/%s/.ssh/authorized_keys", user)
	if _, err := os.Stat(keyFile); err != nil {
		if os.IsNotExist(err) {
			return nil // Nothing to remove
		}
		return fmt.Errorf("failed to check authorized_keys: %v", err)
	}

	// Read existing keys
	existingKeys, err := os.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("failed to read authorized_keys: %v", err)
	}

	// Filter out the key with matching fingerprint
	var newLines []string
	for _, line := range strings.Split(string(existingKeys), "\n") {
		if line == "" {
			continue
		}
		parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
		if err != nil {
			// Keep lines we can't parse
			newLines = append(newLines, line)
			continue
		}
		if ssh.FingerprintSHA256(parsed) != fingerprint {
			newLines = append(newLines, line)
		}
	}

	// Write back the filtered keys
	err = os.WriteFile(keyFile, []byte(strings.Join(newLines, "\n")+"\n"), 0600)
	if err != nil {
		return fmt.Errorf("failed to write authorized_keys: %v", err)
	}

	return nil
}
