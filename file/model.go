package file

import (
	"fmt"
	"time"
)

type Token struct {
	ID     uint64 `gorm:"primary_key"`
	UUID   string `gorm:"not null"`
	Salt   string `gorm:"not null;type:varchar(32)"`
	Name   string `gorm:"null;type:varchar(255)"`
	UserID uint64 `gorm:"type:bigint not null REFERENCES users(id)"`

	WriteOnly        bool   `gorm:"not null"`
	NumberOfArchives uint64 `gorm:"not null"`
	SizeUsed         uint64 `gorm:"not null;default:0"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type FileMetadata struct {
	ID            uint64    `gorm:"primary_key" json:"-"`
	TokenID       uint64    `gorm:"type:bigint not null REFERENCES tokens(id)" json:"-"`
	LastVersionID uint64    `gorm:"type:bigint" json:"-"`
	Path          string    `gorm:"not null" json:"path"`
	Filename      string    `gorm:"not null" json:"filename"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (fm *FileMetadata) FullPath() string {
	if fm.Path == "/" {
		return fmt.Sprintf("/%s", fm.Filename)
	}
	return fmt.Sprintf("%s/%s", fm.Path, fm.Filename)
}

type FileVersion struct {
	ID             uint64 `gorm:"primary_key" json:"id"`
	DuplicatedSave uint64 `gorm:"not null" json:"duplicate_save"`

	// denormalized for simplicity
	TokenID        uint64 `gorm:"type:bigint not null REFERENCES tokens(id)" json:"-"`
	FileMetadataID uint64 `gorm:"type:bigint not null REFERENCES file_metadata(id)" json:"-"`

	Size    uint64 `gorm:"not null" json:"size"`
	SHA256  string `gorm:"not null" json:"sha"`
	StoreID string `gorm:"not null" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
