package util

import (
	"context"
	"net/http"
	"secure-chat-client/client"
	"secure-chat-client/env"
)

func Is2xx(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func CreateClient(accessToken *string) (*client.ClientWithResponses, error) {
	return client.NewClientWithResponses(env.Opts.WebServiceUrl, client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		if accessToken == nil {
			return nil
		}

		req.Header.Set("Authorization", "Bearer "+*accessToken)
		return nil
	}))
}
