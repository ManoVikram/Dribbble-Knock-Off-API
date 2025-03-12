package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/ManoVikram/flexibbble-api/database"
	"github.com/ManoVikram/flexibbble-api/models"
	"github.com/gin-gonic/gin"
)

func GetAllProjectsHandler(c *gin.Context) {
	category := c.Query("category")

	// Extract pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Query to get the total number of projects
	var totalProjects int
	countQuery := "SELECT COUNT(*) FROM projects"
	if category != "" {
		countQuery += " WHERE category = $1"
		err = database.DB.QueryRow(countQuery, category).Scan(&totalProjects)
	} else {
		err = database.DB.QueryRow(countQuery).Scan(&totalProjects)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := `
		SELECT id, title, description, image, live_site_url, github_url, category, created_by, created_at, updated_at
		FROM projects
	`

	// If category is provided, add a WHERE clause
	var rows *sql.Rows
	if category != "" {
		query += " WHERE category = $1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3;"
		rows, err = database.DB.Query(query, category, limit, offset)
	} else {
		query += " ORDER BY updated_at DESC LIMIT $1 OFFSET $2;"
		rows, err = database.DB.Query(query, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	enrichedProjects := make([]gin.H, 0)

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

	c.JSON(http.StatusOK, gin.H{
		"data":         enrichedProjects,
		"prev_page":    page - 1,
		"current_page": page,
		"next_page":    page + 1,
		"count":        len(enrichedProjects),
		"total":        totalProjects,
	})
}
