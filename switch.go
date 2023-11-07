package mystrom

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Switch holds all info and logic to talk your myStrom Switch device.
type Switch struct {
	baseURL *url.URL
	client  *Client
}

// NewSwitch creates a new Switch instance.
func (c *Client) NewSwitch(baseURL *url.URL) *Switch {
	return &Switch{
		baseURL: baseURL,
		client:  c,
	}
}

// NewSwitch creates a new Switch instance with a default client.
func NewSwitch(baseURL *url.URL) *Switch {
	return NewClient().NewSwitch(baseURL)
}

// Toggle toggles the power state of the Switch.
func (s Switch) Toggle(ctx context.Context) error {
	req, err := s.client.newRequest(ctx, s.baseURL, http.MethodGet, "toggle", nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req)
	return err
}

type RelaySwitchState string

const (
	RelaySwitchStateOn  RelaySwitchState = "1"
	RelaySwitchStateOff RelaySwitchState = "0"
)

// On turns the power of the Switch on.
func (s Switch) On(ctx context.Context) error {
	return s.Relay(ctx, RelaySwitchStateOn)
}

// Off turns the power of the Switch off.
func (s Switch) Off(ctx context.Context) error {
	return s.Relay(ctx, RelaySwitchStateOff)
}

func (s Switch) Relay(ctx context.Context, state RelaySwitchState) error {
	req, err := s.client.newRequest(ctx, s.baseURL, http.MethodGet, "relay", url.Values{"state": []string{string(state)}}, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req)
	return err
}

// SwitchReport represets the content of a report of the Switch.
type SwitchReport struct {
	Power float64 `json:"power"` // current power consumption in watts
	Relay bool    `json:"relay"` // state of the Switch, true is on, false is off
}

// Report returns a report of the current statut of the Switch.
func (s Switch) Report(ctx context.Context) (*SwitchReport, error) {
	report := SwitchReport{}

	req, err := s.client.newRequest(ctx, s.baseURL, http.MethodGet, "report", nil, nil)
	if err != nil {
		return &report, err
	}

	_, err = s.client.doJSON(req, &report)
	return &report, err
}

// SwitchTemperature represets the content of a temperature response of
// the Switch. All temperatures are provided in °C.
type SwitchTemperature struct {
	Measured     float64 `json:"measured"`     // measured temp
	Compensation float64 `json:"compensation"` // assumed gap
	Compensated  float64 `json:"compensated"`  // real temp
}

// Temperature returns the current temperature in °C.
func (s Switch) Temperature(ctx context.Context) (*SwitchTemperature, error) {
	temp := SwitchTemperature{}

	req, err := s.client.newRequest(ctx, s.baseURL, http.MethodGet, "api/v1/temperature", nil, nil)
	if err != nil {
		return &temp, err
	}

	_, err = s.client.doJSON(req, &temp)
	return &temp, err
}

// PowerCycle turns the switch off, waits for a specified amount of time (max 1h), then starts it again.
// The switch has to be turned on in order for this call to work.
func (s Switch) PowerCycle(ctx context.Context, wait time.Duration) error {
	req, err := s.client.newRequest(
		ctx,
		s.baseURL,
		http.MethodGet,
		"power_cycle",
		url.Values{"time": []string{strconv.Itoa(int(wait.Seconds()))}},
		nil,
	)
	if err != nil {
		return err
	}

	_, err = s.client.do(req)
	return err
}

type SwitchTimerMode string

const (
	SwitchTimerModeOn     SwitchTimerMode = "on"
	SwitchTimerModeOff    SwitchTimerMode = "off"
	SwitchTimerModeToggle SwitchTimerMode = "toggle"
	SwitchTimerModeNone   SwitchTimerMode = "none"
)

// Timer sets the relay state for a given time, after the time has elapsed the state of the relay is reversed.
func (s Switch) Timer(ctx context.Context, mode SwitchTimerMode, time time.Duration) error {
	req, err := s.client.newRequest(
		ctx,
		s.baseURL,
		http.MethodPost,
		"timer",
		url.Values{
			"time": []string{strconv.Itoa(int(time.Seconds()))},
			"mode": []string{string(mode)},
		},
		nil,
	)
	if err != nil {
		return err
	}

	_, err = s.client.do(req)
	return err
}

func (s Switch) URL() url.URL {
	return *s.baseURL
}
