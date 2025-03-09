package handlers

import (
	"net/http"

	"github.com/ManoVikram/flexibbble-api/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeleteProjectHandler(c *gin.Context) {
	projectIDString := c.Param("id")
	projectID, err := uuid.Parse(projectIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	query := `
		DELETE FROM projects WHERE id = $1;
	`
	result, err := database.DB.Exec(query, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify deletion"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
