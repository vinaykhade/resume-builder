// internal/services/resume.go

package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/vinaykhade/resume-builder-backend/api/resume"
)

type ResumeService struct {
	db *sql.DB
}

func NewResumeService(db *sql.DB) *ResumeService {
	return &ResumeService{
		db: db,
	}
}

func (s *ResumeService) UploadResume(ctx context.Context, req *resume.ResumeRequest) (*resume.ResumeResponse, error) {
	fileName := fmt.Sprintf("%s.pdf", uuid.New().String())
	filePath := fmt.Sprintf("uploads/%s", fileName)
	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return nil, err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(req.File))
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		return nil, err
	}

	query := `
		INSERT INTO resumes (user_id, name, description, file_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int32
	err = s.db.QueryRow(query, 1, req.Name, req.Description, filePath).Scan(&id)
	if err != nil {
		log.Printf("Failed to insert resume: %v", err)
		return nil, err
	}

	return &resume.ResumeResponse{Id: id}, nil
}

func (s *ResumeService) ListResumes(ctx context.Context, req *resume.ListRequest) (*resume.ListResponse, error) {
	query := `
		SELECT id, name, description, file_path, created_at, updated_at
		FROM resumes
		WHERE user_id = $1
	`
	rows, err := s.db.Query(query, req.UserId)
	if err != nil {
		log.Printf("Failed to query resumes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []*resume.ResumeItem
	for rows.Next() {
		var item resume.ResumeItem
		var createdAt, updatedAt time.Time
		err = rows.Scan(&item.Id, &item.Name, &item.Description, &item.FileUrl, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Failed to scan resume: %v", err)
			return nil, err
		}
		item.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		item.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)
		items = append(items, &item)
	}

	return &resume.ListResponse{Items: items}, nil
}
