package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/flexibbble-api/database"
	"github.com/ManoVikram/flexibbble-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetProjectsByUserHandler(c *gin.Context) {
	userIDString := c.Param("id")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User

	query := `
		SELECT id, name, email, "emailVerified", image, description, github_url, linkedin_url
		FROM users
		WHERE id = $1
	`

	var name, email, emailVerified, image, description, githubURL, linkedInURL sql.NullString
	err = database.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&name,
		&email,
		&emailVerified,
		&image,
		&description,
		&githubURL,
		&linkedInURL,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert sql.NullString to regular strings
	user.Name = name.String
	user.Email = email.String
	user.EmailVerified = emailVerified.String
	user.Image = image.String
	user.Description = description.String
	user.GitHubURL = githubURL.String
	user.LinkedInURL = linkedInURL.String

	var projects []models.Project

	query = `
		SELECT id, title, description, image, live_site_url, github_url, category, created_by, created_at, updated_at
		FROM projects
		WHERE created_by = $1;
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var project models.Project
		if err := rows.Scan(
			&project.ID,
			&project.Title,
			&project.Description,
			&project.Image,
			&project.LiveSiteURL,
			&project.GitHubURL,
			&project.Category,
			&project.CreatedBy,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning project data"})
			return
		}

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             user.ID,
		"name":           user.Name,
		"email":          user.Email,
		"email_verified": user.EmailVerified,
		"image":          user.Image,
		"description":    user.Description,
		"github_url":     user.GitHubURL,
		"linkedin_url":   user.LinkedInURL,
		"projects":       projects,
	})
}
