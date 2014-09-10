package main

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/kyokomi/quick-image-cli/dropbox"
	"strings"
	"net/http"
	"net/http/httptest"
)

func newStub(jsonPath string) (*httptest.Server, *DropBox) {
	stub, _ := ioutil.ReadFile(jsonPath)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(stub))
	}))
	d := NewDropBox("aaaaaaaaaaaa")
	fmt.Println(ts.URL)
	baseURL = ts.URL
	baseContentURL = ts.URL
	return ts, d
}

func TestNewDropBox(t *testing.T) {
	d := NewDropBox("aaaaaaaaaaaa")
	if d.AccessToken != "aaaaaaaaaaaa" {
		t.Error("not equal accessToken")
	}
}

func TestCreateFolder(t *testing.T) {
	ts, d := newStub("test/create_folder.json")

	res, err := d.CreateFolder("hoge_test")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(res))

	defer ts.Close()
}

func TestAddImage(t *testing.T) {
	ts, d := newStub("test/file_put.json")

	res, err := d.AddImage("", "test/gopher.png")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)

	defer ts.Close()
}

func TestMetaData(t *testing.T) {
	ts, d := newStub("test/meta.json")

	res, err := d.metaData("")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)

	defer ts.Close()
}

func TestAccountInfo(t *testing.T) {
	ts, d := newStub("test/account_info.json")

	res, err := d.accountInfo()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)

	defer ts.Close()
}

func TestReadImageList(t *testing.T) {

	var meta dropbox.Metadata
	metaData, err := ioutil.ReadFile("test/meta.json")
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(metaData, &meta)

	var a dropbox.AccountInfo
	aData, err := ioutil.ReadFile("test/account_info.json")
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(aData, &a)

	images, err := readImageList(&meta, &a, false)
	if err != nil {
		t.Error(err)
	}

	for _, image := range images {
		fmt.Println(image.Name)
	}

	if len(images) != 3 {
		t.Error("images len error")
	}
}

func TestReplacePublicFileName(t *testing.T) {
	fileName := replacePublicFileName("hogehoge/Public/fugafuga")
	if fileName != "hogehoge/fugafuga" {
		t.Errorf("replace error %s", fileName)
	}
}

func TestCreateImageUrl(t *testing.T) {
	imageUrl := createImageURL("", "/User/kyokomi/hoge/image.png")
	if imageUrl != strings.Join([]string{addURL, "", "image.png"}, "/") {
		t.Errorf("create image url error %s", imageUrl)
	}
}
