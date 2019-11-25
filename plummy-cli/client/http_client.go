package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const paramsHeaderName = "x-plummy-params"

func NewHttpClient(baseURL string) Client {
	return &httpClient{
		c:       http.DefaultClient,
		baseURL: strings.TrimSuffix(baseURL, "/"), // Normalize base URL
	}
}

type httpClient struct {
	c       *http.Client
	baseURL string
}

func ToBody(fis []FileInput) []byte {
	b, _ := ioutil.ReadAll(newMultiplexEncoderReader(fis))
	return b
}

func (h httpClient) Render(ctx context.Context, engine string, req *RenderRequest) (*RenderResponse, error) {
	url := h.baseURL + "/v1/" + engine + "/render"
	hr, err := http.NewRequest(http.MethodPost, url, newMultiplexEncoderReader(req.Files))
	if err != nil {
		return nil, err
	}

	// Add context to allow timeout, etc
	hr = hr.WithContext(ctx)

	hr.Header.Set("Content-Type", multiplexContentType)
	hr.Header.Set(paramsHeaderName, base64.StdEncoding.EncodeToString(req.RawParams))

	resp, err := h.c.Do(hr)
	if err != nil {
		return nil, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if strings.Split(resp.Header.Get("Content-Type"), ";")[0] == "application/json" {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("cannot read http error response body: %w", err)
			}
			var errResp ErrorResponse
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				return nil, fmt.Errorf("cannot parse json error response: %w", err)
			}
			return nil, fmt.Errorf("%s [HTTP %d]: %s", errResp.Type, resp.StatusCode, errResp.Description)
		}
		return nil, fmt.Errorf("bad response status %d (%s)", resp.StatusCode, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != multiplexContentType {
		return nil, fmt.Errorf("unknown response content-type '%s'", contentType)
	}

	var rawParams []byte
	paramsEncoded := resp.Header.Get(paramsHeaderName)
	if paramsEncoded != "" {
		rawParams, err = base64.StdEncoding.DecodeString(paramsEncoded)
		if err != nil {
			return nil, fmt.Errorf("error decoding response params: %w", err)
		}
	}

	// TODO: Read files
	defer resp.Body.Close()
	files, err := readFiles(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response files: %w", err)
	}

	return &RenderResponse{
		RawParams: rawParams,
		Files: files,
	}, nil
}

func (h httpClient) checkHealth(ctx context.Context) error {
	hr, err := http.NewRequest(http.MethodGet, h.baseURL+"/healthz", nil)
	if err != nil {
		return err
	}

	// Here context would usually have a timeout or deadline.
	hr = hr.WithContext(ctx)

	// Call API
	resp, err := h.c.Do(hr)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code on daemon health API (%d)", resp.StatusCode)
	}

	return nil
}

func (h httpClient) WaitReady(ctx context.Context) error {
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		// One time check if no deadline is set.
		return h.checkHealth(ctx)
	}

	// Retry health check with a slight delay until deadline is reached
	var err error
	for time.Now().Before(deadline) {
		err = h.checkHealth(ctx)
		if err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	return err
}
