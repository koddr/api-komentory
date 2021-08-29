package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/pkg/helpers"
	"Komentory/api/platform/cdn"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Komentory/repository"
	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// GetFileListFromCDN func for return a list of files from CDN.
func GetFileListFromCDN(c *fiber.Ctx) error {
	// Get claims from JWT.
	_, errExtractTokenMetaData := utilities.ExtractTokenMetaData(c)
	if errExtractTokenMetaData != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errExtractTokenMetaData.Error(),
		})
	}

	// Create CDN connection.
	connDOSpaces, errDOSpacesConnection := cdn.DOSpacesConnection()
	if errDOSpacesConnection != nil {
		// Return status 500 and CDN connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errDOSpacesConnection.Error(),
		})
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
			// Return status 400 and bad request error.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   object.Err,
			})
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

// PutImageFileToCDN func for upload a file to CDN.
func PutImageFileToCDN(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, errExtractTokenMetaData := utilities.ExtractTokenMetaData(c)
	if errExtractTokenMetaData != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errExtractTokenMetaData.Error(),
		})
	}

	// Define user ID.
	userID := claims.UserID.String()

	// Create new FilePath struct
	filePath := &models.LocalFilePath{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(filePath); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new DO Spaces connection.
	connDOSpaces, errDOSpacesConnection := cdn.DOSpacesConnection()
	if errDOSpacesConnection != nil {
		// Return status 500 and CDN connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errDOSpacesConnection.Error(),
		})
	}

	// Upload image file process.
	uploadImageInfo, errUploadFileToCDN := cdn.UploadFileToCDN(connDOSpaces, filePath.Path, "image", userID)
	if errUploadFileToCDN != nil {
		// Return status 400 and bad request error.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errUploadFileToCDN.Error(),
		})
	}

	// Return status 201 created.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"info": fiber.Map{
			"key":        uploadImageInfo.Key,
			"etag":       uploadImageInfo.ETag,
			"size":       uploadImageInfo.Size,
			"version_id": uploadImageInfo.VersionID,
		},
		"url": fmt.Sprintf("%v/%v", os.Getenv("CDN_PUBLIC_URL"), uploadImageInfo.Key),
	})
}

// PutDocumentFileToCDN func for upload a file to CDN.
func PutDocumentFileToCDN(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, errExtractTokenMetaData := utilities.ExtractTokenMetaData(c)
	if errExtractTokenMetaData != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errExtractTokenMetaData.Error(),
		})
	}

	// Define user ID.
	userID := claims.UserID.String()

	// Create new FilePath struct
	filePath := &models.LocalFilePath{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(filePath); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new DO Spaces connection.
	connDOSpaces, errDOSpacesConnection := cdn.DOSpacesConnection()
	if errDOSpacesConnection != nil {
		// Return status 500 and CDN connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errDOSpacesConnection.Error(),
		})
	}

	// Upload image file process.
	uploadImageInfo, errUploadFileToCDN := cdn.UploadFileToCDN(connDOSpaces, filePath.Path, "document", userID)
	if errUploadFileToCDN != nil {
		// Return status 400 and bad request error.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errUploadFileToCDN.Error(),
		})
	}

	// Return status 201 created.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"info": fiber.Map{
			"key":        uploadImageInfo.Key,
			"etag":       uploadImageInfo.ETag,
			"size":       uploadImageInfo.Size,
			"version_id": uploadImageInfo.VersionID,
		},
		"url": fmt.Sprintf("%v/%v", os.Getenv("CDN_PUBLIC_URL"), uploadImageInfo.Key),
	})
}

// RemoveFileFromCDN func for remove exists file from CDN.
func RemoveFileFromCDN(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, errExtractTokenMetaData := utilities.ExtractTokenMetaData(c)
	if errExtractTokenMetaData != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errExtractTokenMetaData.Error(),
		})
	}

	// Get user ID from JWT.
	userID := claims.UserID.String()

	// Create new FileFromCDN struct
	fileToDelete := &models.FileFromCDN{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(fileToDelete); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get user ID from the user's upload folder on CDN.
	fileUserID, errGetUserIDFromCDNFileKey := helpers.GetUserIDFromCDNFileKey(fileToDelete.Key)
	if errGetUserIDFromCDNFileKey != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errGetUserIDFromCDNFileKey.Error(),
		})
	}

	// Check, if user ID from JWT is equal to user's upload folder on CDN.
	if userID != fileUserID {
		// Return status 403 and error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   repository.GenerateErrorMessage(403, "user", "it's not your file"),
		})
	}

	// Create a new connection to DO Spaces CDN.
	connDOSpaces, errDOSpacesConnection := cdn.DOSpacesConnection()
	if errDOSpacesConnection != nil {
		// Return status 500 and CDN connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errDOSpacesConnection.Error(),
		})
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
		// Return status 400 and bad request error.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errRemoveObject.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
