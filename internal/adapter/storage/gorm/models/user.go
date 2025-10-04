package models

import (
	"time"
)

type User struct {
	BaseModel

	Email              string     `json:"email" gorm:"uniqueIndex;not null"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	DisplayName        string     `json:"display_name"`
	Bio                string     `json:"bio" gorm:"type:text;default:null"`
	Avatar             string     `json:"avatar" gorm:"default:null"`
	DateOfBirth        *time.Time `json:"date_of_birth"`
	Gender             string     `json:"gender" gorm:"default:null"`
	PhoneNumber        string     `json:"phone_number" gorm:"uniqueIndex;default:null"`
	LanguagePreference string     `json:"language_preference" gorm:"default:'en'"`
	TimeZone           string     `json:"time_zone" gorm:"default:'UTC'"`
	IsEmailVerified    bool       `json:"is_email_verified" gorm:"default:false"`
	IsPhoneVerified    bool       `json:"is_phone_verified" gorm:"default:false"`
	LastLoginAt        *time.Time `json:"last_login_at" gorm:"index"`

	AuthMethods         []UserAuthMethod     `json:"auth_methods,omitempty" gorm:"foreignKey:UserID"`
	OrganizationMembers []OrganizationMember `json:"organization_members,omitempty" gorm:"foreignKey:UserID"`
	Preferences         []UserPreference     `json:"preferences,omitempty" gorm:"foreignKey:UserID"`
	OwnedProjects       []Project            `json:"owned_projects,omitempty" gorm:"foreignKey:OwnerID"`
	ProjectMemberships  []ProjectMember      `json:"project_memberships,omitempty" gorm:"foreignKey:UserID"`
	AssignedTickets     []Ticket             `json:"assigned_tickets,omitempty" gorm:"foreignKey:AssigneeID"`
	ReportedTickets     []Ticket             `json:"reported_tickets,omitempty" gorm:"foreignKey:ReporterID"`
	TicketComments      []TicketComment      `json:"ticket_comments,omitempty" gorm:"foreignKey:UserID"`
	WatchedTickets      []Ticket             `json:"watched_tickets,omitempty" gorm:"many2many:ticket_watchers"`
	TimeLogs            []TimeLog            `json:"time_logs,omitempty" gorm:"foreignKey:UserID"`
	UploadedFiles       []TicketAttachment   `json:"uploaded_files,omitempty" gorm:"foreignKey:UploadedBy"`
}

type UserAuthMethod struct {
	BaseModel

	UserID       uint    `json:"user_id" gorm:"not null;index"`
	AuthType     string  `json:"auth_type" gorm:"not null;index"`
	AuthProvider *string `json:"auth_provider" gorm:"index;default:null"`
	ProviderID   *string `json:"provider_id" gorm:"default:null"`
	IsPrimary    bool    `json:"is_primary" gorm:"default:false"`

	PasswordHash string `json:"-" gorm:"default:null"`
	PasswordSalt string `json:"-" gorm:"default:null"`

	AccessToken  *string    `json:"-" gorm:"default:null"`
	RefreshToken *string    `json:"-" gorm:"default:null"`
	TokenExpiry  *time.Time `json:"-" gorm:"default:null"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type UserPreference struct {
	BaseModel

	UserID  uint   `json:"user_id" gorm:"not null;index"`
	Key     string `json:"key" gorm:"not null;index"`
	Value   string `json:"value" gorm:"type:text"`
	Context string `json:"context" gorm:"index"`

	// Relationships
	// User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}


