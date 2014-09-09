package main

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/kyokomi/quick-image-cli/dropbox"
	"strings"
)

func TestNewDropBox(t *testing.T) {
	d := NewDropBox("aaaaaaaaaaaa")
	if d.AccessToken != "aaaaaaaaaaaa" {
		t.Error("not equal accessToken")
	}
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

	images, err := readImageList(meta, a, false)
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
