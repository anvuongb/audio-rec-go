package voice

import (
	"time"
)

const (
	VoiceMetadataTable = "voice_metadata"
)

type VoiceMetadata struct {
	RequestId      string    `json:"request_id"`
	FileId         string    `json:"file_id"`
	FilepathMask   string    `json:"filepath_mask"`
	FilepathNoMask string    `json:"filepath_no_mask"`
	GeneratedText  string    `json:"generated_text"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (VoiceMetadata) TableName() string {
	return VoiceMetadataTable
}
