package handlers

import (
	"net/http"
	"strconv"

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

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.ErrorResponse(w, "Missing category ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.CategoryRepo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, "Failed to load category", http.StatusInternalServerError)
		return
	}

	if category == nil {
		utils.ErrorResponse(w, "Category not found", http.StatusNotFound)
		return
	}

	utils.JSONResponse(w, category, http.StatusOK)
}

