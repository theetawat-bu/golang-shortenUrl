package model

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	UrlFull  string `json:"url_full"`
	UrlShort string `json:"url_short" gorm:"unique"`
}
