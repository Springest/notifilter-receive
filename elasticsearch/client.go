package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ESPayload struct {
	Key        string                 `json:"key"`
	ReceivedAt time.Time              `json:"received_at"`
	Data       map[string]interface{} `json:"data"`
}

func Persist(key string, data map[string]interface{}) error {
	fmt.Printf("[ES] Persisting to %s: %v\n", key, data)

	payload := ESPayload{
		Key:        key,
		ReceivedAt: time.Now(),
		Data:       data,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:9200/notifilter/event/?pretty", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("[ES] Success")
	log.Println("[ES]", string(body))
	return nil
}

func EventCount() (int, error) {
	type response struct {
		Hits struct {
			Total int `json:"total"`
		} `json:"hits"`
	}

	resp, err := http.Get("http://localhost:9200/notifilter/event/_search?size=0")
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	fmt.Println("[ES] response: ", string(body))

	var parsed response
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return 0, err
	}

	fmt.Println(parsed)
	return parsed.Hits.Total, nil
}
