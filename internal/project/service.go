package project

import (
	"errors"
	"strings"

	"github.com/shubham071122/collab/internal/collaboration"
	"github.com/shubham071122/collab/internal/subscription"
	"github.com/shubham071122/collab/internal/user"
)

type Service struct {
	projectRepo *Repository
	userRepo    *user.Repository
	hub         *collaboration.Hub
	subService  *subscription.Service
}

func NewService(
	projectRepo *Repository,
	userRepo *user.Repository,
	hub *collaboration.Hub,
	subService *subscription.Service,
) *Service {
	return &Service{
		projectRepo: projectRepo,
		userRepo:    userRepo,
		hub:         hub,
		subService:  subService,
	}
}

func (s *Service) GetProjectByID(projectID string, userID string) (*Project, error) {
	if projectID == "" {
		return nil, errors.New("Project ID is required")
	}

	if userID == "" {
		return nil, errors.New("User ID is required")
	}

	project, err := s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, errors.New("Project not found")
	}

	isOwner, err := s.projectRepo.IsProjectOwner(project.ID, userID)
	if err != nil {
		return nil, errors.New("Failed to check project ownership")
	}
	if isOwner {
		return project, nil
	}

	isCollab, err := s.projectRepo.IsProjectCollaborator(project.ID, userID)
	if err != nil {
		return nil, errors.New("Failed to check project collaborator")
	}
	if isCollab {
		return project, nil
	}

	return nil, errors.New("Unauthorized: access denied")
}

func (s *Service) CreateProject(req CreateProjectRequest) (*Project, error) {

	if req.Name == "" {
		return nil, errors.New("Project name is required")
	}

	if req.UserID == "" {
		return nil, errors.New("User ID is required")
	}

	ownedCount, err := s.projectRepo.GetOwnedProjectsCount(req.UserID)
	if err != nil {
		return nil, errors.New("Failed to check project limits")
	}

	canCreate, err := s.subService.CheckProjectLimit(req.UserID, ownedCount)
	if err != nil {
		return nil, err
	}
	if !canCreate {
		return nil, subscription.ErrLimitExceeded
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

	ownerOwnedCount, err := s.projectRepo.GetOwnedProjectsCount(existingProject.OwnerID)
	if err != nil {
		return nil, errors.New("Failed to verify project limits")
	}

	canEdit, err := s.subService.CheckProjectLimit(existingProject.OwnerID, ownerOwnedCount-1)
	if err != nil {
		return nil, err
	}
	if !canEdit {
		return nil, subscription.ErrLimitExceeded
	}

	isOwner, err := s.projectRepo.IsProjectOwner(existingProject.ID, userID)
	if err != nil {
		return nil, errors.New("Failed to check project ownership")
	}
	if !isOwner {
		isCollab, err := s.projectRepo.IsProjectCollaborator(existingProject.ID, userID)
		if err != nil {
			return nil, errors.New("Failed to check collaborator status")
		}
		if !isCollab {
			return nil, errors.New("Unauthorized: only owner or collaborators with edit permission can update")
		}

		permission, err := s.projectRepo.GetCollaboratorPermission(existingProject.ID, userID)
		if err != nil {
			return nil, errors.New("Failed to get collaborator permission")
		}

		if permission != "edit" {
			return nil, errors.New("Unauthorized: insufficient permission to update project")
		}
	}

	if req.Name != nil {
		existingProject.Name = *req.Name
	}
	if req.Description != nil {
		existingProject.Description = *req.Description
	}
	if req.Canvas != nil {
		existingProject.Canvas = *req.Canvas
	}
	if req.IsArchived != nil {
		if !isOwner {
			return nil, errors.New("Unauthorized: only the owner can archive or unarchive the project")
		}
		existingProject.IsArchived = *req.IsArchived
	}

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

	emailStr := strings.ToLower(strings.TrimSpace(req.Email))
	targetUser, err := s.userRepo.GetUserByEmail(emailStr)

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

	collaborators, err := s.projectRepo.GetCollaborators(existingProject.ID)
	if err != nil {
		return errors.New("Failed to verify collaborator count")
	}

	isAlreadyCollaborator := false
	for _, c := range collaborators {
		if c.ID == targetUser.ID {
			isAlreadyCollaborator = true
			break
		}
	}

	if !isAlreadyCollaborator {
		canShare, err := s.subService.CheckShareLimit(ownerID, len(collaborators))
		if err != nil {
			return err
		}
		if !canShare {
			return subscription.ErrLimitExceeded
		}
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

func (s *Service) GetProjectMembers(projectID string, userID string) (*MemberInfo, error) {
	if projectID == "" {
		return nil, errors.New("Project ID is required")
	}
	if userID == "" {
		return nil, errors.New("User ID is required")
	}
	return s.projectRepo.GetProjectMembers(projectID, userID)
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

	err = s.projectRepo.UpdateCollaboratorPermission(existingProject.ID, userID, permission)
	if err != nil {
		return err
	}

	s.hub.UpdateUserPermission(existingProject.ID, userID, permission)
	return nil
}

func (s *Service) RemoveCollaborator(projectID string, userID string, callerID string) error {
	if projectID == "" {
		return errors.New("Project ID is required")
	}

	if callerID == "" {
		return errors.New("User ID is required")
	}

	existingProject, err := s.projectRepo.GetProjectByID(projectID)

	if err != nil {
		return errors.New("Project not found")
	}

	isOwner, err := s.projectRepo.IsProjectOwner(existingProject.ID, callerID)
	if err != nil {
		return errors.New("Failed to check project ownership")
	}

	if !isOwner && callerID != userID {
		return errors.New("Unauthorized: only the project owner or the collaborator themselves can remove this collaboration")
	}

	err = s.projectRepo.RemoveCollaborator(existingProject.ID, userID)
	if err != nil {
		return err
	}

	s.hub.UpdateUserPermission(existingProject.ID, userID, "")
	return nil
}

func (s *Service) GetProjects(userID string) ([]Project, error) {
	if userID == "" {
		return nil, errors.New("User ID is required")
	}
	return s.projectRepo.GetProjectsByUserID(userID)
}
