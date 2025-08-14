package models

type Artist struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"unique;not null" json:"name"`
	Songs []Song `gorm:"many2many:artist_songs;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"songs"`
}
