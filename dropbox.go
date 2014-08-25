package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"code.google.com/p/leveldb-go/leveldb"
	"code.google.com/p/leveldb-go/leveldb/db"
	"io"
	"os"
	"bytes"
)

const (
	//	accountInfoUrl = "https://api.dropbox.com/1/account/info"
	listUrl  = "https://api.dropbox.com/1/metadata/auto"
	addUrl   = "https://api-content.dropbox.com/1/files_put/auto"
	mediaUrl = "https://api.dropbox.com/1/media/auto"

	authHeader = "Bearer %s"
)

type DropBox struct {
	client      *http.Client
	accessToken string
	level       *leveldb.DB
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

func (d *DropBox) SetupCache(cacheDirPath string) error {
	level, err := leveldb.Open(cacheDirPath, &db.Options{})
	if err != nil {
		return err
	}

	d.level = level
	return nil
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

func (d *DropBox) Post(imageUrl, filePath string) ([]byte, error) {

	var b bytes.Buffer
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(&b, f); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", imageUrl, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf(authHeader, d.accessToken))

	res, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (d *DropBox) ReadImageList() ([]Image, error) {
	ld, err := d.Get(listUrl)
	if err != nil {
		return nil, err
	}

	var meta metadata
	json.Unmarshal(ld, &meta)

	//	fmt.Println(meta)

	l := make([]Image, 0, len(meta.Contents))
	for _, content := range meta.Contents {
		if content.IsDir {
			continue
		}

		// TODO: ちゃんと有効期限でキャッシュしたいお
		ad, err := d.level.Get([]byte(content.Path), &db.ReadOptions{})
		if err != nil {
			// send request
			ad, err = d.Get(mediaUrl + content.Path)
			if err != nil {
				continue
			}
			// cache
			if err := d.level.Set([]byte(content.Path), ad, &db.WriteOptions{}); err != nil {
				continue
			}
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

func (d *DropBox) AddImage(filePath string) error {

	index := strings.LastIndex(filePath, "/")
	fileName := filePath[index+1:]
	url := strings.Join([]string{addUrl, fileName}, "/")

	pd, err := d.Post(url, filePath)
	if err != nil {
		return err
	}
	fmt.Println(string(pd))

	return nil
}
