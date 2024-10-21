package utils

import (
	"fmt"
	"os"
	"rr/domain"
)

func UrlCom(media []domain.Media, api string, vid string) {
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	for i := range media {
		media[i].Video = fmt.Sprintf("http://%s:%s/%s/%s/%s", ip, port, api, vid, media[i].Video)
	}
}
