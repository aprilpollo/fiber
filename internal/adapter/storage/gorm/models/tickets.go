package models

import (
	"time"
)

type Ticket struct {
	BaseModel

	ProjectID      uint       `json:"project_id" gorm:"not null;index"`
	Title          string     `json:"title" gorm:"not null;size:500"`
	Description    string     `json:"description" gorm:"type:text"`
	TicketKey      string     `json:"ticket_key" gorm:"not null;unique;size:20"`
	TypeID         uint       `json:"type_id" gorm:"not null;index;default:1"`
	StatusID       uint       `json:"status_id" gorm:"not null;index;default:1"`
	PriorityID     uint       `json:"priority_id" gorm:"not null;index;default:2"` // FK to priorities (2=medium)
	AssigneeID     *uint      `json:"assignee_id" gorm:"index"`
	ReporterID     uint       `json:"reporter_id" gorm:"not null;index"`
	ParentID       *uint      `json:"parent_id" gorm:"index"`
	EstimatedHours *float64   `json:"estimated_hours"`
	ActualHours    *float64   `json:"actual_hours"`
	DueDate        *time.Time `json:"due_date"`

	// Story Points for Agile
	StoryPoints *int `json:"story_points"`

	// Resolution
	ResolutionID *uint      `json:"resolution_id" gorm:"index"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	ResolvedBy   *uint      `json:"resolved_by" gorm:"index"`

	// Relationships
	Project     Project            `json:"project" gorm:"foreignKey:ProjectID"`
	Type        TicketType         `json:"type" gorm:"foreignKey:TypeID"`
	Status      TicketStatus       `json:"status" gorm:"foreignKey:StatusID"`
	Priority    Priority           `json:"priority" gorm:"foreignKey:PriorityID"`
	Resolution  *Resolution        `json:"resolution,omitempty" gorm:"foreignKey:ResolutionID"`
	Assignee    *User              `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	Reporter    User               `json:"reporter" gorm:"foreignKey:ReporterID"`
	Resolver    *User              `json:"resolver,omitempty" gorm:"foreignKey:ResolvedBy"`
	Parent      *Ticket            `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []Ticket           `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Comments    []TicketComment    `json:"comments,omitempty"`
	Attachments []TicketAttachment `json:"attachments,omitempty"`
	Labels      []Label            `json:"labels,omitempty" gorm:"many2many:ticket_labels"`
	Components  []Component        `json:"components,omitempty" gorm:"many2many:ticket_components"`
	//Versions    []Version          `json:"versions,omitempty" gorm:"many2many:ticket_versions"`
	Watchers    []User             `json:"watchers,omitempty" gorm:"many2many:ticket_watchers"`
	TimeLogs    []TimeLog          `json:"time_logs,omitempty"`
}

type TicketStatus struct {
	BaseModel

	ProjectID   uint   `json:"project_id" gorm:"not null;index;uniqueIndex:idx_project_status_name"`
	Name        string `json:"name" gorm:"not null;size:50;uniqueIndex:idx_project_status_name"` // unique per project
	Description string `json:"description" gorm:"type:text"`
	Color       string `json:"color" gorm:"size:7;default:'#42526E'"`
	Position    int    `json:"position" gorm:"not null;default:0;index"`
	IsDefault   bool   `json:"is_default" gorm:"default:false"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`

	// Status behavior flags
	IsClosed   bool `json:"is_closed" gorm:"default:false"`   
	IsResolved bool `json:"is_resolved" gorm:"default:false"` 

	// Relationships
	Project Project  `json:"project" gorm:"foreignKey:ProjectID"`
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:StatusID"`
}

type TicketComment struct {
	BaseModel

	TicketID uint   `json:"ticket_id" gorm:"not null;index"`
	UserID   uint   `json:"user_id" gorm:"not null;index"`
	Content  string `json:"content" gorm:"not null;type:text"`

	// Relationships
	Ticket Ticket `json:"ticket" gorm:"foreignKey:TicketID"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
}

type TicketAttachment struct {
	BaseModel

	TicketID   uint   `json:"ticket_id" gorm:"not null;index"`
	Filename   string `json:"filename" gorm:"not null;size:255"`
	FilePath   string `json:"file_path" gorm:"not null;size:500"`
	FileSize   int64  `json:"file_size" gorm:"not null"`
	MimeType   string `json:"mime_type" gorm:"size:100"`
	UploadedBy uint   `json:"uploaded_by" gorm:"not null;index"`

	// Relationships
	Ticket   Ticket `json:"ticket" gorm:"foreignKey:TicketID"`
	Uploader User   `json:"uploader" gorm:"foreignKey:UploadedBy"`
}

type Label struct {
	BaseModel

	ProjectID   uint   `json:"project_id" gorm:"not null;index"`
	Name        string `json:"name" gorm:"not null;size:100;index"`
	Color       string `json:"color" gorm:"size:7"` // hex color
	Description string `json:"description" gorm:"size:500"`

	// Relationships
	Project Project  `json:"project" gorm:"foreignKey:ProjectID"`
	Tickets []Ticket `json:"tickets,omitempty" gorm:"many2many:ticket_labels"`
}

type TimeLog struct {
	BaseModel

	TicketID    uint      `json:"ticket_id" gorm:"not null;index"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Hours       float64   `json:"hours" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	LoggedDate  time.Time `json:"logged_date" gorm:"not null;index"`

	// Relationships
	Ticket Ticket `json:"ticket" gorm:"foreignKey:TicketID"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
}

type Priority struct {
	ID          uint   `json:"id" gorm:"uniqueIndex;not null"`
	Name        string `json:"name" gorm:"not null;unique;size:50"` // low, medium, high, critical, blocker
	Description string `json:"description" gorm:"type:text"`
	Color       string `json:"color" gorm:"size:7"`          // hex color for UI
	Level       int    `json:"level" gorm:"not null;unique"` // 1=Low, 2=Medium, 3=High, 4=Critical, 5=Blocker
	IsActive    bool   `json:"is_active" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Relationships
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:PriorityID"`
}

type TicketType struct {
	BaseModel

	Name        string `json:"name" gorm:"not null;unique;size:50"` // task, bug, feature, improvement, epic
	Description string `json:"description" gorm:"type:text"`
	Icon        string `json:"icon" gorm:"size:50"` // icon name for UI
	Color       string `json:"color" gorm:"size:7"` // hex color for UI
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	Position    int    `json:"position" gorm:"default:0;index"`

	// Relationships
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:TypeID"`
}

type Component struct {
	BaseModel

	ProjectID   uint   `json:"project_id" gorm:"not null;index"`
	Name        string `json:"name" gorm:"not null;size:100"`
	Description string `json:"description" gorm:"type:text"`
	LeadID      *uint  `json:"lead_id" gorm:"index"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`

	// Relationships
	Project Project  `json:"project" gorm:"foreignKey:ProjectID"`
	Lead    *User    `json:"lead,omitempty" gorm:"foreignKey:LeadID"`
	Tickets []Ticket `json:"tickets,omitempty" gorm:"many2many:ticket_components"`
}

// type Version struct {
// 	BaseModel

// 	ProjectID   uint       `json:"project_id" gorm:"not null;index"`
// 	Name        string     `json:"name" gorm:"not null;size:100"`
// 	Description string     `json:"description" gorm:"type:text"`
// 	ReleaseDate *time.Time `json:"release_date"`
// 	IsReleased  bool       `json:"is_released" gorm:"default:false"`
// 	IsActive    bool       `json:"is_active" gorm:"default:true"`

// 	// Relationships
// 	Project Project  `json:"project" gorm:"foreignKey:ProjectID"`
// 	Tickets []Ticket `json:"tickets,omitempty" gorm:"many2many:ticket_versions"`
// }

type Resolution struct {
	BaseModel

	Name        string `json:"name" gorm:"not null;unique;size:50"`
	Description string `json:"description" gorm:"type:text"`
	IsDefault   bool   `json:"is_default" gorm:"default:false"`

	// Relationships
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:ResolutionID"`
}
