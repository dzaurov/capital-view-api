package handlers

import (
	"capital-view-api/db"     // Adjust import path if needed
	"capital-view-api/models" // Adjust import path if needed
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- Member Handlers ---

// CreateMember godoc
// @Summary Create a new member entry
// @Description Add a new member record to the database
// @Tags members
// @Accept json
// @Produce json
// @Param member body models.Member true "Member data to create"
// @Success 201 {object} models.Member
// @Failure 400 {object} HTTPError "Bad Request - Invalid input data"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /members [post]
func CreateMember(c *gin.Context) {
	var input models.Member
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	result := db.DB.Create(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GetMembers godoc
// @Summary Get all member entries
// @Description Retrieve a list of all member records
// @Tags members
// @Produce json
// @Success 200 {array} models.Member
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /members [get]
func GetMembers(c *gin.Context) {
	var members []models.Member
	result := db.DB.Find(&members)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}
	c.JSON(http.StatusOK, members)
}

// GetMember godoc
// @Summary Get a single member entry by ID
// @Description Retrieve details of a specific member record using its ID
// @Tags members
// @Produce json
// @Param id path int true "Member ID"
// @Success 200 {object} models.Member
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Member not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /members/{id} [get]
func GetMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var member models.Member
	result := db.DB.First(&member, uint(id))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("member not found")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		}
		return
	}

	c.JSON(http.StatusOK, member)
}

// UpdateMember godoc
// @Summary Update an existing member entry
// @Description Modify the details of an existing member record by ID
// @Tags members
// @Accept json
// @Produce json
// @Param id path int true "Member ID"
// @Param member body models.Member true "Member data to update"
// @Success 200 {object} models.Member
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format or input data"
// @Failure 404 {object} HTTPError "Not Found - Member not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /members/{id} [put]
func UpdateMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var existingMember models.Member
	if err := db.DB.First(&existingMember, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("member not found to update")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		}
		return
	}

	var input models.Member
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	result := db.DB.Model(&existingMember).Updates(input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusOK, existingMember)
}

// DeleteMember godoc
// @Summary Delete a member entry by ID
// @Description Remove a member record from the database using its ID
// @Tags members
// @Produce json
// @Param id path int true "Member ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Member not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /members/{id} [delete]
func DeleteMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	result := db.DB.Delete(&models.Member{}, uint(id))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, NewHTTPError(errors.New("member not found to delete")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}
