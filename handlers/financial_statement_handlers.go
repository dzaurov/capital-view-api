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

// --- FinancialStatement Handlers ---

// CreateFinancialStatement godoc
// @Summary Create a new financial statement entry
// @Description Add a new financial statement record to the database
// @Tags financial-statements
// @Accept json
// @Produce json
// @Param financialStatement body models.FinancialStatement true "Financial Statement data to create"
// @Success 201 {object} models.FinancialStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid input data"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /financial-statements [post]
func CreateFinancialStatement(c *gin.Context) {
	var input models.FinancialStatement
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

// GetFinancialStatements godoc
// @Summary Get all financial statement entries
// @Description Retrieve a list of all financial statement records
// @Tags financial-statements
// @Produce json
// @Success 200 {array} models.FinancialStatement
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /financial-statements [get]
func GetFinancialStatements(c *gin.Context) {
	var financialStatements []models.FinancialStatement
	result := db.DB.Find(&financialStatements)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}
	c.JSON(http.StatusOK, financialStatements)
}

// GetFinancialStatement godoc
// @Summary Get a single financial statement entry by ID
// @Description Retrieve details of a specific financial statement record using its ID
// @Tags financial-statements
// @Produce json
// @Param id path int true "Financial Statement ID"
// @Success 200 {object} models.FinancialStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Financial Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /financial-statements/{id} [get]
func GetFinancialStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var financialStatement models.FinancialStatement
	result := db.DB.First(&financialStatement, uint(id))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("financial statement not found")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		}
		return
	}

	c.JSON(http.StatusOK, financialStatement)
}

// UpdateFinancialStatement godoc
// @Summary Update an existing financial statement entry
// @Description Modify the details of an existing financial statement record by ID
// @Tags financial-statements
// @Accept json
// @Produce json
// @Param id path int true "Financial Statement ID"
// @Param financialStatement body models.FinancialStatement true "Financial Statement data to update"
// @Success 200 {object} models.FinancialStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format or input data"
// @Failure 404 {object} HTTPError "Not Found - Financial Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /financial-statements/{id} [put]
func UpdateFinancialStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var existingFinancialStatement models.FinancialStatement
	if err := db.DB.First(&existingFinancialStatement, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("financial statement not found to update")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		}
		return
	}

	var input models.FinancialStatement
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	result := db.DB.Model(&existingFinancialStatement).Updates(input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusOK, existingFinancialStatement)
}

// DeleteFinancialStatement godoc
// @Summary Delete a financial statement entry by ID
// @Description Remove a financial statement record from the database using its ID
// @Tags financial-statements
// @Produce json
// @Param id path int true "Financial Statement ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Financial Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /financial-statements/{id} [delete]
func DeleteFinancialStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	result := db.DB.Delete(&models.FinancialStatement{}, uint(id))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, NewHTTPError(errors.New("financial statement not found to delete")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Financial Statement deleted successfully"})
}
