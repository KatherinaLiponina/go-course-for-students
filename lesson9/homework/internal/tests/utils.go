package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"strconv"
	"time"

	"homework9/internal/ports"
)

type adData struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type adResponse struct {
	Data adData `json:"data"`
}

type adsResponse struct {
	Data []adData `json:"data"`
}

type userData struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type userResponse struct {
	Data userData `json:"data"`
}

var (
	ErrBadRequest = fmt.Errorf("bad request")
	ErrForbidden  = fmt.Errorf("forbidden")
	ErrNotFound   = fmt.Errorf("not found")
)

type testClient struct {
	cancelFunc context.CancelFunc
	ch         chan int
	baseURL    string
}

func getTestClient() *testClient {
	ctx, cancelFunc := context.WithCancel(context.Background())
	endChan := make(chan int)
	server, _ := ports.CreateServer(ctx, endChan)
	testServer := httptest.NewServer(server.Handler)

	return &testClient{
		cancelFunc: cancelFunc,
		ch:         endChan,
		baseURL:    testServer.URL,
	}
}

func (t *testClient) cancelTestClient() {
	t.cancelFunc()
	<-t.ch
}

func (tc *testClient) getResponse(req *http.Request, out any) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return ErrBadRequest
		}
		if resp.StatusCode == http.StatusForbidden {
			return ErrForbidden
		}
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}

	err = json.Unmarshal(respBody, out)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	return nil
}

func (tc *testClient) createAd(userID int64, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) changeAdStatus(userID int64, adID int64, published bool) (adResponse, error) {
	body := map[string]any{
		"user_id":   userID,
		"published": published,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d/status", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateAd(userID int64, adID int64, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAds(r io.Reader) (adsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", r)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) createUser(nickname string, email string) (userResponse, error) {
	body := map[string]any{
		"nickname": nickname,
		"email":    email,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/users", bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateUser(id int64, nickname string, email string) (userResponse, error) {
	body := map[string]any{
		"nickname": nickname,
		"email":    email,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", id), bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAdsByAuthor(authorID int64) (adsResponse, error) {
	body := map[string]any{
		"by_author":     true,
		"author_id":     authorID,
		"by_creation":   false,
		"creation_time": nil,
		"all":           false,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	return tc.listAds(bytes.NewReader(data))
}

func (tc *testClient) listAdsByTime(time time.Time) (adsResponse, error) {
	body := map[string]any{
		"by_author":     false,
		"author_id":     -1,
		"by_creation":   true,
		"creation_time": time,
		"all":           false,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	return tc.listAds(bytes.NewReader(data))
}

func (tc *testClient) listAll() (adsResponse, error) {
	body := map[string]any{
		"by_author":     false,
		"author_id":     -1,
		"by_creation":   false,
		"creation_time": nil,
		"all":           true,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	return tc.listAds(bytes.NewReader(data))
}

func (tc *testClient) getAd(adID int64) (adResponse, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), nil)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) findAdByTitle(title string) (adsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads/title?title="+title, nil)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) GetUserByID(ID int64) (userResponse, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", ID), nil)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) DeleteUser(ID int64) (userResponse, error) {

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", ID), nil)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) DeleteAd(ID int64, authorID int64) (adResponse, error) {

	author := strconv.FormatInt(authorID, 10)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d?author=%s", ID, author), nil)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}
