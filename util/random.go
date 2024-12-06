package util

import (
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const charset = alphabet + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(max int) int {
	if max <= 0 {
		log.Fatal("無効な最大値です: 0 より大きい必要があります")
	}
	return rand.Intn(max)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := RandomInt(k)
		sb.WriteByte(alphabet[c])
	}
	return sb.String()
}

func RandomPassword(len int) string {
	pd, _ := password.Generate(len, 10, 10, false, false)
	return pd
}

func RandomEra() int32 {
	num := RandomInt(9)
	return int32((num + 1) * 10)
}

func RandomCheckCode() string {
	b := make([]byte, 8)
	for i := range b {
		idx := RandomInt(len(charset))
		b[i] = charset[idx]
	}
	return string(b)
}

func SwitchAge(Y, M, D int) int32 {
	currentDate := time.Now()
	currentYear := currentDate.Year()
	currentMonth := int(currentDate.Month())
	currentDay := currentDate.Day()

	age := int64(currentYear - Y)
	if age < 0 || age > math.MaxInt32 {
		log.Fatal("無効な年齢計算です: オーバーフローまたは負の年齢")
	}

	if currentMonth < M || (currentMonth == M && currentDay < D) {
		age--
	}

	return int32(age)
}

func RandomBool() bool {
	bo := []bool{true, false}
	idx := RandomInt(len(bo))
	return bo[idx]
}

func RandomRole() string {
	bo := []string{"user", "admin"}
	idx := RandomInt(len(bo))
	return bo[idx]
}

func RandomGender() string {
	gender := []string{"男性", "女性", "その他"}
	idx := RandomInt(len(gender))
	return gender[idx]
}

func RandomDate() time.Time {
	year := RandomInt(25) + 2000 // Random year between 2000 and 2024
	month := RandomInt(12) + 1   // Random month between 1 and 12

	var day int
	if month == 2 {
		if (year%4 == 0 && year%100 != 0) || (year%400 == 0) {
			day = RandomInt(29)
		} else {
			day = RandomInt(28)
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		day = RandomInt(30)
	} else {
		day = RandomInt(31)
	}

	day += 1

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
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

	idx := RandomInt(len(diseases))
	return diseases[idx]
}

func RandomCondition() string {
	conditions := []string{
		"安定している",
		"軽度の不快感",
		"中度の痛み",
		"重度の症状",
		"危険な状態",
	}

	idx := RandomInt(len(conditions))
	return conditions[idx]
}

func RandomStatus() string {
	statuses := []string{
		"処理中",
		"封鎖されました",
		"封鎖却下",
	}

	idx := RandomInt(len(statuses))
	return statuses[idx]
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

	idx := RandomInt(len(moods))
	return moods[idx]
}
