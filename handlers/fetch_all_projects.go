package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/flexibbble-api/database"
	"github.com/ManoVikram/flexibbble-api/models"
	"github.com/gin-gonic/gin"
)

func FetchAllProjectsHandler(c *gin.Context) {
	query := `
		SELECT id, title, description, image, live_site_url, github_url, category, created_by, created_at, updated_at
		FROM projects
		ORDER BY updated_at DESC;
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var enrichedProjects []gin.H

	for rows.Next() {
		var project models.Project

		err := rows.Scan(
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
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		query = "SELECT id, name, email, image, description, github_url, linkedin_url FROM users WHERE id = $1"
		// NOTE: If any of the selected fields contain NULL, .Scan() will fail
		// Handle nullable fields
		var name, email, image, description, githubURL, linkedInURL sql.NullString
		err = database.DB.QueryRow(query, project.CreatedBy).Scan(
			&user.ID,
			&name,
			&email,
			&image,
			&description,
			&githubURL,
			&linkedInURL,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching details"})
				return
			}

			user = models.User{}
		}

		// Convert sql.NullString to regular strings
		user.Name = name.String
		user.Email = email.String
		user.Image = image.String
		user.Description = description.String
		user.GitHubURL = githubURL.String
		user.LinkedInURL = linkedInURL.String

		enrichedProjects = append(enrichedProjects, gin.H{
			"id":            project.ID,
			"title":         project.Title,
			"description":   project.Description,
			"image":         project.Image,
			"live_site_url": project.LiveSiteURL,
			"github_url":    project.GitHubURL,
			"category":      project.Category,
			"created_at":    project.CreatedAt,
			"updated_at":    project.UpdatedAt,
			"created_by":    user, // Embed user details instead of just UUID
		})
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, enrichedProjects)
}
