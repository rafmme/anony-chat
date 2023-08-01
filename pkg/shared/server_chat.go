package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
)

func ServerChat() string {
	apiUrls := []string{
		"https://api.api-ninjas.com/v1/jokes",
		"https://api.api-ninjas.com/v1/facts",
		"https://api.api-ninjas.com/v1/riddles",
		"https://api.api-ninjas.com/v1/dadjokes",
	}
	apiURL := apiUrls[GenerateRandomNumber(len(apiUrls))]

	r, err := http.NewRequest("GET", apiURL, nil)
	r.Header.Add("X-Api-Key", os.Getenv("BOT_KEY"))

	if err != nil {
		barf.Logger().Error(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		barf.Logger().Error(err.Error())
		return ""
	}

	defer resp.Body.Close()

	var result []map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		barf.Logger().Error(err.Error())
	}

	return formatResponse(result[0])
}

func formatResponse(result map[string]string) string {
	if result == nil {
		return ""
	}

	question := result["question"]
	if question != "" {
		return fmt.Sprintf("Riddle_Q: %s\nRiddle_A: %s", question, result["answer"])
	}

	if joke := result["joke"]; joke != "" {
		return joke
	}

	return result["fact"]
}
