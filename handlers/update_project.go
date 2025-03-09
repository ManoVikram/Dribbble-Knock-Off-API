package handlers

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"

	"github.com/ManoVikram/flexibbble-api/database"
	"github.com/ManoVikram/flexibbble-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func isBase64ImageURL(str string) bool {
	// Regular expression to check if the string is a valid base64 image data URL
	base64ImageRegex := regexp.MustCompile(`^data:image/(png|jpeg|jpg|gif|webp|bmp|svg\+xml);base64,[A-Za-z0-9+/=]+$`)

	if !base64ImageRegex.MatchString(str) {
		return false
	}

	// Extract the base64 part
	parts := strings.Split(str, ",")
	if len(parts) < 2 {
		return false
	}

	// Validate the base64 data
	_, err := base64.StdEncoding.DecodeString(parts[1])
	return err == nil
}

func UpdateProjectHandler(c *gin.Context) {
	projectIDString := c.Param("id")
	projectID, err := uuid.Parse(projectIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var updateProject models.Project
	if err := c.ShouldBindJSON(&updateProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if isBase64ImageURL(updateProject.Image) {
		imageURL, err := UploadToImageKitHandler("image", updateProject.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updateProject.Image = imageURL
	}

	updateProject.ID = projectID

	query := `
		UPDATE projects
		SET title = $1, description = $2, image = $3, live_site_url = $4, github_url = $5, category = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	_, err = database.DB.Exec(query,
		updateProject.Title,
		updateProject.Description,
		updateProject.Image,
		updateProject.LiveSiteURL,
		updateProject.GitHubURL,
		updateProject.Category,
		updateProject.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}
