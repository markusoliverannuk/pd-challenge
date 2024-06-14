package server

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"gitlab.com/0x4149/logz"
)

type RateLimiter struct {
	mu        sync.Mutex
	lastCall  time.Time
	callDelay time.Duration
}

func NewRateLimiter(callDelay time.Duration) *RateLimiter {
	return &RateLimiter{
		lastCall:  time.Now(),
		callDelay: callDelay,
	}
}

func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// calculating time before next call
	waitTime := rl.callDelay - time.Since(rl.lastCall)
	if waitTime > 0 {
		time.Sleep(waitTime)
	}

	// updating last call
	rl.lastCall = time.Now()
}

// here we can conf the rate limiter, currently set to 125 ms, so 8 requests per sec or 16 per 2 sec (pipedrive limit for my plan is 20req/2sec)
var rateLimiter = NewRateLimiter(100 * time.Millisecond)

func CreatePipedriveDeal(username, description, id, originalID string) error {
	pipedriveAPIKey := os.Getenv("PIPEDRIVE_API_KEY")
	if pipedriveAPIKey == "" {
		return fmt.Errorf("PIPEDRIVE_API_KEY is not set")
	}

	rateLimiter.Wait()

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"title": fmt.Sprintf("Gist from %s", username),
			"fa317076d6d60293665ecdb667b151e748c5b512": fmt.Sprintf("https://gist.github.com/%s", username),
			"703e11f921b0a818e5c10aef045d520258d070bd": fmt.Sprintf("%s", description),
			"facb8d7e9ffc45bbb55b66e28d98c976a3b257b5": fmt.Sprintf("https://gist.github.com/%s", originalID),
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
