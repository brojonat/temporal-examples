package temporal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RunPollCompleteWebhook(ctx context.Context, endpoint string, pr PollResult) error {
	b, err := json.Marshal(pr)
	if err != nil {
		return nil
	}
	r, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("bad response (%d) and error reading body: %w", res.StatusCode, err)
	}
	return fmt.Errorf("bad response (%d) and error: %s", res.StatusCode, b)
}
