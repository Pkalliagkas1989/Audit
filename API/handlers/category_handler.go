package handlers

import (
	"net/http"

	"forum/repository"
	"forum/utils"
)

// CategoryHandler handles category related requests
type CategoryHandler struct {
	CategoryRepo *repository.CategoryRepository
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{CategoryRepo: repo}
}

// GetCategories returns all categories as JSON
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := h.CategoryRepo.GetAll()
	if err != nil {
		utils.ErrorResponse(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, categories, http.StatusOK)
}
