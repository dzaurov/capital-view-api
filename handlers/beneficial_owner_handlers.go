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

// --- BeneficialOwner Handlers ---

// CreateBeneficialOwner godoc
// @Summary Create a new beneficial owner entry
// @Description Add a new beneficial owner record to the database
// @Tags beneficial-owners
// @Accept json
// @Produce json
// @Param beneficialOwner body models.BeneficialOwner true "Beneficial Owner data to create"
// @Success 201 {object} models.BeneficialOwner
// @Failure 400 {object} HTTPError "Bad Request - Invalid input data"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /beneficial-owners [post]
func CreateBeneficialOwner(c *gin.Context) {
	var input models.BeneficialOwner
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

// GetBeneficialOwners godoc
// @Summary Get all beneficial owner entries
// @Description Retrieve a list of all beneficial owner records
// @Tags beneficial-owners
// @Produce json
// @Success 200 {array} models.BeneficialOwner
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /beneficial-owners [get]
func GetBeneficialOwners(c *gin.Context) {
	var beneficialOwners []models.BeneficialOwner
	result := db.DB.Find(&beneficialOwners)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}
	c.JSON(http.StatusOK, beneficialOwners)
}

// GetBeneficialOwner godoc
// @Summary Get a single beneficial owner entry by ID
// @Description Retrieve details of a specific beneficial owner record using its ID
// @Tags beneficial-owners
// @Produce json
// @Param id path int true "Beneficial Owner ID"
// @Success 200 {object} models.BeneficialOwner
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Beneficial Owner not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /beneficial-owners/{id} [get]
func GetBeneficialOwner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var beneficialOwner models.BeneficialOwner
	result := db.DB.First(&beneficialOwner, uint(id))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("beneficial owner not found")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		}
		return
	}

	c.JSON(http.StatusOK, beneficialOwner)
}

// UpdateBeneficialOwner godoc
// @Summary Update an existing beneficial owner entry
// @Description Modify the details of an existing beneficial owner record by ID
// @Tags beneficial-owners
// @Accept json
// @Produce json
// @Param id path int true "Beneficial Owner ID"
// @Param beneficialOwner body models.BeneficialOwner true "Beneficial Owner data to update"
// @Success 200 {object} models.BeneficialOwner
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format or input data"
// @Failure 404 {object} HTTPError "Not Found - Beneficial Owner not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /beneficial-owners/{id} [put]
func UpdateBeneficialOwner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var existingBeneficialOwner models.BeneficialOwner
	if err := db.DB.First(&existingBeneficialOwner, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("beneficial owner not found to update")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		}
		return
	}

	var input models.BeneficialOwner
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	result := db.DB.Model(&existingBeneficialOwner).Updates(input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusOK, existingBeneficialOwner)
}

// DeleteBeneficialOwner godoc
// @Summary Delete a beneficial owner entry by ID
// @Description Remove a beneficial owner record from the database using its ID
// @Tags beneficial-owners
// @Produce json
// @Param id path int true "Beneficial Owner ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Beneficial Owner not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /beneficial-owners/{id} [delete]
func DeleteBeneficialOwner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	result := db.DB.Delete(&models.BeneficialOwner{}, uint(id))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, NewHTTPError(errors.New("beneficial owner not found to delete")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beneficial Owner deleted successfully"})
}
