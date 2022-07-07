package bonusly

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type ctrl struct {
	statusCode int
	response   string
}

func Test_gettingAchievements(t *testing.T) {
	server := httpMock("/api/v1/achievements", http.StatusOK, `{
		"success": true,
		"result" : [
			{
			  "id": "5b16f45e9fb5ba8225bc55ef",
			  "headline": "Max earned the best bonus tagged #help-out for the Engineering Department!",
			  "title": "best-bonus-received-tagged-#help-out-for-the-engineering-department",
			  "importance": 0.7,
			  "bonus_id": "12345abcde",
			  "scope": {
				"department": "engineering"
			  },
			  "receiver": {
				"id": 1,
				"display_name": "Max",
				"username": "max.mustermann",
				"email": "max.mustermann@example.com"
			  }
			}
		  ]
	}`)
	defer server.Close()
	achievementsArray, err := createMockClient(server.URL).GetAchivements()
	if len(achievementsArray) == 0 || err != nil {
		t.Errorf("GetAchivements should return an array of achivements = %v, want %v", achievementsArray, `{}`)
	}
}

func Test_gettingAchievementsWithInvalidToken(t *testing.T) {
	server := httpMock("/api/v1/achievements", http.StatusUnauthorized, `{}`)
	defer server.Close()
	achievementsArray, err := createMockClient(server.URL).GetAchivements()
	if len(achievementsArray) != 0 {
		t.Errorf("GetAchivements should return an Http-Unauthorized '401' but didn't")
	}
	if err != nil {
		t.Errorf("No error was thrown but was expected")
	}
}

func createMockClient(baseUrl string) *Client {
	mockClientConfig := Configuration{Token: "test"}
	return New(mockClientConfig, WithEndpoint(Endpoint(baseUrl+"/api/v1")))
}

func (c *ctrl) mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(c.statusCode)
	w.Write([]byte(c.response))
}

func httpMock(pattern string, statusCode int, response string) *httptest.Server {
	c := &ctrl{statusCode, response}
	handler := http.NewServeMux()
	handler.HandleFunc(pattern, c.mockHandler)
	return httptest.NewServer(handler)
}
