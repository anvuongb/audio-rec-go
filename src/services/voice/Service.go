package voice

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

type Service interface {
	InitMetadata(ctx context.Context, request GenericRequest) (GenericResponse, error)
	SaveAudio(ctx context.Context, request VoiceFile) (GenericResponse, error)
}

type service struct {
	repository Repository
	logger     log.Logger
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s service) InitMetadata(ctx context.Context, request GenericRequest) (GenericResponse, error) {
	logger := log.With(s.logger, "method", "InitMetadata", "request_id", request.RequestId)
	v, err := s.repository.InitMetadata(request)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return GenericResponse{RequestId: request.RequestId, ResultCode: -1, ResultMessage: err.Error()}, nil
	}
	return GenericResponse{RequestId: request.RequestId, ResultCode: 1, ResultMessage: "OK", GeneratedText: v.GeneratedText, FileId: v.FileId}, nil
}

func (s service) SaveAudio(ctx context.Context, request VoiceFile) (GenericResponse, error) {
	logger := log.With(s.logger, "method", "SaveAudio", "request_id", request.RequestId)
	err := s.repository.SaveAudio(request)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return GenericResponse{RequestId: request.RequestId, ResultCode: -1, ResultMessage: err.Error()}, nil
	}
	return GenericResponse{RequestId: request.RequestId, ResultCode: 1, ResultMessage: "OK"}, nil
}
