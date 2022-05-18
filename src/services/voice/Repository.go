package voice

import (
	"audio-rec-go/src/config"
	"fmt"
	"io/ioutil"
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
	SaveAudio(requestId string, fileId string, voiceDecode []byte, masked bool, country string, gender string, maskType string) error
	GetStats() (int, int, int, int)
	GetVoiceFile(fileId string, masked bool) (VoiceFile, error)
	GetVoiceRecords(request GenericRequest) ([]VoiceMetadata, error)
}

func NewRepository(db *gorm.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: logger,
	}
}

var src = rand.NewSource(time.Now().UnixNano())

var sentences = [60]string{
	"A blessing in disguise.",
	"Beat around the bush.",
	"Call it a day.",
	"Get your act together.",
	"Hang in there.",
	"Make a long story short.",
	"Pull yourself together.",
	"The best of both worlds.",
	"To get bent out of shape.",
	"Your guess is as good as mine.",
	"What's the weather like?",
	"When is the train leaving?",
	"I am sorry to interrupt you.",
	"Don’t embarrass me.",
	"He has been to the market.",
	"Don’t talk nonsense.",
	"He always keeps telling a lie.",
	"Key feels like eating funnel cake.",
	"Don’t you disturb me?",
	"I am having my haircut.",
	"He has gone to market.",
	"He is good at English.",
	"I have had breakfast.",
	"He is so called a teacher.",
	"He is very stingy.",
	"I have been out since morning.",
	"Which one of these umbrellas is yours?",
	"I am not a riff-riff person.",
	"You have no manner.",
	"Where does this road lead to?",
	"Why are you staring?",
	"Will you do me a favor?",
	"You are a deceiver and cheater.",
	"What is your intention?",
	"You should see a doctor.",
	"A blessing in disguise.",
	"Beat around the bush.",
	"Call it a day.",
	"Get your act together.",
	"Hand in there.",
	"Make a long story short.",
	"Pull yourself together.",
	"The best of both worlds.",
	"To get bent out of shape.",
	"Your guess is as good as mine.",
	"What is done is done.",
	"What is going on?",
	"There is less salt in the curry.",
	"There is no sense in it.",
	"What a shame!",
	"This is selling at throw away price.",
	"What does it matter?",
	"There is nothing but water.",
	"Two and two make four.",
	"Use your brain sometimes.",
	"Use your brains/use your mind.",
	"That boy David reads in your class.",
	"The examination is about to start.",
	"The train is about to come.",
	"What is the time by your watch?",
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

func (repo repo) GetStats() (int, int, int, int) {
	// query := `
	// 	SELECT count(*) FROM %s WHERE datetime(created_at) >= datetime('now', '-%d hours')
	// `
	time.Now().Add(-1 * time.Hour)
	var count1hour int64
	repo.db.Model(&VoiceMetadata{}).Where("nomasked_file_uploaded = ?", 1).Where("masked_file_uploaded = ?", 1).Where("created_at >= ?", time.Now().Add(-1*time.Hour)).Count(&count1hour)

	var count3hour int64
	repo.db.Model(&VoiceMetadata{}).Where("nomasked_file_uploaded = ?", 1).Where("masked_file_uploaded = ?", 1).Where("created_at >= ?", time.Now().Add(-3*time.Hour)).Count(&count3hour)

	var count24hour int64
	repo.db.Model(&VoiceMetadata{}).Where("nomasked_file_uploaded = ?", 1).Where("masked_file_uploaded = ?", 1).Where("created_at >= ?", time.Now().Add(-24*time.Hour)).Count(&count24hour)

	var count int64
	repo.db.Model(&VoiceMetadata{}).Where("nomasked_file_uploaded = ?", 1).Where("masked_file_uploaded = ?", 1).Count(&count)

	return int(count1hour), int(count3hour), int(count24hour), int(count)
}

func (repo repo) GetVoiceRecords(request GenericRequest) ([]VoiceMetadata, error) {
	logger := log.With(repo.logger, "method", "GetVoiceRecords", "request_id", request.RequestId)

	limit := 0
	if request.Limit <= 0 || request.Limit >= 200 {
		limit = 200
	}

	pageNumber := 0
	recordsPerPage := 0
	if request.PageNumber > 0 {
		pageNumber = request.PageNumber
		if request.RecordsPerPage > 0 {
			recordsPerPage = request.RecordsPerPage
		} else {
			recordsPerPage = 10
		}
	}

	var v []VoiceMetadata
	tx := repo.db.Table(VoiceMetadataTable)
	tx = tx.Where("nomasked_file_uploaded = ?", 1)

	if pageNumber > 0 {
		tx = tx.Limit(recordsPerPage).Offset((pageNumber - 1) * recordsPerPage).Order("created_at desc")
	} else {
		tx = tx.Limit(limit).Order("created_at desc")
	}

	err := tx.Find(&v).Error
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return []VoiceMetadata{}, err
	}
	for i := range v {
		v[i].CreatedAtStr = v[i].CreatedAt.Add(-14 * time.Hour).Format(config.DateLayout)
	}

	return v, nil
}

func (repo repo) GetAudioByFilepath(filepath string) ([]byte, error) {
	logger := log.With(repo.logger, "method", "SaveAudio", "filepath", filepath)
	audio, err := ioutil.ReadFile(filepath)
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return nil, err
	}
	return audio, nil
}

func (repo repo) GetVoiceFile(fileId string, masked bool) (VoiceFile, error) {
	logger := log.With(repo.logger, "method", "CheckIfAudioExists", "file_id", fileId)
	var v []VoiceMetadata
	tx := repo.db.Table(VoiceMetadataTable)
	tx = tx.Where("file_id = ?", fileId)
	if masked {
		tx = tx.Where("masked_file_uploaded = ?", 1)
	} else {
		tx = tx.Where("nomasked_file_uploaded = ?", 1)
	}
	err := tx.Find(&v).Error
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return VoiceFile{}, err
	}
	if len(v) != 1 {
		err := fmt.Errorf("found %d matches file_id, must be 1", len(v))
		level.Error(logger).Log("err", err.Error())
		return VoiceFile{}, err
	}
	var audio []byte
	if masked {
		audio, err = repo.GetAudioByFilepath(v[0].FilepathMask)
		if err != nil {
			level.Error(logger).Log("err", err.Error())
			return VoiceFile{}, err
		}
	} else {
		audio, err = repo.GetAudioByFilepath(v[0].FilepathNoMask)
		if err != nil {
			level.Error(logger).Log("err", err.Error())
			return VoiceFile{}, err
		}
	}
	vf := VoiceFile{
		RequestId:     v[0].RequestId,
		FileId:        v[0].FileId,
		GeneratedText: v[0].GeneratedText,
		Masked:        masked,
		AudioByte:     audio,
		ResultCode:    1,
		ResultMessage: "OK",
		MaskType:      v[0].MaskType,
		Country:       v[0].Country,
		Gender:        v[0].Gender,
		CreatedAtStr:  v[0].CreatedAt.Add(-14 * time.Hour).Format(config.DateLayout),
	}

	return vf, nil
}

func (repo repo) SaveAudio(requestId string, fileId string, voiceDecode []byte, masked bool, country string, gender string, maskType string) error {
	logger := log.With(repo.logger, "method", "SaveAudio", "request_id", requestId)

	// get filepath from db
	var v []VoiceMetadata
	tx := repo.db.Table(VoiceMetadataTable)
	tx = tx.Where("file_id = ?", fileId)
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

	if masked {
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
		err = repo.db.Model(VoiceMetadata{}).Where("file_id = ?", fileId).Update("masked_file_uploaded", 1).Error
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
		// update country at first submission
		err = repo.db.Model(VoiceMetadata{}).Where("file_id = ?", fileId).Updates(VoiceMetadata{NomaskedFileUploaded: 1, Country: country, MaskType: maskType, Gender: gender}).Error
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
		RequestId:            request.RequestId,
		FileId:               fileId,
		FilepathMask:         "data/recordings/" + fileId + "_mask.wav",
		FilepathNoMask:       "data/recordings/" + fileId + "_no_mask.wav",
		GeneratedText:        sentences[rand.Intn(60)],
		MaskedFileUploaded:   0,
		NomaskedFileUploaded: 0,
		MaskType:             "N/A",
		Country:              "N/A",
		Gender:               "N/A",
	}

	// write to db
	err := repo.db.Create(&v).Error
	if err != nil {
		level.Error(logger).Log("err", err.Error())
		return VoiceMetadata{}, err
	}
	return v, nil
}
