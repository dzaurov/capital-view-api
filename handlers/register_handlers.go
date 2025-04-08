package handlers

import (
	"capital-view-api/db"     // Adjust import path
	"capital-view-api/models" // Adjust import path
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- Helper for Error Response ---
type HTTPError struct {
	Error string `json:"error"`
}

func NewHTTPError(err error) HTTPError {
	return HTTPError{Error: err.Error()}
}

// --- Register Handlers ---

// CreateRegister godoc
// @Summary Create a new register entry
// @Description Add a new register record to the database
// @Tags registers
// @Accept json
// @Produce json
// @Param register body models.Registers true "Register data to create"
// @Success 201 {object} models.Registers
// @Failure 400 {object} HTTPError "Bad Request - Invalid input data"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /registers [post]
func CreateRegister(c *gin.Context) {
	var input models.Registers
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	// GORM automatically handles ID generation if primaryKey;autoIncrement is set
	result := db.DB.Create(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GetRegisters godoc
// @Summary Get all register entries
// @Description Retrieve a list of all register records
// @Tags registers
// @Produce json
// @Success 200 {array} models.Registers
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /registers [get]
func GetRegisters(c *gin.Context) {
	var registers []models.Registers
	result := db.DB.Find(&registers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}
	c.JSON(http.StatusOK, registers)
}

// GetRegister godoc
// @Summary Get a single register entry by ID
// @Description Retrieve details of a specific register record using its ID
// @Tags registers
// @Produce json
// @Param id path int true "Register ID"
// @Success 200 {object} models.Registers
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Register not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /registers/{id} [get]
func GetRegister(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Use ParseUint for auto-increment IDs
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var register models.Registers
	result := db.DB.First(&register, uint(id)) // GORM expects the correct type for the ID

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("register not found")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		}
		return
	}

	c.JSON(http.StatusOK, register)
}

// UpdateRegister godoc
// @Summary Update an existing register entry
// @Description Modify the details of an existing register record by ID
// @Tags registers
// @Accept json
// @Produce json
// @Param id path int true "Register ID"
// @Param register body models.Registers true "Register data to update"
// @Success 200 {object} models.Registers
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format or input data"
// @Failure 404 {object} HTTPError "Not Found - Register not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /registers/{id} [put]
func UpdateRegister(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	// Check if record exists
	var existingRegister models.Registers
	if err := db.DB.First(&existingRegister, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("register not found to update")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		}
		return
	}

	// Bind JSON data to update
	// Use a map or a dedicated update struct if you only want to allow specific fields to be updated
	var input models.Registers
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	// GORM's Updates method only updates non-zero fields by default.
	// Use .Model(&existingRegister).Updates(input) to update based on the input struct
	// Or use .Model(&existingRegister).Select("*").Updates(input) to update all fields including clearing them if null/zero in input
	// For partial updates with nulls, often a map[string]interface{} is better.
	// Let's stick to updating based on the provided input struct for simplicity here.
	result := db.DB.Model(&existingRegister).Updates(input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusOK, existingRegister) // Return the updated record
}

// DeleteRegister godoc
// @Summary Delete a register entry by ID
// @Description Remove a register record from the database using its ID
// @Tags registers
// @Produce json
// @Param id path int true "Register ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Register not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /registers/{id} [delete]
func DeleteRegister(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	// Delete the record
	result := db.DB.Delete(&models.Registers{}, uint(id))

	if result.Error != nil {
		// Handle potential errors during delete, though gorm.ErrRecordNotFound isn't typically returned on Delete
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, NewHTTPError(errors.New("register not found to delete")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Register deleted successfully"})
}
