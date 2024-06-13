package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"gitlab.com/0x4149/logz"
)

func CreatePipedriveDeal(username, description, id string) error {
	pipedriveAPIKey := os.Getenv("PIPEDRIVE_API_KEY")
	if pipedriveAPIKey == "" {
		return fmt.Errorf("PIPEDRIVE_API_KEY is not set")
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"title": fmt.Sprintf("Gist from %s", username),
			"fa317076d6d60293665ecdb667b151e748c5b512": fmt.Sprintf("https://gist.github.com/%s", username),
			"703e11f921b0a818e5c10aef045d520258d070bd": fmt.Sprintf("%s", description),
		}).
		Post(fmt.Sprintf("https://api.pipedrive.com/v1/deals?api_token=%s", pipedriveAPIKey))

	if err != nil {
		return fmt.Errorf("error creating Pipedrive deal: %v", err)
	}

	if resp.StatusCode() != http.StatusCreated && resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("error creating Pipedrive deal, received status code %d: %s", resp.StatusCode(), resp.String())
	}

	logz.Info("Successfully created Pipedrive deal for gist:", description)
	return nil
}
