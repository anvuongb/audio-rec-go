package voice

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

type Service interface {
	InitMetadata(ctx context.Context, request GenericRequest) (GenericResponse, error)
	SaveAudio(ctx context.Context, request VoiceFile) (GenericResponse, error)
	GetStats(ctx context.Context, request GenericRequest) (StatsResponse, error)
	GetAudioByFileId(ctx context.Context, request GenericRequest) (VoiceFile, error)
	GetAudioByFileIdByte(ctx context.Context, request GenericRequest) ([]byte, error)
	GetVoiceRecords(ctx context.Context, request GenericRequest) ([]VoiceMetadata, error)
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
	level.Error(logger).Log("info", "received")
	v, err := s.repository.InitMetadata(request)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return GenericResponse{RequestId: request.RequestId, ResultCode: -1, ResultMessage: err.Error()}, nil
	}
	return GenericResponse{RequestId: request.RequestId, ResultCode: 1, ResultMessage: "OK", GeneratedText: v.GeneratedText, FileId: v.FileId}, nil
}

func (s service) GetVoiceRecords(ctx context.Context, request GenericRequest) ([]VoiceMetadata, error) {
	logger := log.With(s.logger, "method", "GetVoiceRecords", "request_id", request.RequestId)
	level.Error(logger).Log("page_number", request.PageNumber, "records_per_page", request.RecordsPerPage, "limit", request.Limit)
	v, err := s.repository.GetVoiceRecords(request)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return []VoiceMetadata{}, err
	}
	return v, nil
}

func (s service) SaveAudio(ctx context.Context, request VoiceFile) (GenericResponse, error) {
	logger := log.With(s.logger, "method", "SaveAudio", "request_id", request.RequestId)
	level.Error(logger).Log("info", "received")
	// level.Error(logger).Log("file", request.File)
	voiceBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(voiceBuffer, request.File); err != nil {
		level.Error(logger).Log("error", err)
		return GenericResponse{}, fmt.Errorf("failed decoding voice")
	}

	voiceDecode := voiceBuffer.Bytes()

	err := s.repository.SaveAudio(request.RequestId, request.FileId, voiceDecode, request.Masked)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return GenericResponse{RequestId: request.RequestId, ResultCode: -1, ResultMessage: err.Error()}, nil
	}

	return GenericResponse{RequestId: request.RequestId, ResultCode: 1, ResultMessage: "OK"}, nil
}

func (s service) GetStats(ctx context.Context, request GenericRequest) (StatsResponse, error) {
	c1, c2, c3, c4 := s.repository.GetStats()
	return StatsResponse{Count1Hour: c1, Count2Hour: c2, Count24Hour: c3, Count: c4}, nil
}

func (s service) GetAudioByFileId(ctx context.Context, request GenericRequest) (VoiceFile, error) {
	logger := log.With(s.logger, "method", "GetAudioByFileId", "request_id", request.RequestId)
	level.Error(logger).Log("info", "received")
	vf, err := s.repository.GetVoiceFile(request.FileId, request.Masked)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return VoiceFile{RequestId: request.RequestId, ResultCode: -1, ResultMessage: err.Error()}, nil
	}
	return vf, nil
}

func (s service) GetAudioByFileIdByte(ctx context.Context, request GenericRequest) ([]byte, error) {
	logger := log.With(s.logger, "method", "GetAudioByFileId", "request_id", request.RequestId)
	level.Error(logger).Log("info", "received")
	vf, err := s.repository.GetVoiceFile(request.FileId, request.Masked)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return nil, err
	}
	return vf.AudioByte, nil
}
