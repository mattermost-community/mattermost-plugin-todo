package main

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type TelemetryAPIRequest struct {
	Event      string
	Properties map[string]interface{}
}

func GetTelemetryPayloadFromJSON(data io.Reader) (*TelemetryAPIRequest, error) {
	body := &TelemetryAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsTelemetryPayloadValid(t *TelemetryAPIRequest) error {
	if t == nil {
		return errors.New("invalid request body")
	}

	if t.Event == "" {
		return errors.New("event is required")
	}

	return nil
}

type AddAPIRequest struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	SendTo      string `json:"send_to"`
	PostID      string `json:"post_id"`
}

func GetAddIssuePayloadFromJSON(data io.Reader) (*AddAPIRequest, error) {
	body := &AddAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsAddIssuePayloadValid(a *AddAPIRequest) error {
	if a == nil {
		return errors.New("invalid request body")
	}

	if a.Message == "" {
		return errors.New("message is required")
	}

	return nil
}

type EditAPIRequest struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func GetEditIssuePayloadFromJSON(data io.Reader) (*EditAPIRequest, error) {
	body := &EditAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsEditIssuePayloadValid(e *EditAPIRequest) error {
	if e == nil {
		return errors.New("invalid request body")
	}

	if e.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type ChangeAssignmentAPIRequest struct {
	ID     string `json:"id"`
	SendTo string `json:"send_to"`
}

func GetChangeAssignmentPayloadFromJSON(data io.Reader) (*ChangeAssignmentAPIRequest, error) {
	body := &ChangeAssignmentAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsChangeAssignmentPayloadValid(c *ChangeAssignmentAPIRequest) error {
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

type AcceptAPIRequest struct {
	ID string `json:"id"`
}

func GetAcceptRequestPayloadFromJSON(data io.Reader) (*AcceptAPIRequest, error) {
	body := &AcceptAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsAcceptRequestPayloadValid(a *AcceptAPIRequest) error {
	if a == nil {
		return errors.New("invalid request body")
	}

	if a.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type CompleteAPIRequest struct {
	ID string `json:"id"`
}

func GetCompleteIssuePayloadFromJSON(data io.Reader) (*CompleteAPIRequest, error) {
	body := &CompleteAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsCompleteIssuePayloadValid(c *CompleteAPIRequest) error {
	if c == nil {
		return errors.New("invalid request body")
	}

	if c.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type RemoveAPIRequest struct {
	ID string `json:"id"`
}

func GetRemoveIssuePayloadFromJSON(data io.Reader) (*RemoveAPIRequest, error) {
	body := &RemoveAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsRemoveIssuePayloadValid(r *RemoveAPIRequest) error {
	if r == nil {
		return errors.New("invalid request body")
	}

	if r.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type BumpAPIRequest struct {
	ID string `json:"id"`
}

func GetBumpIssuePayloadFromJSON(data io.Reader) (*BumpAPIRequest, error) {
	body := &BumpAPIRequest{}
	if err := json.NewDecoder(data).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func IsBumpIssuePayloadValid(b *BumpAPIRequest) error {
	if b == nil {
		return errors.New("invalid request body")
	}

	if b.ID == "" {
		return errors.New("id is required")
	}

	return nil
}
