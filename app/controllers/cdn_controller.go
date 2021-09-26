package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/pkg/helpers"
	"Komentory/api/platform/cdn"
	"context"
	"fmt"
	"os"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// PutFileToCDN func for upload a file to CDN.
// Allowed types: image, document.
func PutFileToCDN(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utilities.TokenValidateExpireTime(c)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
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
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "cdn", err.Error())
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
		"status": fiber.StatusCreated,
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
	claims, err := utilities.TokenValidateExpireTime(c)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Get user ID from JWT.
	userID := claims.UserID.String()

	// Create new FileFromCDN struct
	fileToDelete := &models.FileFromCDN{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(fileToDelete); err != nil {
		return utilities.CheckForError(c, err, 400, "file object", err.Error())
	}

	// Get owner's user ID for file from the upload folder on CDN.
	fileOwnerUserID, err := helpers.GetUserIDFromCDNFileKey(fileToDelete.Key)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 400, "user id", err.Error())
	}

	// Check, if user ID from JWT is equal to getted owner user ID from upload folder on CDN.
	if userID != fileOwnerUserID {
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "file object", "you have no permissions")
	}

	// Create a new connection to DO Spaces CDN.
	connDOSpaces, err := cdn.DOSpacesConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "cdn", err.Error())
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
		return utilities.CheckForErrorWithStatusCode(c, err, 400, "cdn remove object", err.Error())
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
