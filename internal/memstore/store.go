package memstore

import (
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID        uuid.UUID
	Author    string
	Title     string
	Summary   string
	Content   string
	Tags      []string
	Source    *url.URL
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Store struct {
	lock sync.RWMutex
	news []*News
}

func NewStore() *Store {
	return &Store{
		news: make([]*News, 0),
		lock: sync.RWMutex{},
	}
}

func (s *Store) Create(new *News) *News {
	newNews := &News{
		ID:        uuid.New(),
		Author:    new.Author,
		Title:     new.Title,
		Summary:   new.Summary,
		Content:   new.Content,
		Tags:      new.Tags,
		Source:    new.Source,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.news = append(s.news, newNews)

	return newNews

}

func (s *Store) Get(id uuid.UUID) *News {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, news := range s.news {
		if news.ID == id && news.DeletedAt.IsZero() {
			return news
		}
	}
	return nil
}

func (s *Store) GetAll() []*News {
	s.lock.RLock()
	defer s.lock.RUnlock()

	result := make([]*News, 0)
	for _, news := range s.news {
		if news.DeletedAt.IsZero() {
			result = append(result, news)
		}
	}
	return result
}
