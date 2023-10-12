package whispervoice

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"encoding/json"

)

type TextRep struct {
	Text     string `json:"text"`
	Segments []struct {
		ID               int     `json:"id"`
		Seek             int     `json:"seek"`
		Start            float64 `json:"start"`
		End              float64 `json:"end"`
		Text             string  `json:"text"`
		Tokens           []int   `json:"tokens"`
		Temperature      float64 `json:"temperature"`
		AvgLogprob       float64 `json:"avg_logprob"`
		CompressionRatio float64 `json:"compression_ratio"`
		NoSpeechProb     float64 `json:"no_speech_prob"`
	} `json:"segments"`
	Language string `json:"language"`
}

func Voice2Text(file string,url string) (string, error) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	fw, err := writer.CreateFormFile("audio_file", filepath.Base(file))
	if err != nil {
		log.Fatal(err)
	}
	fd, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
		log.Fatal(err)
	}

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST",url, form)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s\n", bodyText)

	var resVoice TextRep

	jsonErr := json.Unmarshal(bodyText, &resVoice)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return "", fmt.Errorf("error parce json")
	}

	if resVoice.Text != "" {
		return resVoice.Text, nil
	} else {
		return "", fmt.Errorf("error from api")
	}


}