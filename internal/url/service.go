package url

import (
	"errors"
	"math/rand/v2"
	"net/url"
)

type urlService struct{
	store *urlStore
}

func NewUrlService(store *urlStore) *urlService{
	return &urlService{store: store}
}


func (us *urlService) ShortenUrl(originalUrl string) (sanitizedUrl, shortCode string, err error) {
	sanitizedUrl, err = us.normalizeUrl(originalUrl)
	if err != nil {
		return "", "", err
	}

	maxRetries := 5
	for range maxRetries{
		shortCode = us.generateUrl()

		err = us.store.SetUrl(sanitizedUrl, shortCode)
		if err != nil {
			switch {
			case errors.Is(err, ErrUrlDuplicated):
				continue
			default:
				return "", "", err
			}
		}

		return sanitizedUrl, shortCode, nil
	}

	return "", "", ErrShortUrlFailedGeneration
}

func (us *urlService) ResolveUrl(short_code string) (string, error) {
	url, err := us.store.GetUrl(short_code)
	if err != nil {
		return "", err
	}
	return url.OriginalUrl, nil
}

func (us *urlService) normalizeUrl(ori_url string) (string, error) {
	u, err := url.Parse(ori_url)
	if err != nil {
		return "", ErrUrlInvalid
	}

	if u.Scheme == "" {
		ori_url = "https://" + ori_url
		u, err = url.Parse(ori_url)
		if err != nil {
			return "", ErrUrlInvalid
		}
	}
	if !(u.Scheme == "https" || u.Scheme == "http") {
		return "", ErrUrlInvalid
	}

	if u.Host == "" {
		return "", ErrUrlInvalid
	}
	return ori_url, nil
}

func (us *urlService) generateUrl() string {
	lenShortCode := 7
	url := generator(lenShortCode)
	return url
}

func generator(digitString int) string {
	alphaNumeric := []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	tempURL := make([]byte, digitString)
	for i := range tempURL {
		tempURL[i] = alphaNumeric[rand.IntN(len(alphaNumeric))]
	}

	return string(tempURL)
}