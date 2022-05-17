package voice

import (
	"mime/multipart"
	"time"
)

type GenericResponse struct {
	RequestId     string `json:"request_id,omitempty" gorm:"-"`
	ResultCode    int    `json:"result_code,omitempty"`
	ResultMessage string `json:"result_message,omitempty"`
	GeneratedText string `json:"generated_text,omitempty"`
	FileId        string `json:"file_id,omitempty"`
}

type GenericRequest struct {
	RequestId      string    `json:"request_id,omitempty" gorm:"-"`
	FileId         string    `json:"file_id,omitempty"`
	Masked         bool      `json:"masked,omitempty"`
	PageNumber     int       `json:"page,omitempty"`
	Limit          int       `json:"limit,omitempty"`
	RecordsPerPage int       `json:"records_per_page,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

type StatsResponse struct {
	Count1Hour  int `json:"count_1_hour"`
	Count2Hour  int `json:"count_2_hour"`
	Count24Hour int `json:"count_24_hour"`
	Count       int `json:"count_total"`
}

type VoiceFile struct {
	RequestId     string         `json:"request_id,omitempty" gorm:"-"`
	FileId        string         `json:"file_id,omitempty"`
	File          multipart.File `json:"file,omitempty"`
	GeneratedText string         `json:"generated_text,omitempty"`
	SoundRate     int32          `json:"sound_rate,omitempty"`
	Masked        bool           `json:"masked,omitempty"`
	AudioByte     []byte         `json:"audio_byte,omitempty"`
	ResultCode    int            `json:"result_code,omitempty"`
	ResultMessage string         `json:"result_message,omitempty"`
}
