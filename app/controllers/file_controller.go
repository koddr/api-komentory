package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/pkg/helpers"
	"Komentory/api/platform/cdn"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// GetFileListFromCDN func for return a list of files from CDN.
func GetFileListFromCDN(c *fiber.Ctx) error {
	// Get claims from JWT.
	_, err := utilities.ExtractTokenMetaData(c)
	if err != nil {
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Create CDN connection.
	connDOSpaces, err := cdn.DOSpacesConnection()
	if err != nil {
		return utilities.CheckForError(c, err, 500, "cdn", err.Error())
	}

	// Create context with cancel.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // auto close

	// Get list of files from CDN.
	listObjectsChannel := connDOSpaces.ListObjects(
		ctx,
		os.Getenv("DO_SPACES_BUCKET_NAME"),
		minio.ListObjectsOptions{
			Prefix:    os.Getenv("DO_SPACES_UPLOADS_FOLDER_NAME"),
			Recursive: true,
		},
	)

	// Define File struct for object list.
	objects := []*models.FileFromCDN{}

	// Range object list from CDN for create a new Go object for JSON serialization.
	for object := range listObjectsChannel {
		// Check, if received object is valid.
		if object.Err != nil {
			return utilities.CheckForError(c, err, 400, "cdn object", object.Err.Error())
		}

		// Skip upload folder from list, only files.
		if !strings.HasSuffix(object.Key, "/") {
			// Create a new File struct from object info.
			file := &models.FileFromCDN{
				Key:       object.Key,
				ETag:      object.ETag,
				VersionID: object.VersionID,
				URL:       fmt.Sprintf("%v/%v", os.Getenv("CDN_PUBLIC_URL"), object.Key),
			}

			// Add this file to objects list.
			objects = append(objects, file)
		}
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"count":   len(objects),
		"objects": objects,
	})
}

// PutFileToCDN func for upload a file to CDN.
// Allowed types: image, document.
func PutFileToCDN(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utilities.ExtractTokenMetaData(c)
	if err != nil {
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Define user ID.
	userID := claims.UserID.String()

	// Create new LocalFile struct.
	localFile := &models.LocalFile{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(localFile); err != nil {
		return utilities.CheckForError(c, err, 400, "local file", err.Error())
	}

	// Create a new DO Spaces connection.
	connDOSpaces, err := cdn.DOSpacesConnection()
	if err != nil {
		return utilities.CheckForError(c, err, 500, "cdn", err.Error())
	}

	// Upload file process.
	uploadFileInfo, err := cdn.UploadFileToCDN(connDOSpaces, localFile.Path, localFile.Type, userID)
	if err != nil {
		return utilities.CheckForError(
			c, err, 400, fmt.Sprintf("cdn upload %s object", localFile.Type), err.Error(),
		)
	}

	// Return status 201 created.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"info": fiber.Map{
			"key":        uploadFileInfo.Key,
			"etag":       uploadFileInfo.ETag,
			"size":       uploadFileInfo.Size,
			"version_id": uploadFileInfo.VersionID,
		},
		"url": fmt.Sprintf("%v/%v", os.Getenv("CDN_PUBLIC_URL"), uploadFileInfo.Key),
	})
}

// RemoveFileFromCDN func for remove exists file from CDN.
func RemoveFileFromCDN(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utilities.ExtractTokenMetaData(c)
	if err != nil {
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Get user ID from JWT.
	userID := claims.UserID.String()

	// Create new FileFromCDN struct
	fileToDelete := &models.FileFromCDN{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(fileToDelete); err != nil {
		return utilities.CheckForError(c, err, 400, "file object", err.Error())
	}

	// Get user ID from the user's upload folder on CDN.
	fileUserID, err := helpers.GetUserIDFromCDNFileKey(fileToDelete.Key)
	if err != nil {
		return utilities.CheckForError(c, err, 400, "user id", err.Error())
	}

	// Check, if user ID from JWT is equal to user's upload folder on CDN.
	if userID != fileUserID {
		return utilities.ThrowJSONError(c, 403, "file object", "you have no permissions to interact")
	}

	// Create a new connection to DO Spaces CDN.
	connDOSpaces, err := cdn.DOSpacesConnection()
	if err != nil {
		return utilities.CheckForError(c, err, 500, "cdn", err.Error())
	}

	// Remove file from CDN by key.
	if errRemoveObject := connDOSpaces.RemoveObject(
		context.Background(),
		os.Getenv("DO_SPACES_BUCKET_NAME"),
		fileToDelete.Key,
		minio.RemoveObjectOptions{
			VersionID: fileToDelete.VersionID,
		},
	); errRemoveObject != nil {
		return utilities.CheckForError(c, err, 400, "cdn remove object", err.Error())
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
