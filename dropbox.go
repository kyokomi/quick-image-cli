package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/kyokomi/quick-image-cli/dropbox"
)

const (
	accountInfoURL = "https://api.dropbox.com/1/account/info"

	listURL = "https://api.dropbox.com/1/metadata/auto/Public"
	addURL  = "https://api-content.dropbox.com/1/files_put/auto/Public"

	mediaURL = "https://api.dropbox.com/1/media/auto"

	publicURL  = "https://dl.dropbox.com/u/%.0f"
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

func (d *DropBox) PostFile(url, filePath string) ([]byte, error) {
	var b bytes.Buffer
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(&b, f); err != nil {
		return nil, err
	}

	return d.Post(url, bytes.NewReader(b.Bytes()), nil)
}

func (d *DropBox) newPostRequest(url string, body io.Reader, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf(authHeader, d.AccessToken))
	for key, value := range params {
		req.Header.Set(key, value)
	}
	return req, nil
}

func (d *DropBox) Post(url string, body io.Reader, params map[string]string) ([]byte, error) {

	req, err := d.newPostRequest(url, body, params)
	if err != nil {
		return nil, err
	}

	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (d *DropBox) ReadImageList() ([]Image, error) {
	var meta dropbox.Metadata
	ld, err := d.Get(listURL)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(ld, &meta)

	a, err := d.accountInfo()
	if err != nil {
		return nil, err
	}
	return readImageList(meta, a, isDir)
}

func readImageList(meta dropbox.Metadata, a dropbox.AccountInfo, isDir bool) ([]Image, error) {
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
	ad, err := d.Post(mediaURL+contentPath, nil, nil)
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

	url := createImageURL(filePath)

	var p dropbox.FilePut
	pd, err := d.PostFile(url, filePath)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(pd, &p)

	fileName := replacePublicFileName(p.Path)
	image := Image{
		Name: fileName,
		URL:  fmt.Sprintf(publicURL, a.Uid) + fileName,
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
	info, err := d.Get(accountInfoURL)
	if err != nil {
		return dropbox.AccountInfo{}, err
	}
	json.Unmarshal(info, &a)

	return a, nil
}
