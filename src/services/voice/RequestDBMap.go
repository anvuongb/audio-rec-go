package voice

import (
	"time"
)

const (
	VoiceMetadataTable = "voice_metadata"
)

type VoiceMetadata struct {
	RequestId            string    `json:"request_id"`
	FileId               string    `json:"file_id"`
	FilepathMask         string    `json:"filepath_mask"`
	FilepathNoMask       string    `json:"filepath_no_mask"`
	GeneratedText        string    `json:"generated_text"`
	MaskedFileUploaded   int       `json:"masked_file_uploaded"`
	NomaskedFileUploaded int       `json:"nomasked_file_uploaded"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedAtStr         string    `json:"created_at_str" gorm:"-"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (VoiceMetadata) TableName() string {
	return VoiceMetadataTable
}
