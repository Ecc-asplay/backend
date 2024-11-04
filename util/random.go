package util

import (
	"fmt"
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

func RandomSexual() string {
	genders := []string{"男性", "女性", "その他"}
	return genders[seededRand.Intn(len(genders))]
}

func RandomCheckCode() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func BirthStringToInt(Bday string) (map[string]int, error) {
	BirthInt := map[string]int{
		"year":  0,
		"month": 0,
		"day":   0,
	}

	date, err := time.Parse("2006-01-02", Bday)
	if err != nil {
		fmt.Println("日付の解析エラー:", err)
		return nil, err
	}

	BirthInt["year"] = date.Year()
	BirthInt["month"] = int(date.Month())
	BirthInt["day"] = date.Day()

	return BirthInt, nil
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

func RandomGender() string {
	gender := []string{"男性", "女性", "その他"}
	return gender[rand.Intn(len(gender))]
}
