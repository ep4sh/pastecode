package paste

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

const PasteMaxLivenessHours = 72

type Pastecode struct {
	UID       string
	Username  string
	CreatedAt string
	Code      string
}

func NewPastecode(user, code string) (*Pastecode, error) {
	if user == "" {
		return &Pastecode{}, fmt.Errorf("user is missing")
	}

	if code == "" {
		return &Pastecode{}, fmt.Errorf("code is missing")
	}

	return &Pastecode{
		UID:       uuid.New().String(),
		Username:  user,
		CreatedAt: time.Now().UTC().Format(time.DateTime),
		Code:      code,
	}, nil
}

func (p *Pastecode) DeathTime() (bool, error) {
	createdAt, err := time.Parse(time.DateTime, p.CreatedAt)
	if err != nil {
		return false, err
	}

	cutoff := time.Now().UTC().Add(-PasteMaxLivenessHours * time.Hour)

	return createdAt.Before(cutoff), nil
}

type Pastecodes map[string]*Pastecode

func NewPastecodes() Pastecodes {
	m := make(Pastecodes, 8)
	return m
}

func (ps Pastecodes) GC() {
	for id, paste := range ps {
		timeHasCome, err := paste.DeathTime()
		if err != nil {
			log.Printf("GC failed: %s", err)
			continue
		}

		if timeHasCome {
			delete(ps, id)
			log.Printf("GC remove paste: %s, author: %s", id, paste.Username)
		}
	}
}

func (ps Pastecodes) Add(paste *Pastecode) error {
	id := paste.UID
	if id != "" {
		ps[id] = paste
		return nil
	}

	return fmt.Errorf("uuid is empty string")
}

func (ps Pastecodes) FindPaste(id string) (*Pastecode, error) {
	if ps[id] != nil {
		return ps[id], nil
	}

	return nil, fmt.Errorf("paste not found")
}

func ParseUUID(uuidString string) (string, error) {
	u, err := uuid.Parse(uuidString)
	if err != nil {
		return "", fmt.Errorf("invalid uuid")
	}
	return u.String(), nil
}
