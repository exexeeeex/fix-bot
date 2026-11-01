package maxauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	testEnabled bool
}

func NewService() *Service {
	return &Service{
		testEnabled: os.Getenv("APP_ENV") == "development",
	}
}

func (s *Service) Authenticate(data string) (MaxUser, error) {
	if data == "" {
		return MaxUser{}, errors.New("data is required")
	}

	params, err := s.ValidateInitData(data)
	if err != nil {
		return MaxUser{}, err
	}

	authDate, err := strconv.ParseInt(
		params.Get("auth_date"),
		10,
		64,
	)
	if err != nil || authDate <= 0 {
		return MaxUser{}, errors.New("invalid auth_date")
	}

	createdAt := time.Unix(authDate, 0)
	now := time.Now()

	if createdAt.After(now.Add(time.Minute)) {
		return MaxUser{}, errors.New("auth_date is not the future")
	}

	if !s.testEnabled {
		if now.Sub(createdAt) > time.Hour {
			return MaxUser{}, errors.New("init data expired")
		}
	}

	userJSON := params.Get("user")
	if userJSON == "" {
		return MaxUser{}, errors.New("user not found")
	}

	var user MaxUser

	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return MaxUser{}, errors.New("invalid user data")
	}

	if user.Id <= 0 {
		return MaxUser{}, errors.New("invalid user id")
	}

	return user, nil
}

func (s *Service) ValidateInitData(
	data string,
) (url.Values, error) {
	params, err := url.ParseQuery(data)
	if err != nil {
		return nil, errors.New("invalid query string")
	}

	for _, values := range params {
		if len(values) != 1 {
			return nil, errors.New("duplicate parameter")
		}
	}

	hashValues, ok := params["hash"]
	if !ok || len(hashValues) != 1 || hashValues[0] == "" {
		return nil, errors.New("hash is mising")
	}

	originalHash := hashValues[0]

	if _, err := hex.DecodeString(originalHash); err != nil {
		return nil, errors.New("Invalid hash format")
	}

	keys := make([]string, 0, len(params))

	for key := range params {
		if key != "hash" {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	lines := make([]string, 0, len(keys))

	for _, key := range keys {
		lines = append(lines, key+"="+params.Get(key))
	}

	secretMAC := hmac.New(sha256.New, []byte("WebAppData"))
	secretMAC.Write([]byte(os.Getenv("MAX_BOT_TOKEN")))
	secretKey := secretMAC.Sum(nil)

	dataMAC := hmac.New(sha256.New, secretKey)
	dataMAC.Write([]byte(strings.Join(lines, "\n")))
	expectedHash := dataMAC.Sum(nil)

	actualHash, err := hex.DecodeString(originalHash)
	if err != nil {
		return nil, errors.New("invalid hash")
	}

	if !hmac.Equal(expectedHash, actualHash) {
		return nil, errors.New("signature mismatch")
	}

	return params, nil
}

func (s *Service) GenerateTestInitialData(
	user MaxUser,
	createdAt time.Time,
) (string, error) {
	if os.Getenv("MAX_BOT_TOKEN") == "" {
		return "", errors.New("MAX_BOT_TOKEN is not configured")
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("marshal user: %w", err)
	}

	params := url.Values{}
	params.Set("auth_date", strconv.FormatInt(createdAt.Unix(), 10))
	params.Set("user", string(userJSON))

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, key+"="+params.Get(key))
	}

	dataCheckString := strings.Join(lines, "\n")

	secretMAC := hmac.New(sha256.New, []byte("WebAppData"))
	_, _ = secretMAC.Write([]byte(os.Getenv("MAX_BOT_TOKEN")))
	secretKey := secretMAC.Sum(nil)

	dataMAC := hmac.New(sha256.New, secretKey)
	_, _ = dataMAC.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(dataMAC.Sum(nil))

	params.Set("hash", hash)

	return params.Encode(), nil
}
