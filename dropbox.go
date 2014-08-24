package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	//	accountInfoUrl = "https://api.dropbox.com/1/account/info"
	listUrl  = "https://api.dropbox.com/1/metadata/auto"
	mediaUrl = "https://api.dropbox.com/1/media/auto"

	authHeader = "Bearer %s"
)

type DropBox struct {
	client      *http.Client
	accessToken string
}

func NewDropBox(accessToken string) *DropBox {
	dropBox := &DropBox{
		client:      &http.Client{},
		accessToken: accessToken,
	}
	return dropBox
}

type Image struct {
	Name string
	URL  string
}

func (d *DropBox) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf(authHeader, d.accessToken))

	res, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	//	fmt.Println(res)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return body, nil
}

func (d *DropBox) ReadImageList() ([]Image, error) {
	ld, err := d.Get(listUrl)
	if err != nil {
		return nil, err
	}

	var meta metadata
	json.Unmarshal(ld, &meta)

	//	fmt.Println(meta)

	// TODO: goroutineにするか・・・
	l := make([]Image, 0, len(meta.Contents))
	for _, content := range meta.Contents {
		if content.IsDir {
			continue
		}

		// TODO: 有効期限でキャッシュしたいお LevelDB?
		ad, err := d.Get(mediaUrl + content.Path)
		if err != nil {
			return nil, err
		}

		var m media
		json.Unmarshal(ad, &m)
		image := Image{
			Name: content.Path,
			URL:  m.URL,
		}
		l = append(l, image)
	}
	return l, nil
}
