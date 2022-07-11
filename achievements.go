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

	bodyMessage, readErr := readBody(resp)

	if readErr != nil {
		return make([]Achievement, 0), readErr
	}

	var archivements achivementResponse
	jsonParsingErr := json.Unmarshal([]byte(bodyMessage), &archivements)

	if jsonParsingErr != nil {
		return make([]Achievement, 0), jsonParsingErr
	}

	return archivements.Results, nil
}

func readBody(resp *http.Response) (body []byte, err error) {
	defer func() {
		if tempErr := close(resp); tempErr != nil {
			err = tempErr
		}
	}()
	return ioutil.ReadAll(resp.Body)
}

func close(resp *http.Response) error {
	return resp.Body.Close()
}
