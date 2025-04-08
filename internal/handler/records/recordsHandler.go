package records

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/SRIRAMGJ007/Health-Sync/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

var encryptionKeyStr []byte

type GetEncryptedFileParams struct {
	UserID pgtype.UUID
	ID     pgtype.UUID
}

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}
	encryptionKeyStr = []byte(os.Getenv("ENCRYPTION_KEY")) // Initialize
}

func UploadMedicalRecord(ctx *gin.Context, queries *repository.Queries) {
	userID := ctx.Param("userid")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userid missing"})
		log.Println("update user profile failed -> query parameter id missing, bad request")
		return
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("update user profile failed -> invalid userid, bad request: { %v }", err)
		return
	}

	pgUUID := pgtype.UUID{Bytes: parsedID, Valid: true}

	// Get the file from the request
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Open and read the uploaded file
	fileData, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileData.Close()

	data, err := io.ReadAll(fileData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	if len(encryptionKeyStr) != 32 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid encryption key length"})
		return
	}

	encryptionKey := []byte(encryptionKeyStr)
	// Encrypt file data
	encryptedData, err := utils.EncryptData(data, encryptionKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
		return
	}

	// Store in the database using SQLC
	_, err = queries.StoreEncryptedFile(ctx, repository.StoreEncryptedFileParams{
		UserID:   pgUUID, // Using parsed pgtype.UUID
		FileName: file.Filename,
		FileData: encryptedData,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store record"})
		return
	}

	// Respond with success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Medical record uploaded successfully"})
}

func ListMedicalRecords(ctx *gin.Context, queries *repository.Queries) {
	userID := ctx.Param("userid")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID missing"})
		return
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pgUUID := pgtype.UUID{Bytes: parsedID, Valid: true}
	records, err := queries.GetUserFiles(ctx, pgUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	if len(records) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No records found"})
		return
	}

	var fileList []gin.H
	for _, record := range records {
		fileList = append(fileList, gin.H{
			"file_id":    record.ID, // Assuming you have a unique FileID column
			"file_name":  record.FileName,
			"created_at": record.CreatedAt, // If you have a timestamp field
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userID, "records": fileList})
}

func DownloadMedicalRecord(ctx *gin.Context, queries *repository.Queries) {
	userID := ctx.Param("userid")
	fileID := ctx.Param("fileid") // Expecting fileid as a query parameter
	log.Println(userID)
	log.Println(fileID)
	if userID == "" || fileID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID and File ID are required"})
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	parsedFileID, err := uuid.Parse(fileID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}
	pgFileID := pgtype.UUID{Bytes: parsedFileID, Valid: true}

	params := GetEncryptedFileParams{
		ID:     pgFileID,
		UserID: pgUserID,
	}

	record, err := queries.GetEncryptedFile(ctx, repository.GetEncryptedFileParams(params)) // Fetch a specific file
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	encryptionKey := []byte(encryptionKeyStr)

	decryptedData, err := utils.DecryptData(record.FileData, encryptionKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename="+record.FileName)
	ctx.Header("Content-Type", http.DetectContentType(decryptedData))
	ctx.Data(http.StatusOK, "application/octet-stream", decryptedData)
}

func ViewMedicalRecord(ctx *gin.Context, queries *repository.Queries) {
	userID := ctx.Param("userid")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID missing"})
		log.Println("View medical record failed -> User ID missing")
		return
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("View medical record failed -> Invalid user ID: %v", err)
		return
	}

	pgUUID := pgtype.UUID{Bytes: parsedID, Valid: true}
	records, err := queries.GetUserFiles(ctx, pgUUID)
	if err != nil || len(records) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No records found"})
		log.Printf("View medical record failed -> No records found: %v", err)
		return
	}

	record := records[0] // Modify this if you want to allow users to select a file

	if len(encryptionKeyStr) != 32 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid encryption key length"})
		return
	}
	encryptionKey := []byte(encryptionKeyStr)

	decryptedData, err := utils.DecryptData(record.FileData, encryptionKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
		log.Printf("Decryption failed: %v", err)
		return
	}

	// Detect MIME type based on file extension
	mimeType := "application/octet-stream" // Default
	ext := strings.ToLower(filepath.Ext(record.FileName))
	switch ext {
	case ".pdf":
		mimeType = "application/pdf"
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	}

	// Set headers for inline display
	ctx.Header("Content-Type", mimeType)
	ctx.Header("Content-Disposition", "inline; filename="+record.FileName)
	ctx.Data(http.StatusOK, mimeType, decryptedData)
}
