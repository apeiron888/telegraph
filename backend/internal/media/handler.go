package media

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MediaType represents the type of media
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeDocument MediaType = "document"
)

// FileMetadata represents uploaded file information
type FileMetadata struct {
	ID          string    `json:"id" bson:"_id"`
	FileName    string    `json:"file_name" bson:"file_name"`
	ContentType string    `json:"content_type" bson:"content_type"`
	Size        int64     `json:"size" bson:"size"`
	MD5Hash     string    `json:"md5_hash" bson:"md5_hash"`
	MediaType   MediaType `json:"media_type" bson:"media_type"`
	UploaderID  string    `json:"uploader_id" bson:"uploader_id"`
	StoragePath string    `json:"storage_path" bson:"storage_path"`
	URL         string    `json:"url" bson:"url"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

// MediaHandler handles file storage operations
type MediaHandler interface {
	Upload(file io.Reader, filename, contentType, uploaderID string, size int64) (*FileMetadata, error)
	GetURL(fileID string) string
	Delete(fileID string) error
	Serve(fileID string) (string, error)
}

// LocalMediaHandler implements MediaHandler using local filesystem
type LocalMediaHandler struct {
	baseDir string
	baseURL string
}

func NewLocalMediaHandler(baseDir, baseURL string) *LocalMediaHandler {
	// Create base directory if it doesn't exist
	os.MkdirAll(baseDir, 0755)
	return &LocalMediaHandler{
		baseDir: baseDir,
		baseURL: baseURL,
	}
}

func (h *LocalMediaHandler) Upload(file io.Reader, filename, contentType, uploaderID string, size int64) (*FileMetadata, error) {
	// Generate unique ID
	fileID := uuid.New().String()

	// Determine media type
	mediaType := determineMediaType(contentType)

	// Create subdirectory for media type
	mediaDir := filepath.Join(h.baseDir, string(mediaType))
	os.MkdirAll(mediaDir, 0755)

	// Generate storage path
	ext := filepath.Ext(filename)
	storagePath := filepath.Join(mediaDir, fileID+ext)

	// Create file
	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Calculate MD5 while copying
	hash := md5.New()
	multiWriter := io.MultiWriter(dst, hash)

	written, err := io.Copy(multiWriter, file)
	if err != nil {
		os.Remove(storagePath)
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	md5Hash := hex.EncodeToString(hash.Sum(nil))

	metadata := &FileMetadata{
		ID:          fileID,
		FileName:    filename,
		ContentType: contentType,
		Size:        written,
		MD5Hash:     md5Hash,
		MediaType:   mediaType,
		UploaderID:  uploaderID,
		StoragePath: storagePath,
		URL:         h.GetURL(fileID + ext),
		CreatedAt:   time.Now(),
	}

	return metadata, nil
}

func (h *LocalMediaHandler) GetURL(fileID string) string {
	return fmt.Sprintf("%s/media/%s", h.baseURL, fileID)
}

func (h *LocalMediaHandler) Delete(fileID string) error {
	// Find file in all media type directories
	mediaTypes := []MediaType{MediaTypeImage, MediaTypeVideo, MediaTypeAudio, MediaTypeDocument}
	
	for _, mediaType := range mediaTypes {
		pattern := filepath.Join(h.baseDir, string(mediaType), fileID+"*")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		
		for _, match := range matches {
			if err := os.Remove(match); err != nil {
				return fmt.Errorf("failed to delete file: %w", err)
			}
			return nil
		}
	}
	
	return fmt.Errorf("file not found")
}

func (h *LocalMediaHandler) Serve(fileID string) (string, error) {
	// Find file in all media type directories
	mediaTypes := []MediaType{MediaTypeImage, MediaTypeVideo, MediaTypeAudio, MediaTypeDocument}
	
	for _, mediaType := range mediaTypes {
		pattern := filepath.Join(h.baseDir, string(mediaType), fileID+"*")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		
		if len(matches) > 0 {
			return matches[0], nil
		}
	}
	
	return "", fmt.Errorf("file not found")
}

func determineMediaType(contentType string) MediaType {
	contentType = strings.ToLower(contentType)
	
	if strings.HasPrefix(contentType, "image/") {
		return MediaTypeImage
	} else if strings.HasPrefix(contentType, "video/") {
		return MediaTypeVideo
	} else if strings.HasPrefix(contentType, "audio/") {
		return MediaTypeAudio
	}
	
	return MediaTypeDocument
}
