package voice

import "mime/multipart"

type GenericResponse struct {
	RequestId     string `json:"request_id,omitempty" gorm:"-"`
	ResultCode    int    `json:"result_code,omitempty"`
	ResultMessage string `json:"result_message,omitempty"`
	GeneratedText string `json:"generated_text,omitempty"`
	FileId        string `json:"file_id,omitempty"`
}

type GenericRequest struct {
	RequestId string `json:"request_id,omitempty" gorm:"-"`
}

type VoiceFile struct {
	RequestId     string         `json:"request_id,omitempty" gorm:"-"`
	FileId        string         `json:"file_id,omitempty"`
	File          multipart.File `json:"file"`
	GeneratedText string         `json:"generated_text"`
	SoundRate     int32          `json:"sound_rate"`
	Masked        bool           `json:"masked"`
}
