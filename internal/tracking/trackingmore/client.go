package trackingmore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const baseURL = "https://api.trackingmore.com"

type Event struct {
	Status     string
	Location   string
	Note       string
	OccurredAt time.Time
	RawDate    string
}

type Result struct {
	Status         string
	CarrierCode    string
	TrackingNumber string
	Events         []Event
}

type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(apiKey), "=")),
		http:   &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.apiKey != ""
}

type metaEnvelope struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
}

type createResponse struct {
	metaEnvelope
	Data json.RawMessage `json:"data"`
}

type trackingData struct {
	Status         string      `json:"status"`
	CarrierCode    string      `json:"carrier_code"`
	TrackingNumber string      `json:"tracking_number"`
	OriginInfo     *originInfo `json:"origin_info"`
}

type originInfo struct {
	Trackinfo []trackinfoRow `json:"trackinfo"`
}

type trackinfoRow struct {
	Date              string `json:"Date"`
	StatusDescription string `json:"StatusDescription"`
	Details           string `json:"Details"`
	CheckpointStatus  string `json:"checkpoint_status"`
}

type DetectedCarrier struct {
	Code string
	Name string
}

func (c *Client) DetectCarriers(ctx context.Context, trackingNumber string) ([]DetectedCarrier, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("trackingmore api key not configured")
	}
	body, err := json.Marshal(map[string]string{
		"tracking_number": strings.TrimSpace(trackingNumber),
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v2/carriers/detect", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Trackingmore-Api-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var env struct {
		metaEnvelope
		Data []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, err
	}
	if env.Meta.Code != 200 || len(env.Data) == 0 {
		return nil, fmt.Errorf("ไม่สามารถระบุขนส่งจากเลขพัสดุได้")
	}
	out := make([]DetectedCarrier, 0, len(env.Data))
	for _, row := range env.Data {
		code := strings.TrimSpace(row.Code)
		if code == "" {
			continue
		}
		out = append(out, DetectedCarrier{Code: code, Name: strings.TrimSpace(row.Name)})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("ไม่สามารถระบุขนส่งจากเลขพัสดุได้")
	}
	return out, nil
}

func (c *Client) CreateTracking(ctx context.Context, carrierCode, trackingNumber string) (*Result, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("trackingmore api key not configured")
	}
	body, err := json.Marshal(map[string]string{
		"tracking_number": strings.TrimSpace(trackingNumber),
		"carrier_code":    strings.TrimSpace(carrierCode),
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v2/trackings/post", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Trackingmore-Api-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var env createResponse
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, err
	}
	switch env.Meta.Code {
	case 200:
		var data trackingData
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return nil, err
		}
		return parseTrackingData(&data), nil
	case 4016:
		return c.GetTracking(ctx, carrierCode, trackingNumber)
	case 4032:
		return nil, fmt.Errorf("ไม่พบเลขพัสดุ กรุณาตรวจสอบเลขอีกครั้ง")
	case 401, 4001, 4002:
		return nil, fmt.Errorf("TrackingMore authentication failed")
	default:
		if env.Meta.Message != "" {
			return nil, fmt.Errorf("%s", env.Meta.Message)
		}
		return nil, fmt.Errorf("trackingmore error code %d", env.Meta.Code)
	}
}

func (c *Client) GetTracking(ctx context.Context, carrierCode, trackingNumber string) (*Result, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("trackingmore api key not configured")
	}
	carrierCode = strings.TrimSpace(carrierCode)
	trackingNumber = strings.TrimSpace(trackingNumber)
	url := fmt.Sprintf("%s/v2/trackings/%s/%s", baseURL, carrierCode, trackingNumber)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Trackingmore-Api-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var env struct {
		metaEnvelope
		Data trackingData `json:"data"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, err
	}
	if env.Meta.Code != 200 {
		if env.Meta.Message != "" {
			return nil, fmt.Errorf("%s", env.Meta.Message)
		}
		return nil, fmt.Errorf("trackingmore error code %d", env.Meta.Code)
	}
	return parseTrackingData(&env.Data), nil
}

func parseTrackingData(data *trackingData) *Result {
	if data == nil {
		return &Result{Status: "pending"}
	}
	out := &Result{
		Status:         normalizeStatus(data.Status),
		CarrierCode:    strings.TrimSpace(data.CarrierCode),
		TrackingNumber: strings.TrimSpace(data.TrackingNumber),
	}
	if data.OriginInfo != nil {
		for _, row := range data.OriginInfo.Trackinfo {
			ev := Event{
				Status:     firstNonEmpty(normalizeStatus(row.CheckpointStatus), normalizeStatus(data.Status)),
				Location:   strings.TrimSpace(row.Details),
				Note:       strings.TrimSpace(row.StatusDescription),
				RawDate:    strings.TrimSpace(row.Date),
				OccurredAt: parseTrackDate(row.Date),
			}
			out.Events = append(out.Events, ev)
		}
	}
	return out
}

func normalizeStatus(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return "pending"
}

func parseTrackDate(raw string) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t
		}
	}
	return time.Time{}
}
