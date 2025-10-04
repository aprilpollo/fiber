package models

import "time"

type Organization struct {
	BaseModel

	Name        string `json:"name" gorm:"not null;size:255"`
	Slug        string `json:"slug" gorm:"not null;uniqueIndex;size:100"`
	Description string `json:"description" gorm:"type:text"`
	LogoURL     string `json:"logo_url"`
	PlanType    string `json:"plan_type" gorm:"not null;default:'free';index"` // free, pro, enterprise
	StatusID    uint   `json:"status_id" gorm:"not null;index;default:1"`
	Settings    string `json:"settings" gorm:"type:json"` // JSON settings

	// Relationships
	Status   OrganizationStatus   `json:"status" gorm:"foreignKey:StatusID"`
	Members  []OrganizationMember `json:"members,omitempty" gorm:"foreignKey:OrganizationID"`
	Projects []Project            `json:"projects,omitempty" gorm:"foreignKey:OrganizationID"`
}

type OrganizationMember struct {
	BaseModel

	OrganizationID uint       `json:"organization_id" gorm:"not null;index"`
	UserID         uint       `json:"user_id" gorm:"not null;index"`
	RoleID         uint       `json:"role_id" gorm:"not null;index"`
	StatusID       uint       `json:"status_id" gorm:"not null;index;default:1"`
	InvitedAt      *time.Time `json:"invited_at"`
	JoinedAt       *time.Time `json:"joined_at"`
	InvitedBy      *uint      `json:"invited_by" gorm:"index"`

	// Relationships
	Role         OrganizationMemberRole `json:"role" gorm:"foreignKey:RoleID"`
	Status       MemberStatus           `json:"status" gorm:"foreignKey:StatusID"`
	Organization *Organization          `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	User         User                   `json:"user" gorm:"foreignKey:UserID"`
	Inviter      *User                  `json:"inviter,omitempty" gorm:"foreignKey:InvitedBy"`
}

type OrganizationStatus struct {
	ID          uint   `json:"id" gorm:"uniqueIndex;not null"`
	Name        string `json:"name" gorm:"not null;unique;size:50"`
	Description string `json:"description" gorm:"type:text"`

	Organizations []Organization `json:"organizations,omitempty" gorm:"foreignKey:StatusID"`
}

type MemberStatus struct {
	ID          uint   `json:"id" gorm:"uniqueIndex;not null"`
	Name        string `json:"name" gorm:"not null;unique;size:50"`
	Description string `json:"description" gorm:"type:text"`

	OrganizationMembers []OrganizationMember `json:"organization_members,omitempty" gorm:"foreignKey:StatusID"`
}

type OrganizationMemberRole struct {
	ID          uint   `json:"id" gorm:"uniqueIndex;not null"`
	Name        string `json:"name" gorm:"not null;unique;size:50"`
	Description string `json:"description" gorm:"type:text"`
	IsDefault   bool   `json:"is_default" gorm:"not null;default:false"`
	IsPreview   bool   `json:"is_preview" gorm:"not null;default:true"`

	// Permissions
	CanManageOrganization bool `json:"can_manage_organization" gorm:"not null;default:false"`
	CanManageMembers      bool `json:"can_manage_members" gorm:"not null;default:false"`
	CanManageProjects     bool `json:"can_manage_projects" gorm:"not null;default:false"`
	CanCreateProjects     bool `json:"can_create_projects" gorm:"not null;default:false"`
	CanViewAllProjects    bool `json:"can_view_all_projects" gorm:"not null;default:false"`
	CanManageTasks        bool `json:"can_manage_tasks" gorm:"not null;default:false"`
	CanViewReports        bool `json:"can_view_reports" gorm:"not null;default:false"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	OrganizationMembers []OrganizationMember `json:"organization_members,omitempty" gorm:"foreignKey:RoleID"`
}

type VWOrganizationMemberRole struct {
	OrganizationID        uint   `json:"organization_id" gorm:"primaryKey"`
	OrganizationName      string `json:"organization_name"`
	Slug                  string `json:"slug"`
	UserID                uint   `json:"user_id"`
	Name                  string `json:"name" gorm:"not null;unique;size:50"`
	Description           string `json:"description" gorm:"type:text"`
	IsDefault             bool   `json:"is_default" gorm:"not null;default:false"`
	IsPreview             bool   `json:"is_preview" gorm:"not null;default:true"`
	CanManageOrganization bool   `json:"can_manage_organization" gorm:"not null;default:false"`
	CanManageMembers      bool   `json:"can_manage_members" gorm:"not null;default:false"`
	CanManageProjects     bool   `json:"can_manage_projects" gorm:"not null;default:false"`
	CanCreateProjects     bool   `json:"can_create_projects" gorm:"not null;default:false"`
	CanViewAllProjects    bool   `json:"can_view_all_projects" gorm:"not null;default:false"`
	CanManageTasks        bool   `json:"can_manage_tasks" gorm:"not null;default:false"`
	CanViewReports        bool   `json:"can_view_reports" gorm:"not null;default:false"`
}
