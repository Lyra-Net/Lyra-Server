package models

import (
	"github.com/lib/pq"
)

type Song struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"not null" json:"title"`
	Artists    []Artist       `gorm:"many2many:artist_songs;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"artists"`
	TitleToken pq.StringArray `gorm:"type:text[]" json:"title_token"`
	Categories []string       `gorm:"type:text[]" json:"categories"`
}
