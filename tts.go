package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	TextTemplateName = "ssml.xml"
)

var (
	Endpoint  = ""
	SecretKey = ""

	TextTemplate *template.Template
)

func init() {
	Endpoint = os.Getenv("ENDPOINT")
	SecretKey = os.Getenv("SECRET_KEY")
	if Endpoint == "" || SecretKey == "" {
		log.Fatalln("Env need ENDPOINT AND SECRET_KEY")
	}
	TextTemplate = template.Must(template.ParseFiles(TextTemplateName))
}

type VoiceParams struct {
	Lang        string `json:"lang"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Voice       string `json:"voice"`
	Pitch       string `json:"pitch"`
	Rate        string `json:"rate"`
	Style       string `json:"style"`
	StyleDegree string `json:"style_degree"`
}

func Voice(params *VoiceParams) (string, error) {
	params.Voice = escapeXmlText(params.Voice)
	buf := bytes.NewBuffer([]byte{})
	if err := TextTemplate.ExecuteTemplate(buf, TextTemplateName, params); err != nil {
		return "", err
	}
	println(buf.String())
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/cognitiveservices/v1", Endpoint), buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Microsoft-OutputFormat", "audio-48khz-192kbitrate-mono-mp3")
	req.Header.Set("Ocp-Apim-Subscription-Key", SecretKey)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("User-Agent", "zhuyst")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err.Error())
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status_code: %d", resp.StatusCode)
		println(err.Error())
		return "", err
	}

	fileName := uuid.NewV4().String() + ".mp3"
	file, err := os.OpenFile("./result/"+fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		println(err.Error())
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		println(err.Error())
		return "", err
	}
	return fileName, nil
}

type VoiceItem struct {
	Name            string   `json:"Name"`
	DisplayName     string   `json:"DisplayName"`
	LocalName       string   `json:"LocalName"`
	ShortName       string   `json:"ShortName"`
	Gender          string   `json:"Gender"`
	Locale          string   `json:"Locale"`
	LocaleName      string   `json:"LocaleName"`
	StyleList       []string `json:"StyleList"`
	SampleRateHertz string   `json:"SampleRateHertz"`
	VoiceType       string   `json:"VoiceType"`
	Status          string   `json:"Status"`
	RolePlayList    []string `json:"RolePlayList"`
	WordsPerMinute  string   `json:"WordsPerMinute"`
}

func MapList() (map[string]map[string]*VoiceItem, error) {
	items, err := List()
	if err != nil {
		return nil, err
	}

	m := make(map[string]map[string]*VoiceItem)
	for i := range items {
		item := items[i]
		_, ok := m[item.Locale]
		if !ok {
			m[item.Locale] = map[string]*VoiceItem{}
		}
		m[item.Locale][item.ShortName] = item
	}

	return m, nil
}

func List() ([]*VoiceItem, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/cognitiveservices/voices/list", Endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", SecretKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items []*VoiceItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

func escapeXmlText(text string) string {
	var escapeTextBuf bytes.Buffer
	xml.Escape(&escapeTextBuf, []byte(text))
	return escapeTextBuf.String()
}
