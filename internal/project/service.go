package project

import (
	"errors"

	"github.com/shubham071122/collab/internal/user"
)

type Service struct {
	projectRepo *Repository
	userRepo    *user.Repository
}

func NewService(projectRepo *Repository, userRepo *user.Repository) *Service {
	return &Service{
		projectRepo: projectRepo,
		userRepo:    userRepo,
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

func (s *Service) UpdateProject(projectID string, userID string, req UpdateProjectRequest) (*Project, error) {
	if projectID == "" {
		return nil, errors.New("Project ID is required")
	}

	if userID == "" {
		return nil, errors.New("User ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)

	if err != nil {
		return nil, errors.New("Project not found")
	}

	// isOwner, err := s.projectRepo.IsProjectOwner(existingProject.ID, userID)
	// if err != nil {
	// 	return nil, errors.New("Failed to check project ownership")
	// }

	// if !isOwner {
	// 	return nil, errors.New("Unauthorized: only project owner can update")
	// }

	existingProject.Name = req.Name
	existingProject.Description = req.Description
	existingProject.Canvas = req.Canvas
	updatedProject, err := s.projectRepo.UpdateProject(existingProject)
	if err != nil {
		return nil, errors.New("Failed to update project")
	}
	return updatedProject, nil
}

func (s *Service) DeleteProject(projectID string, userID string) error {
	if projectID == "" {
		return errors.New("Project ID is required")
	}

	if userID == "" {
		return errors.New("User ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		return errors.New("Project not found")
	}

	isOwner, err := s.projectRepo.IsProjectOwner(existingProject.ID, userID)
	if err != nil {
		return errors.New("Failed to check project ownership")
	}

	if !isOwner {
		return errors.New("Unauthorized: only project owner can delete")
	}

	return s.projectRepo.DeleteProject(existingProject.ID)
}

func (s *Service) ShareProject(projectID string, req ShareProjectRequest, ownerID string) error {
	if projectID == "" {
		return errors.New("Project ID is required")
	}

	if ownerID == "" {
		return errors.New("User ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)

	if err != nil {
		return errors.New("Project not found")
	}

	targetUser, err := s.userRepo.GetUserByEmail(req.Email)

	if err != nil {
		return errors.New("User not found")
	}

	isOwner, err := s.projectRepo.IsProjectOwner(existingProject.ID, ownerID)
	if err != nil {
		return errors.New("Failed to check project ownership")
	}

	if !isOwner {
		return errors.New("Unauthorized: only project owner can share")
	}

	if targetUser.ID == ownerID {
		return errors.New("Cannot share project with yourself")
	}

	return s.projectRepo.ShareProject(existingProject.ID, targetUser.ID, req.Permission)
}

func (s *Service) GetCollaborators(projectID string, userID string) ([]Collaborator, error) {
	if projectID == "" {
		return nil, errors.New("Project ID is required")
	}

	if userID == "" {
		return nil, errors.New("User ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)

	if err != nil {
		return nil, errors.New("Project not found")
	}

	isOwner, err := s.projectRepo.IsProjectOwner(existingProject.ID, userID)
	if err != nil {
		return nil, errors.New("Failed to check project ownership")
	}

	if !isOwner {
		return nil, errors.New("Unauthorized: only project owner can view collaborators")
	}

	return s.projectRepo.GetCollaborators(existingProject.ID)
}

func (s *Service) UpdateCollaboratorPermission(projectID string, userID string, ownerID string, permission string) error {
	if projectID == "" {
		return errors.New("Project ID is required")
	}

	if userID == "" {
		return errors.New("User ID is required")
	}

	if ownerID == "" {
		return errors.New("Owner ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)

	if err != nil {
		return errors.New("Project not found")
	}

	isOwner, err := s.projectRepo.IsProjectOwner(existingProject.ID, ownerID)
	if err != nil {
		return errors.New("Failed to check project ownership")
	}

	if !isOwner {
		return errors.New("Unauthorized: only project owner can update collaborator permission")
	}

	return s.projectRepo.UpdateCollaboratorPermission(existingProject.ID, userID, permission)
}

func (s *Service) RemoveCollaborator(projectID string, userID string, ownerID string) error {
	if projectID == "" {
		return errors.New("Project ID is required")
	}

	if ownerID == "" {
		return errors.New("User ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)

	if err != nil {
		return errors.New("Project not found")
	}

	isOwner, err := s.projectRepo.IsProjectOwner(existingProject.ID, ownerID)
	if err != nil {
		return errors.New("Failed to check project ownership")
	}

	if !isOwner {
		return errors.New("Unauthorized: only project owner can remove collaborators")
	}

	return s.projectRepo.RemoveCollaborator(existingProject.ID, userID)
}
