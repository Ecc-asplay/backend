package util

import (
	"math/rand"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const charset = alphabet + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[seededRand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomPassword(len int) string {
	pd, _ := password.Generate(len, 10, 10, false, false)
	return pd
}

func RandomEra() int32 {
	return (seededRand.Int31n(9) + 1) * 10
}

func RandomCheckCode() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func SwitchAge(Y, M, D int) int32 {
	currentDate := time.Now()
	currentYear := currentDate.Year()
	currentMonth := int(currentDate.Month())
	currentDay := currentDate.Day()

	age := int32(currentYear - Y)

	if currentMonth < M || (currentMonth == M && currentDay < D) {
		age--
	}

	return age
}

func RandomBool() bool {
	bo := []bool{true, false}
	return bo[rand.Intn(len(bo))]
}

func RandomRole() string {
	bo := []string{"user", "admin"}
	return bo[rand.Intn(len(bo))]
}

func RandomGender() string {
	gender := []string{"男性", "女性", "その他"}
	return gender[rand.Intn(len(gender))]
}

func RandomDate() time.Time {
	rand.Seed(time.Now().UnixNano())
	year := rand.Intn(25) + 2000
	month := time.Month(rand.Intn(12) + 1)

	var day int
	if month == 2 {
		if (year%4 == 0 && year%100 != 0) || (year%400 == 0) {
			day = rand.Intn(29) + 1
		} else {
			day = rand.Intn(28) + 1
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		day = rand.Intn(30) + 1
	} else {
		day = rand.Intn(31) + 1
	}

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func RandomDisease() string {
	diseases := []string{
		"高血圧症",
		"糖尿病",
		"気管支喘息",
		"心筋梗塞",
		"肺炎",
		"認知症",
		"胃潰瘍",
		"アルツハイマー病",
		"花粉症",
		"うつ病",
	}

	rand.Seed(time.Now().UnixNano())
	return diseases[rand.Intn(len(diseases))]
}

func RandomCondition() string {
	conditions := []string{
		"安定している",
		"軽度の不快感",
		"中度の痛み",
		"重度の症状",
		"危険な状態",
	}

	rand.Seed(time.Now().UnixNano())
	return conditions[rand.Intn(len(conditions))]
}

func RandomStatus() string {
	statuses := []string{
		"処理中",
		"封鎖されました",
		"封鎖却下",
	}

	rand.Seed(time.Now().UnixNano())

	return statuses[rand.Intn(len(statuses))]
}

func RandomMood() string {
	moods := []string{
		"幸福",
		"悲しみ",
		"興奮",
		"安らぎ",
		"疲労",
		"怒り",
		"不安",
		"満足",
		"驚き",
		"混乱",
	}

	rand.Seed(time.Now().UnixNano())
	return moods[rand.Intn(len(moods))]
}
