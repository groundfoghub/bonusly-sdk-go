package bonusly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type achivementResponse struct {
	Success bool          `json:"success"`
	Results []Achievement `json:"result"`
}

type Achievement struct {
	Id         string              `json:"id"`
	Headline   string              `json:"headline"`
	Title      string              `json:"title"`
	Importance float32             `json:"importance"`
	BonusId    string              `json:"bonus_id"`
	Scope      AchievementScope    `json:"scope"`
	Receiver   AchievementReceiver `json:"receiver"`
}

type AchievementScope struct {
	Department string `json:"department"`
}

type AchievementReceiver struct {
	Id          int    `json:"id"`
	DisplayName string `json:"display_name"`
	UserName    string `json:"username"`
	Email       string `json:"email"`
}

func (client *Client) GetAchivements() ([]Achievement, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/achievements", client.endpoint), nil)
	resp, requestErr := client.Do(req)

	if requestErr != nil {
		return make([]Achievement, 0), requestErr
	}

	var closingErr error

	defer func() {
		closingErr = resp.Body.Close()
	}()

	responseBody, bodyReadingErr := ioutil.ReadAll(resp.Body)

	if bodyReadingErr != nil {
		return make([]Achievement, 0), bodyReadingErr
	}

	var archivements achivementResponse
	jsonParsingErr := json.Unmarshal([]byte(responseBody), &archivements)

	if jsonParsingErr != nil {
		return make([]Achievement, 0), jsonParsingErr
	}

	if closingErr != nil {
		return archivements.Results, closingErr
	}

	return archivements.Results, nil
}
