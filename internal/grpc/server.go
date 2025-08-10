package grpc

import (
	"context"
	"errors"
	"net/url"

	"github.com/google/uuid"
	newsv1 "github.com/sanjevscet/news-grpc/api/news/v1"
	"github.com/sanjevscet/news-grpc/internal/memstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewsServiceServer is the server API for NewsService service.

type NewsStorer interface {
	Create(news *memstore.News) *memstore.News
	Get(id uuid.UUID) *memstore.News
}
type Server struct {
	newsv1.UnimplementedNewsServiceServer
	storer NewsStorer
}

func NewServer(storer NewsStorer) *Server {
	return &Server{
		storer: storer,
	}
}

func (s *Server) Create(_ context.Context, request *newsv1.CreateRequest) (*newsv1.CreateResponse, error) {
	parsedNews, errs := parseAndValidate(request)
	if errs != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", errs)
	}

	createdNews := s.storer.Create(parsedNews)
	return toNewsResponse(createdNews), nil
}
func (s *Server) Get(_ context.Context, request *newsv1.GetRequest) (*newsv1.GetResponse, error) {
	newsUUID, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ID format: %v", err)
	}
	fetchedNews := s.storer.Get(newsUUID)
	if fetchedNews == nil {
		return nil, status.Errorf(codes.NotFound, "news with ID %s not found", request.Id)
	}

	return &newsv1.GetResponse{
		Id:        fetchedNews.ID.String(),
		Author:    fetchedNews.Author,
		Title:     fetchedNews.Title,
		Content:   fetchedNews.Content,
		Source:    fetchedNews.Source.String(),
		Tags:      fetchedNews.Tags,
		CreatedAt: timestamppb.New(fetchedNews.CreatedAt.UTC()),
		UpdatedAt: timestamppb.New(fetchedNews.UpdatedAt.UTC()),
		DeletedAt: timestamppb.New(fetchedNews.DeletedAt.UTC()),
	}, nil
}

func parseAndValidate(in *newsv1.CreateRequest) (n *memstore.News, errs error) {
	if in == nil {
		return nil, errors.New("request cannot be nil")
	}
	if in.Author == "" {
		errs = errors.Join(errs, errors.New("author cannot be empty"))
	}
	if in.Summary == "" {
		errs = errors.Join(errs, errors.New("summary cannot be empty"))
	}
	if in.Content == "" {
		errs = errors.Join(errs, errors.New("content cannot be empty"))
	}

	if len(in.Tags) == 0 {
		errs = errors.Join(errs, errors.New("at least one tag is required"))
	}

	parsedID, err := uuid.Parse(in.Id)
	if err != nil {
		errs = errors.Join(errs, errors.New("invalid ID format"))
	}

	parsedUrl, err := url.Parse(in.Source)
	if err != nil {
		errs = errors.Join(errs, errors.New("invalid URL format"))
	}

	if err != nil {
		return nil, errs
	}

	return &memstore.News{
		ID:      parsedID,
		Author:  in.Author,
		Title:   in.Title,
		Content: in.Content,
		Source:  parsedUrl,
		Tags:    in.Tags,
	}, nil
}

func toNewsResponse(news *memstore.News) *newsv1.CreateResponse {
	if news == nil {
		return nil
	}

	return &newsv1.CreateResponse{
		Id:        news.ID.String(),
		Author:    news.Author,
		Title:     news.Title,
		Content:   news.Content,
		Source:    news.Source.String(),
		Tags:      news.Tags,
		CreatedAt: timestamppb.New(news.CreatedAt.UTC()),
	}
}
