package main

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type telemetryAPIRequest struct {
	Event      string
	Properties map[string]interface{}
}

func GetTelemetryPayloadFromJSON(data io.Reader) (*telemetryAPIRequest, error) {
	var body *telemetryAPIRequest
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsTelemetryPayloadValid(t *telemetryAPIRequest) error {
	if t == nil {
		return errors.New("invalid request body")
	}

	if t.Event == "" {
		return errors.New("event is required")
	}

	return nil
}

type addAPIRequest struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	SendTo      string `json:"send_to"`
	PostID      string `json:"post_id"`
}

func GetAddIssuePayloadFromJSON(data io.Reader) (*addAPIRequest, error) {
	body := &addAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsAddIssuePayloadValid(a *addAPIRequest) error {
	if a == nil {
		return errors.New("invalid request body")
	}

	if a.Message == "" {
		return errors.New("message is required")
	}

	return nil
}

type editAPIRequest struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func GetEditIssuePayloadFromJSON(data io.Reader) (*editAPIRequest, error) {
	var body *editAPIRequest
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsEditIssuePayloadValid(e *editAPIRequest) error {
	if e == nil {
		return errors.New("invalid request body")
	}

	if e.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type changeAssignmentAPIRequest struct {
	ID     string `json:"id"`
	SendTo string `json:"send_to"`
}

func GetChangeAssignmentPayloadFromJSON(data io.Reader) (*changeAssignmentAPIRequest, error) {
	var body *changeAssignmentAPIRequest
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsChangeAssignmentPayloadValid(c *changeAssignmentAPIRequest) error {
	if c == nil {
		return errors.New("invalid request body")
	}

	if c.ID == "" {
		return errors.New("id is required")
	}

	if c.SendTo == "" {
		return errors.New("no user specified")
	}

	return nil
}

type acceptAPIRequest struct {
	ID string `json:"id"`
}

func GetAcceptRequestPayloadFromJSON(data io.Reader) (*acceptAPIRequest, error) {
	var body *acceptAPIRequest
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsAcceptRequestPayloadValid(a *acceptAPIRequest) error {
	if a == nil {
		return errors.New("invalid request body")
	}

	if a.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type completeAPIRequest struct {
	ID string `json:"id"`
}

func GetCompleteIssuePayloadFromJSON(data io.Reader) (*completeAPIRequest, error) {
	var body *completeAPIRequest
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsCompleteIssuePayloadValid(c *completeAPIRequest) error {
	if c == nil {
		return errors.New("invalid request body")
	}

	if c.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type removeAPIRequest struct {
	ID string `json:"id"`
}

func GetRemoveIssuePayloadFromJSON(data io.Reader) (*removeAPIRequest, error) {
	var body *removeAPIRequest
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsRemoveIssuePayloadValid(r *removeAPIRequest) error {
	if r == nil {
		return errors.New("invalid request body")
	}

	if r.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type bumpAPIRequest struct {
	ID string `json:"id"`
}

func GetBumpIssuePayloadFromJSON(data io.Reader) (*bumpAPIRequest, error) {
	var body *bumpAPIRequest
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsBumpIssuePayloadValid(b *bumpAPIRequest) error {
	if b == nil {
		return errors.New("invalid request body")
	}

	if b.ID == "" {
		return errors.New("id is required")
	}

	return nil
}
