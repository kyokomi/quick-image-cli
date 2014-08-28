package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"bytes"
	"io"
	"os"
	"github.com/kyokomi/quick-image-cli/dropbox"
)

const (
	accountInfoUrl = "https://api.dropbox.com/1/account/info"

	listUrl = "https://api.dropbox.com/1/metadata/auto/Public"
	addUrl  = "https://api-content.dropbox.com/1/files_put/auto/Public"

	mediaUrl = "https://api.dropbox.com/1/media/auto"

	publicUrl  = "https://dl.dropbox.com/u/%.0f"
	authHeader = "Bearer %s"
)

type DropBox struct {
	Client      *http.Client
	AccessToken string
}

func NewDropBox(accessToken string) *DropBox {
	dropBox := &DropBox{
		Client:      &http.Client{},
		AccessToken: accessToken,
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
	req.Header.Set("Authorization", fmt.Sprintf(authHeader, d.AccessToken))
	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return body, nil
}

func (d *DropBox) PostFile(url_, filePath string) ([]byte, error) {

	var b bytes.Buffer
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(&b, f); err != nil {
		return nil, err
	}

	return d.Post(url_, b)
}

func (d *DropBox) Post(url_ string, buf bytes.Buffer) ([]byte, error) {

	req, err := http.NewRequest("POST", url_, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf(authHeader, d.AccessToken))

	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (d *DropBox) ReadImageList() ([]Image, error) {
	var meta dropbox.Metadata
	ld, err := d.Get(listUrl)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(ld, &meta)
	fmt.Println(string(ld))

	a, err := d.accountInfo()
	if err != nil {
		return nil, err
	}
	return d.readImageList(meta, a)
}

func (d *DropBox) readImageList(meta dropbox.Metadata, a dropbox.AccountInfo) ([]Image, error) {
	l := make([]Image, 0, len(meta.Contents))
	for _, content := range meta.Contents {
		if content.IsDir {
			continue
		}

		fileName := replacePublicFileName(content.Path)
		image := Image{
			Name: fileName,
			URL:  fmt.Sprintf(publicUrl, a.Uid) + fileName,
		}
		l = append(l, image)
	}
	return l, nil
}

func (d *DropBox) GetImage(contentPath string) ([]byte, error) {
	// send request
	ad, err := d.Post(mediaUrl+contentPath, bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (d *DropBox) AddImage(filePath string) (*Image, error) {

	a, err := d.accountInfo()
	if err != nil {
		return nil, err
	}

	url_ := createImageUrl(filePath)

	var p dropbox.FilePut
	pd, err := d.PostFile(url_, filePath)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(pd, &p)

	fileName := replacePublicFileName(p.Path)
	image := Image{
		Name: fileName,
		URL:  fmt.Sprintf(publicUrl, a.Uid) + fileName,
	}
	return &image, nil
}

func replacePublicFileName(filePath string) string {
	return strings.Replace(filePath, "/Public", "", 1)
}

func createImageUrl(filePath string) string {
	index := strings.LastIndex(filePath, "/")
	fileName := filePath[index+1:]
	return strings.Join([]string{addUrl, fileName}, "/")
}

func (d *DropBox) accountInfo() (dropbox.AccountInfo, error) {
	var a dropbox.AccountInfo
	info, err := d.Get(accountInfoUrl)
	if err != nil {
		return dropbox.AccountInfo{}, err
	}
	json.Unmarshal(info, &a)

	return a, nil
}
