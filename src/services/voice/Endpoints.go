package voice

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	InitMetadata         endpoint.Endpoint
	SaveAudio            endpoint.Endpoint
	GetStats             endpoint.Endpoint
	GetAudioByFileIdByte endpoint.Endpoint
	GetAudioByFileId     endpoint.Endpoint
	GetVoiceRecords      endpoint.Endpoint
}

func makeGetVoiceRecordsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenericRequest)
		response, err := s.GetVoiceRecords(ctx, req)
		return response, err
	}
}

func makeGetAudioByFileIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenericRequest)
		response, err := s.GetAudioByFileId(ctx, req)
		return response, err
	}
}

func makeGetAudioByFileIdByteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenericRequest)
		response, err := s.GetAudioByFileIdByte(ctx, req)
		return response, err
	}
}

func makeGetStatsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenericRequest)
		response, err := s.GetStats(ctx, req)
		return response, err
	}
}

func makeInitMetadataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenericRequest)
		response, err := s.InitMetadata(ctx, req)
		return response, err
	}
}

func makeSaveAudioEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VoiceFile)
		response, err := s.SaveAudio(ctx, req)
		return response, err
	}
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		InitMetadata:         makeInitMetadataEndpoint(s),
		SaveAudio:            makeSaveAudioEndpoint(s),
		GetStats:             makeGetStatsEndpoint(s),
		GetAudioByFileIdByte: makeGetAudioByFileIdByteEndpoint(s),
		GetAudioByFileId:     makeGetAudioByFileIdEndpoint(s),
		GetVoiceRecords:      makeGetVoiceRecordsEndpoint(s),
	}
}
