package link

import (
	"math/rand"
	"url_shortener/internal/stat"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url   string      `json:"url"`
	Hash  string      `json:"hash" gorm:"uniqueIndex"`
	Stats []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewLink(url string) *Link {
	link := &Link{
		Url: url,
	}
	link.GenerateHash()
	return link
}

var letterRunes = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")

func RandStringRunes(n int) string {
	s := make([]rune, n)
	l := len(letterRunes)
	for i := range s {
		s[i] = letterRunes[rand.Intn(l)]
	}
	return string(s)
}

func (l *Link) GenerateHash() {
	l.Hash = RandStringRunes(6)
}
