package helpers

import (
	"Komentory/api/app/models"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/h2non/filetype"
)

// GetLocalFileInfo func for get local file's info (Content-Type, extension, size, etc).
func GetLocalFileInfo(pathToFile, fileType string) (*models.LocalFileInfo, error) {
	// Get file size.
	fileSize, err := GetFileSize(pathToFile)
	if err != nil {
		return nil, err
	}

	// Define maximum file size in bytes.
	maxFileSize, err := strconv.ParseInt(os.Getenv("MAX_UPLOAD_FILE_SIZE"), 10, 64)
	if err != nil {
		return nil, err
	}

	// If actual file size is greater than max file size, throw error.
	if fileSize > maxFileSize {
		return nil, fmt.Errorf("file is too large for upload (%d bytes)", fileSize)
	}

	// Read given file from file system.
	buf, err := ioutil.ReadFile(filepath.Clean(pathToFile))
	if err != nil {
		return nil, err
	}

	// Create matching for file in buffer.
	kind, err := filetype.Match(buf)
	if err != nil {
		return nil, err
	}

	// Check, if given file has an unknown type.
	if kind == filetype.Unknown {
		return nil, fmt.Errorf("file has an unknown type")
	}

	// Switch file types.
	switch fileType {
	case "image":
		// Check, if given file is image.
		if !filetype.IsImage(buf) {
			return nil, fmt.Errorf("only images are supported")
		}

		// Check, if given image is JPG, PNG, or SVG.
		if kind.Extension != "jpg" && kind.Extension != "png" && kind.Extension != "svg" {
			return nil, fmt.Errorf("only images with *.jpg, *.png, or *.svg extensions are supported")
		}
	case "document":
		// Check, if file is document.
		if !filetype.IsDocument(buf) {
			return nil, fmt.Errorf("only documents are supported")
		}

		// Check, if given document is PDF.
		if kind.Extension != "pdf" {
			return nil, fmt.Errorf("only documents with *.pdf extension is supported")
		}
	default:
		// Throw error, if file is not supported.
		return nil, fmt.Errorf("wrong or unsupported file type (%s)", fileType)
	}

	// Return file info.
	return &models.LocalFileInfo{
		ContentType: kind.MIME.Value,
		Extension:   kind.Extension,
		Size:        fileSize,
	}, nil
}

// GenerateFileName func for generate name for uploaded file.
func GenerateFileName(userID string) (string, error) {
	// Create a new SHA256 hash.
	sha256 := sha256.New()

	// Create a new string with user ID and time string.
	name := userID + time.Now().String()

	// See: https://pkg.go.dev/io#Writer.Write
	_, err := sha256.Write([]byte(name))
	if err != nil {
		// Return error, if generation failed.
		return "", err
	}

	// Return a new file name.
	return hex.EncodeToString(sha256.Sum(nil)), nil
}

// GetFileSize func for getting the file size.
func GetFileSize(pathToFile string) (int64, error) {
	// Get file from system path.
	file, err := os.Open(filepath.Clean(pathToFile))
	if err != nil {
		return 0, err
	}

	// Get file statistic.
	fileStat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	// Check, if file size is zero.
	if fileStat.Size() == 0 {
		// Return error message.
		return 0, fmt.Errorf("file have no size (zero bytes)")
	}

	// Return file size in bytes.
	return fileStat.Size(), nil
}
