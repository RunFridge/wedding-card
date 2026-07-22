package demo

import (
	"bytes"
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/handlers"
	"github.com/RunFridge/wedding-card/internal/imaging"
	"github.com/RunFridge/wedding-card/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	minDaysUntilDemoWedding = 28
	demoWeddingHour         = 13
	gameAssetPhotoCount     = 6
	statsDays               = 14
	demoUploadWindowHours   = "2160"
	seedFilename            = "demo.jpg"
	jpegContentType         = "image/jpeg"
)

//go:embed main_photo.jpg
var mainPhotoJPEG []byte

var guestbookSeeds = []struct {
	nickname string
	message  string
	secret   bool
	hoursAgo int
}{
	{"수진", "민준아 서아야 결혼 축하해!! 너희 둘 소개해준 보람이 있다 😆 오래오래 행복하게 잘 살아~", false, 310},
	{"동현", "드디어 가는구나. 축하한다 친구야. 결혼식 사회 실수 안 하게 열심히 연습 중이다", false, 288},
	{"서아 이모", "우리 서아 어릴 때부터 봐왔는데 벌써 시집을 가네. 두 사람 서로 아껴주며 예쁘게 살아라. 이모가 많이 사랑한다", false, 250},
	{"김대리", "과장님 결혼 진심으로 축하드립니다! 신혼여행 다녀오시면 밥 꼭 사주세요 🙏", false, 221},
	{"하늘", "언니 웨딩사진 실화야...? 세상에서 제일 예쁜 신부다 💐", false, 190},
	{"준영", "축하해!! 축의금은 마음만큼 넣었다. 진짜다", false, 160},
	{"예린", "결혼식에서 부케 나 받기로 한 거 잊지 마 언니 ㅋㅋㅋ", false, 130},
	{"태우", "군대 동기 결혼이라니 감회가 새롭다. 진심으로 축하한다!", false, 101},
	{"민지", "서아야 나 그날 부산에서 첫차 타고 올라가. 조금 늦더라도 꼭 갈게. 미리 축하해, 사랑해", true, 76},
	{"박영길", "김진호 친구 아들 결혼을 축하하네. 두 사람 앞날에 행복만 가득하길 바라네", false, 52},
	{"지원", "둘이 연애할 때부터 지켜본 사람으로서... 이 결혼 인정입니다 ✌️", false, 30},
	{"현우", "형 결혼 축하해요. 형수님한테 잘하세요, 안 그러면 제가 다 이를 거예요 ㅋㅋ", true, 9},
}

var gameScoreSeeds = []struct {
	nickname string
	timeMs   int
}{
	{"SEO", 15400},
	{"KIM", 17250},
	{"LEE", 18930},
	{"PRK", 21080},
	{"CHO", 23640},
	{"JUN", 25120},
	{"HAN", 27450},
	{"YUN", 29900},
}

var hallOfFameSeeds = []string{"수진", "동현", "하늘"}

var guestPhotoSeeds = []struct {
	name   string
	hearts int64
}{
	{"정하늘", 23},
	{"박준영", 17},
	{"최예린", 12},
	{"김태우", 9},
	{"이민지", 8},
	{"한지원", 5},
}

var pageViewSeeds = []int{34, 41, 38, 52, 47, 63, 88, 74, 69, 91, 105, 96, 122, 137}
var gameBeatSeeds = []int{3, 5, 4, 7, 6, 9, 14, 11, 10, 15, 18, 16, 21, 24}

func seed() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(Password), config.Cfg.BcryptCost)
	if err != nil {
		return err
	}

	hash := string(passwordHash)
	if err := seedWeddingConfig(); err != nil {
		return err
	}
	if err := seedGuestbook(hash); err != nil {
		return err
	}
	if err := seedGameScores(); err != nil {
		return err
	}
	if err := seedHallOfFame(); err != nil {
		return err
	}
	if err := seedStats(); err != nil {
		return err
	}
	if err := seedAssetPhotos(); err != nil {
		return err
	}
	return seedGuestPhotos(hash)
}

func seedWeddingConfig() error {
	subwayInfo, err := json.Marshal([]string{"6호선 월드컵경기장역 2번 출구에서 도보 7분"})
	if err != nil {
		return err
	}
	busInfo, err := json.Marshal([]config.BusInfoEntry{
		{Stop: "월드컵공원 정류장", Routes: "271, 710, 7715"},
	})
	if err != nil {
		return err
	}

	return models.SetConfigOverrides(map[string]string{
		"groom_kor_name":            "밥",
		"bride_kor_name":            "앨리스",
		"groom_eng_name":            "Bob",
		"bride_eng_name":            "Alice",
		"groom_father_kor_name":     "데이브",
		"groom_mother_kor_name":     "캐럴",
		"bride_father_kor_name":     "트렌트",
		"bride_mother_kor_name":     "페기",
		"groom_bank_account":        "데모은행 110-1234-567890",
		"bride_bank_account":        "데모은행 220-9876-543210",
		"wedding_datetime":          nextDemoWeddingDate(),
		"venue_name":                "달빛정원 웨딩홀",
		"venue_address":             "서울특별시 마포구 월드컵북로 400",
		"venue_floor":               "5층",
		"venue_hall":                "가든홀",
		"venue_phone":               "02-1234-5678",
		"short_greeting":            "저희, 결혼합니다",
		"main_greet_text":           "서로의 하루를 궁금해하던 마음이\n평생을 함께하고 싶은 약속이 되었습니다.\n저희 두 사람의 시작을 함께 축복해 주세요.",
		"car_info":                  "웨딩홀 지하 주차장 2시간 무료",
		"subway_info":               string(subwayInfo),
		"bus_info":                  string(busInfo),
		"photo_upload_enabled":      "true",
		"photo_upload_hours_before": demoUploadWindowHours,
	})
}

func nextDemoWeddingDate() string {
	kst := time.FixedZone("KST", 9*60*60)
	day := time.Now().In(kst).AddDate(0, 0, minDaysUntilDemoWedding)
	for day.Weekday() != time.Saturday {
		day = day.AddDate(0, 0, 1)
	}
	return time.Date(day.Year(), day.Month(), day.Day(), demoWeddingHour, 0, 0, 0, kst).Format(time.RFC3339)
}

func seedGuestbook(passwordHash string) error {
	for i, s := range guestbookSeeds {
		entry, err := models.CreateGuestbookEntry(s.nickname, s.message, config.IPHash(demoIP(i)), passwordHash, s.secret)
		if err != nil {
			return err
		}
		if err := models.SetGuestbookEvaluated(entry.ID, true, false); err != nil {
			return err
		}
		if _, err := database.DB.Exec(
			`UPDATE guestbook_entries SET created_at = datetime('now', ?) WHERE id = ?`,
			fmt.Sprintf("-%d hours", s.hoursAgo), entry.ID,
		); err != nil {
			return err
		}
	}
	return nil
}

func seedGameScores() error {
	for i, s := range gameScoreSeeds {
		if _, err := models.CreateGameScore(s.nickname, s.timeMs, config.IPHash(demoIP(i))); err != nil {
			return err
		}
	}
	return nil
}

func seedHallOfFame() error {
	for i, nickname := range hallOfFameSeeds {
		if _, err := models.CreateHallOfFameEntry(nickname, config.IPHash(demoIP(i))); err != nil {
			return err
		}
	}
	return nil
}

func seedStats() error {
	for i := 0; i < statsDays; i++ {
		date := time.Now().AddDate(0, 0, i-statsDays+1).Format("2006-01-02")
		if _, err := database.DB.Exec(`INSERT INTO page_views (date, count) VALUES (?, ?)`, date, pageViewSeeds[i]); err != nil {
			return err
		}
		if _, err := database.DB.Exec(`INSERT INTO game_beats (date, count) VALUES (?, ?)`, date, gameBeatSeeds[i]); err != nil {
			return err
		}
	}
	return nil
}

func seedAssetPhotos() error {
	mainID, err := createAssetPhoto("메인 사진", mainPhotoJPEG)
	if err != nil {
		return err
	}
	if err := models.SetAssetPhotoAsMain(mainID); err != nil {
		return err
	}

	for n := 1; n <= gameAssetPhotoCount; n++ {
		data, err := numberJPEG(n)
		if err != nil {
			return err
		}
		id, err := createAssetPhoto("카드 "+strconv.Itoa(n), data)
		if err != nil {
			return err
		}
		if err := models.SetAssetPhotoGameFlag(id, true); err != nil {
			return err
		}
	}
	return nil
}

func createAssetPhoto(label string, data []byte) (int64, error) {
	processed, err := imaging.ProcessAssetUpload(bytes.NewReader(data))
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	base := newObjectName()
	hashname := base + ".jpg"
	thumbHashname := base + "_thumb.jpg"
	if err := handlers.Store.Upload(ctx, "assets/"+hashname, bytes.NewReader(processed.Optimized), jpegContentType); err != nil {
		return 0, err
	}
	if err := handlers.Store.Upload(ctx, "assets/"+thumbHashname, bytes.NewReader(processed.Thumbnail), jpegContentType); err != nil {
		return 0, err
	}

	return models.CreateAssetPhoto(label, hashname, thumbHashname, seedFilename)
}

func seedGuestPhotos(passwordHash string) error {
	ctx := context.Background()
	hearts := make(map[int64]int64, len(guestPhotoSeeds))
	for i, s := range guestPhotoSeeds {
		data, err := placeholderJPEG(i)
		if err != nil {
			return err
		}
		processed, err := imaging.ProcessUpload(bytes.NewReader(data))
		if err != nil {
			return err
		}

		hashname := newObjectName() + ".jpg"
		if err := handlers.Store.Upload(ctx, "photos/"+hashname, bytes.NewReader(processed.Optimized), jpegContentType); err != nil {
			return err
		}

		id, err := models.CreatePhotoUpload(s.name, config.IPHash(demoIP(i)), hashname, "", seedFilename, processed.Thumbnail, passwordHash)
		if err != nil {
			return err
		}
		if err := models.SetPhotoEvaluated(id, true, false); err != nil {
			return err
		}
		hearts[id] = s.hearts
	}
	return models.IncrementPhotoHearts(hearts)
}

func demoIP(i int) string {
	return "203.0.113." + strconv.Itoa(i+1)
}

func newObjectName() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
