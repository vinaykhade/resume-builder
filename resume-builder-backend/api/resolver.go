// api/resolver.go

package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/vinaykhade/resume-builder-backend/internal/models"
	"github.com/vinaykhade/resume-builder-backend/internal/services"
)

type Resolver struct {
	resumeService *services.ResumeService
}

func NewResolver(resumeService *services.ResumeService) *Resolver {
	return &Resolver{
		resumeService: resumeService,
	}
}

func (r *Resolver) ListResumes(ctx context.Context, args struct{ UserId int32 }) ([]*models.Resume, error) {
	resumes, err := r.resumeService.ListResumes(ctx, args.UserId)
	if err != nil {
		log.Printf("Failed to list resumes: %v", err)
		return nil, err
	}

	var results []*models.Resume
	for _, item := range resumes.Items {
		createdAt, _ := time.Parse(time.RFC3339, item.CreatedAt.AsTime().Format(time.RFC3339))
		updatedAt, _ := time.Parse(time.RFC3339, item.UpdatedAt.AsTime().Format(time.RFC3339))
		resume := &models.Resume{
			ID:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			FileURL:     item.FileUrl,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		results = append(results, resume)
	}

	return results, nil
}

func (r *Resolver) UploadResume(ctx context.Context, args struct {
	Name, Description string
	File              multipart.File
}) (*models.Resume, error) {
	fileBytes, err := ioutil.ReadAll(args.File)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return nil, err
	}

	fileName := filepath.Base(args.File.(interface {
		Name() string
	}).Name())
	filePath := fmt.Sprintf("uploads/%s", fileName)
	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return nil, err
	}
	defer f.Close()

	_, err = f.Write(fileBytes)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		return nil, err
	}

	res, err := r.resumeService.UploadResume(ctx, &services.ResumeRequest{
		Name:        args.Name,
		Description: args.Description,
		File:        fileBytes,
	})
	if err != nil {
		log.Printf("Failed to upload resume: %v", err)
		return nil, err
	}

	return &models.Resume{
		ID:          res.Id,
		Name:        args.Name,
		Description: args.Description,
		FileURL:     filePath,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
