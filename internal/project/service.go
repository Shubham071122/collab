package project

import (
	"errors"
)

type Service struct {
	projectRepo *Repository
}

func NewService(projectRepo *Repository) *Service {
	return &Service{
		projectRepo: projectRepo,
	}
}

func (s *Service) GetProjectByID(projectID string) (*Project, error) {
	project, err := s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, errors.New("Project not found")
	}
	return project, nil
}

func (s *Service) CreateProject(req CreateProjectRequest) (*Project, error) {

	if req.Name == "" {
		return nil, errors.New("Project name is required")
	}

	if req.UserID == "" {
		return nil, errors.New("User ID is required")
	}

	project, err := s.projectRepo.CreateProject(req.Name, req.Description, req.UserID)
	if err != nil {
		return nil, errors.New("Failed to create project")
	}
	return project, nil
}

func (s *Service) UpdateProject(projectID string, req UpdateProjectRequest) (*Project, error) {
	if projectID == "" {
		return nil, errors.New("Project ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)

	if err != nil {
		return nil, errors.New("Project not found")
	}

	existingProject.Name = req.Name
	existingProject.Description = req.Description
	updatedProject, err := s.projectRepo.UpdateProject(existingProject)
	if err != nil {
		return nil, errors.New("Failed to update project")
	}
	return updatedProject, nil
}

func (s *Service) DeleteProject(projectID string) error {
	if projectID == "" {
		return errors.New("Project ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)

	if err != nil {
		return errors.New("Project not found")
	}
	return s.projectRepo.DeleteProject(existingProject.ID)
}
