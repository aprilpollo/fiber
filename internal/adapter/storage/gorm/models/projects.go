package models

import (
	"time"
)

type Project struct {
	BaseModel

	OrganizationID *uint  `json:"organization_id" gorm:"index"`
	Name           string `json:"name" gorm:"not null;size:255"`
	Description    string `json:"description" gorm:"type:text"`
	Key            string `json:"key" gorm:"not null;unique;size:10"`
	OwnerID        uint   `json:"owner_id" gorm:"not null;index"`
	StatusID       uint   `json:"status_id" gorm:"not null;index;default:1"` // FK to project_statuses (1=active)

	// Relationships
	Status       ProjectStatus   `json:"status" gorm:"foreignKey:StatusID"`
	Organization *Organization   `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	Owner        *User           `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Members      []ProjectMember `json:"members"`
}

type ProjectMember struct {
	BaseModel

	ProjectID uint       `json:"project_id" gorm:"not null;index;uniqueIndex:idx_project_user"`
	UserID    uint       `json:"user_id" gorm:"not null;index;uniqueIndex:idx_project_user"`
	RoleID    uint       `json:"role_id" gorm:"not null;index"`             // FK to project_member_roles
	StatusID  uint       `json:"status_id" gorm:"not null;index;default:1"` // FK to member_statuses
	InvitedAt *time.Time `json:"invited_at"`
	JoinedAt  *time.Time `json:"joined_at"`
	InvitedBy *uint      `json:"invited_by" gorm:"index"`

	// Relationships
	Role    ProjectMemberRole `json:"role" gorm:"foreignKey:RoleID"`
	Status  MemberStatus      `json:"status" gorm:"foreignKey:StatusID"` 
	Project *Project          `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	User    *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Inviter *User             `json:"inviter,omitempty" gorm:"foreignKey:InvitedBy"`
}

type ProjectStatus struct {
	ID          uint   `json:"id" gorm:"uniqueIndex;not null"`
	Name        string `json:"name" gorm:"not null;unique;size:50"`
	Description string `json:"description" gorm:"type:text"`

	Projects []Project `json:"projects,omitempty" gorm:"foreignKey:StatusID"`
}

type ProjectMemberRole struct {
	ID          uint   `json:"id" gorm:"uniqueIndex;not null"`
	Name        string `json:"name" gorm:"not null;unique;size:50"`
	Description string `json:"description" gorm:"type:text"`
	IsDefault   bool   `json:"is_default" gorm:"not null;default:false"`

	// Permissions
	CanManageProject    bool `json:"can_manage_project" gorm:"not null;default:false"`
	CanManageMembers    bool `json:"can_manage_members" gorm:"not null;default:false"`
	CanCreateTasks      bool `json:"can_create_tasks" gorm:"not null;default:false"`
	CanManageTasks      bool `json:"can_manage_tasks" gorm:"not null;default:false"`
	CanDeleteTasks      bool `json:"can_delete_tasks" gorm:"not null;default:false"`
	CanViewAllTasks     bool `json:"can_view_all_tasks" gorm:"not null;default:false"`
	CanManageComponents bool `json:"can_manage_components" gorm:"not null;default:false"`
	CanManageVersions   bool `json:"can_manage_versions" gorm:"not null;default:false"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	ProjectMembers []ProjectMember `json:"project_members,omitempty" gorm:"foreignKey:RoleID"`
}
