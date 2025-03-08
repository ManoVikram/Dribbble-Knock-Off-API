package handlers

import (
	"net/http"

	"github.com/ManoVikram/flexibbble-api/database"
	"github.com/ManoVikram/flexibbble-api/models"
	"github.com/gin-gonic/gin"
)

func CreateProjectHandler(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageURL, err := UploadToImageKitHandler("image", project.Image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	project.Image = imageURL

	query := `INSERT INTO projects (title, description, image, live_site_url, github_url, category, created_by)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = database.DB.Exec(query, project.Title, project.Description, project.Image, project.LiveSiteURL, project.GitHubURL, project.Category, project.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}
