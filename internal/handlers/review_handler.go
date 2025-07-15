package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"reviews/internal/models"
	"reviews/internal/services"
)

type ReviewHandler struct {
	Service *services.ReviewService
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	industry := r.FormValue("industry")
	service := r.FormValue("service")
	ratingStr := r.FormValue("rating")
	rating, _ := strconv.Atoi(ratingStr)

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "photo is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	pdfPath := ""
	pdfFile, pdfHeader, errPdf := r.FormFile("pdf")
	if errPdf == nil {
		defer pdfFile.Close()
	}

	saveDir := filepath.Join("uploads", "reviews")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		http.Error(w, "unable to create image directory", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("review_%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
	outPath := filepath.Join(saveDir, filename)
	dst, err := os.Create(outPath)
	if err != nil {
		http.Error(w, "unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, "unable to save file", http.StatusInternalServerError)
		return
	}

	if errPdf == nil {
		pdfFilename := fmt.Sprintf("review_%d%s", time.Now().UnixNano(), filepath.Ext(pdfHeader.Filename))
		pdfOut := filepath.Join(saveDir, pdfFilename)
		pdfDst, err := os.Create(pdfOut)
		if err != nil {
			http.Error(w, "unable to save pdf", http.StatusInternalServerError)
			return
		}
		defer pdfDst.Close()
		if _, err := io.Copy(pdfDst, pdfFile); err != nil {
			http.Error(w, "unable to save pdf", http.StatusInternalServerError)
			return
		}
		pdfPath = fmt.Sprintf("/pdfs/reviews/%s", pdfFilename)
	}

	rev := models.Reviews{
		Name: name,

		Photo: fmt.Sprintf("/images/reviews/%s", filename),

		PdfFile:  pdfPath,
		Industry: industry,
		Service:  service,

		Description: description,
		Rating:      rating,
	}

	created, err := h.Service.Create(r.Context(), rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

func (h *ReviewHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.Service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	rev, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rev)
}

func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	name := r.FormValue("name")
	description := r.FormValue("description")
	industry := r.FormValue("industry")
	service := r.FormValue("service")
	ratingStr := r.FormValue("rating")
	rating, _ := strconv.Atoi(ratingStr)

	existing, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	photoPath := existing.Photo
	pdfPath := existing.PdfFile
	file, header, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()
		saveDir := filepath.Join("uploads", "reviews")
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			http.Error(w, "unable to create image directory", http.StatusInternalServerError)
			return
		}
		filename := fmt.Sprintf("review_%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
		outPath := filepath.Join(saveDir, filename)
		dst, err := os.Create(outPath)
		if err != nil {
			http.Error(w, "unable to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "unable to save file", http.StatusInternalServerError)
			return
		}
		photoPath = fmt.Sprintf("/images/reviews/%s", filename)
	}

	pdfFile, pdfHeader, errPdf := r.FormFile("pdf")
	if errPdf == nil {
		defer pdfFile.Close()
		saveDir := filepath.Join("uploads", "reviews")
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			http.Error(w, "unable to create image directory", http.StatusInternalServerError)
			return
		}
		pdfFilename := fmt.Sprintf("review_%d%s", time.Now().UnixNano(), filepath.Ext(pdfHeader.Filename))
		pdfOut := filepath.Join(saveDir, pdfFilename)
		dstPdf, err := os.Create(pdfOut)
		if err != nil {
			http.Error(w, "unable to save pdf", http.StatusInternalServerError)
			return
		}
		defer dstPdf.Close()
		if _, err := io.Copy(dstPdf, pdfFile); err != nil {
			http.Error(w, "unable to save pdf", http.StatusInternalServerError)
			return
		}
		pdfPath = fmt.Sprintf("/pdfs/reviews/%s", pdfFilename)
	}

	rev := models.Reviews{
		ID:          id,
		Name:        name,
		Photo:       photoPath,
		PdfFile:     pdfPath,
		Industry:    industry,
		Service:     service,
		Description: description,
		Rating:      rating,
	}

	updated, err := h.Service.Update(r.Context(), rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if err := h.Service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ServeReviewImage handles GET /images/reviews/:filename requests and serves saved review images.
func (h *ReviewHandler) ServeReviewImage(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get(":filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	imagePath := filepath.Join("uploads", "reviews", filename)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	switch filepath.Ext(imagePath) {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	http.ServeFile(w, r, imagePath)
}

// ServeReviewPDF handles GET /pdfs/reviews/:filename requests and serves saved PDF files.
func (h *ReviewHandler) ServeReviewPDF(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get(":filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	pdfPath := filepath.Join("uploads", "reviews", filename)
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		http.Error(w, "pdf not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, pdfPath)
}
