package voice

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"gorm.io/gorm"
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

type Repository interface {
	InitMetadata(request GenericRequest) (VoiceMetadata, error)
	SaveAudio(request VoiceFile) error
}

func NewRepository(db *gorm.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: logger,
	}
}

var src = rand.NewSource(time.Now().UnixNano())

var sentences = [10]string{
	"A blessing in disguise",
	"Beat around the bush",
	"Call it a day",
	"Get your act together",
	"Hang in there",
	"Make a long story short",
	"Pull yourself together",
	"The best of both worlds",
	"To get bent out of shape",
	"Your guess is as good as mine",
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func (repo repo) SaveAudio(request VoiceFile) error {
	logger := log.With(repo.logger, "method", "SaveAudio", "request_id", request.RequestId)

	// get filepath from db
	var v []VoiceMetadata
	tx := repo.db.Table(VoiceMetadataTable)
	tx = tx.Where("file_id = ?", request.FileId)
	err := tx.Find(&v).Error
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return err
	}
	if len(v) != 1 {
		err := fmt.Errorf("found %d matches file_id, must be 1", len(v))
		level.Error(logger).Log("err", err.Error())
		return err
	}
	filepathMask := v[0].FilepathMask
	filepathNoMask := v[0].FilepathNoMask

	voiceBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(voiceBuffer, request.File); err != nil {
		level.Error(logger).Log("error", err)
		return fmt.Errorf("failed decoding voice")
	}

	voiceDecode := voiceBuffer.Bytes()

	if request.Masked {
		file, err := os.OpenFile(
			filepathMask,
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0666,
		)
		if err != nil {
			level.Error(logger).Log("err", err.Error())
			return err
		}
		defer file.Close()
		_, err = file.Write(voiceDecode)
		if err != nil {
			level.Error(logger).Log("err", err.Error())
			return err
		}
	} else {
		file, err := os.OpenFile(
			filepathNoMask,
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0666,
		)
		if err != nil {
			level.Error(logger).Log("err", err.Error())
			return err
		}
		defer file.Close()
		_, err = file.Write(voiceDecode)
		if err != nil {
			level.Error(logger).Log("err", err.Error())
			return err
		}
	}
	return nil
}

func (repo repo) InitMetadata(request GenericRequest) (VoiceMetadata, error) {
	logger := log.With(repo.logger, "method", "InitMetadata", "request_id", request.RequestId)
	// generate request id
	// requestId := uuid.NewV4().String()

	// generate file id
	fileId := RandStringBytesMaskImprSrcSB(10)

	// create record
	v := VoiceMetadata{
		RequestId:      request.RequestId,
		FileId:         fileId,
		FilepathMask:   "data/recordings/" + fileId + "_mask.wav",
		FilepathNoMask: "data/recordings/" + fileId + "_no_mask.wav",
		GeneratedText:  sentences[rand.Intn(10)],
	}

	// write to db
	err := repo.db.Create(&v).Error
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return VoiceMetadata{}, err
	}
	return v, nil
}
