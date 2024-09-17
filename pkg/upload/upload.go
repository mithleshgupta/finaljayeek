package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// UploadFile saves the uploaded file to the specified directory with a generated filename
// based on the current timestamp and returns file information.
func UploadFile(file *multipart.FileHeader, directory string) (os.FileInfo, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Generate a unique filename based on the current timestamp
	filename := generateFilename(file.Filename)

	dst, err := os.Create(filepath.Join(directory, filename))
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Retrieve file information
	fileInfo, err := dst.Stat()
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

func generateFilename(originalFilename string) string {
	// Generate a unique identifier using a cryptographic hash function
	hash := generateHash(originalFilename)

	// Get the file extension from the original filename
	fileExtension := filepath.Ext(originalFilename)

	// Combine the hash and file extension to create the secure filename
	filename := hash + fileExtension

	return filename
}

// Helper function to generate a cryptographic hash of the input string
func generateHash(input string) string {
	// Create a new SHA256 hash instance
	hash := sha256.New()

	// Convert the input string to bytes and hash it
	hash.Write([]byte(input))

	// Get the resulting hash as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the hash bytes to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
