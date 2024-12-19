package idgen

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

func GenerateRandomValue(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charList := strings.Split(charset, "")
	var randomCharList []string
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charList) - 1)
		randomCharList = append(randomCharList, charList[randomIndex])
	}
	return strings.Join(randomCharList, "")
}

func GenerateTimeStampId(prefix string) string {
	now := time.Now()
	timestamp := strconv.FormatInt(now.UnixMilli(), 10)
	return fmt.Sprintf("%s %s %s", prefix, timestamp, GenerateRandomValue(2))
}

func GenerateShortUUID(length int) string {
	rawUUID := uuid.New()
	base64Encoded := base64.StdEncoding.EncodeToString(rawUUID[:])
	base64Clean := strings.NewReplacer("+", "", "/", "", "=", "").Replace(base64Encoded)
	if length > len(base64Clean) {
		length = len(base64Clean)
	}
	return base64Clean[:length]
}

func GenerateID() func() string {
	var counter int64   // Atomic counter
	var lastDate string // Tracks the last date

	// Base62 character set
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Converts a number to Base62
	toBase62 := func(num int64) string {
		if num == 0 {
			return "0"
		}
		var encoded string
		for num > 0 {
			remainder := num % 62
			encoded = string(base62Chars[remainder]) + encoded
			num /= 62
		}
		return encoded
	}

	return func() string {
		// Current time
		now := time.Now()

		// Format the date as YYMMDD
		currentDate := now.Format("060102")

		// Reset counter at midnight if the date has changed
		if currentDate != lastDate {
			atomic.StoreInt64(&counter, 0)
			lastDate = currentDate
		}

		// Milliseconds since the start of the day
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		milliseconds := now.Sub(midnight).Milliseconds()

		// Atomically increment counter
		currentCounter := atomic.AddInt64(&counter, 1)

		// Convert milliseconds and counter to Base62
		millisecondsBase62 := toBase62(milliseconds)
		counterBase62 := toBase62(currentCounter)

		// Construct the ID
		return fmt.Sprintf("MER-%s-%s-%s", currentDate, millisecondsBase62, counterBase62)
	}
}
