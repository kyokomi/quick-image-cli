package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kyokomi/quick-image-cli/dropbox"
)

var (
	baseURL        = "https://api.dropbox.com/1"
	baseContentURL = "https://api-content.dropbox.com/1"
)

const (
	accountInfoURL  = "/account/info"
	listURL         = "/metadata/auto/Public"
	createFolderURL = "/fileops/create_folder"
	mediaURL        = "/media/auto"

	addURL          = "/files_put/auto/Public"

	publicURL  = "https://dl.dropbox.com/u/%.0f"
	authHeader = "Bearer %s"
)

type DropBox struct {
	Client      *http.Client
	AccessToken string
	BaseURL     string
}

func NewDropBox(accessToken string) *DropBox {
	dropBox := &DropBox{
		Client:      &http.Client{},
		AccessToken: accessToken,
	}
	return dropBox
}

func resourceContentURL(url string) string {
	return baseContentURL + url
}

func resourceURL(url string) string {
	return baseURL + url
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

func (d *DropBox) ReadImageList(path string, isDir bool) ([]Image, error) {

	meta, err := d.metaData(path)
	if err != nil {
		return nil, err
	}

	a, err := d.accountInfo()
	if err != nil {
		return nil, err
	}
	return readImageList(meta, a, isDir)
}

func (d *DropBox) metaData(path string) (*dropbox.Metadata, error) {

	url := strings.Join([]string{resourceURL(listURL), path}, "/")

	var meta dropbox.Metadata
	ld, err := d.Get(url)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(ld, &meta)

	return &meta, nil
}

func readImageList(meta *dropbox.Metadata, a *dropbox.AccountInfo, isDir bool) ([]Image, error) {
	l := make([]Image, 0, len(meta.Contents))
	for _, content := range meta.Contents {
		if !isDir && content.IsDir {
			continue
		}

		fileName := replacePublicFileName(content.Path)
		url := ""
		if !content.IsDir {
			url = fmt.Sprintf(publicURL, a.Uid) + fileName
		}

		image := Image{
			Name: fileName,
			URL:  url,
		}
		l = append(l, image)
	}
	return l, nil
}

func (d *DropBox) GetImage(contentPath string) ([]byte, error) {

	url := resourceURL(mediaURL) + contentPath

	// send request
	ad, err := d.Post(url, nil, nil)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (d *DropBox) AddImage(name, dirPath, filePath string) (*Image, error) {

	a, err := d.accountInfo()
	if err != nil {
		return nil, err
	}

	url := resourceContentURL(createImageURL(name, dirPath, filePath))

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

func (d *DropBox) accountInfo() (*dropbox.AccountInfo, error) {

	url := resourceURL(accountInfoURL)

	var a dropbox.AccountInfo
	info, err := d.Get(url)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(info, &a)

	return &a, nil
}

func (d *DropBox) CreateFolder(path string) ([]byte, error) {
	v := url.Values{
		"root": []string{"auto"},
		"path": []string{"Public/" + path},
	}
	params := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	url := resourceURL(createFolderURL)
	return d.Post(url, strings.NewReader(v.Encode()), params)
}

func replacePublicFileName(filePath string) string {
	return strings.Replace(filePath, "/Public", "", 1)
}

func createImageURL(name, dirPath, filePath string) string {

	if name == "" {
		index := strings.LastIndex(filePath, "/")
		name = filePath[index+1:]
	}
	return strings.Join([]string{addURL, dirPath, name}, "/")
}
