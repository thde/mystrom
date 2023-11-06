package mystrom_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"

	"thde.io/mystrom"
)

func TestSwitch_Toggle(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
			},
			want: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/toggle" {
					t.Errorf("expected /toggle path, got %s", r.URL.Path)
				}

				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			baseURL, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := mystrom.NewClient()
			s := client.NewSwitch(baseURL)

			err = s.Toggle(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Switch.Toggle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSwitch_Relay(t *testing.T) {
	t.Parallel()

	type args struct {
		state mystrom.RelaySwitchState
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "success on",
			args: args{
				state: mystrom.RelaySwitchStateOn,
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "success off",
			args: args{
				state: mystrom.RelaySwitchStateOff,
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				state: mystrom.RelaySwitchStateOff,
			},
			want: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/relay" {
					t.Errorf("expected /relay path, got %s", r.URL.Path)
				}

				values := r.URL.Query()
				if values.Get("state") != string(tt.args.state) {
					t.Errorf("expected state=%s query parameter, got %s", tt.args.state, values.Get("state"))
				}

				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			baseURL, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := mystrom.NewClient()
			s := client.NewSwitch(baseURL)

			err = s.Relay(context.Background(), tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("Switch.Relay() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSwitch_Report(t *testing.T) {
	t.Parallel()

	type args struct {
		body       []byte
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		want    *mystrom.SwitchReport
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				body:       []byte(`{"power": 100.0, "relay": true}`),
				statusCode: http.StatusOK,
			},
			want:    &mystrom.SwitchReport{Power: 100, Relay: true},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				statusCode: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/report" {
					t.Errorf("expected /report path, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.args.statusCode)
				_, _ = w.Write(tt.args.body)
			}))
			defer ts.Close()

			baseURL, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := mystrom.NewClient()
			s := client.NewSwitch(baseURL)

			report, err := s.Report(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Switch.Report() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(report, tt.want) {
				t.Errorf("Switch.Report() = %v, want %v", report, tt.want)
			}
		})
	}
}

func TestSwitch_Temperature(t *testing.T) {
	t.Parallel()

	type args struct {
		body       []byte
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		want    *mystrom.SwitchTemperature
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				body:       []byte(`{"measured": 22.5, "compensation": 0.5, "compensated": 23}`),
				statusCode: http.StatusOK,
			},
			want:    &mystrom.SwitchTemperature{Measured: 22.5, Compensation: 0.5, Compensated: 23},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				statusCode: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/api/v1/temperature" {
					t.Errorf("expected /api/v1/temperature path, got %s", r.URL.Path)
				}

				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(tt.args.body)
			}))
			defer ts.Close()

			baseURL, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := mystrom.NewClient()
			s := client.NewSwitch(baseURL)

			temp, err := s.Temperature(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Switch.Temperature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(temp, tt.want) {
				t.Errorf("Switch.Temperature() = %v, want %v", temp, tt.want)
			}
		})
	}
}

func TestSwitch_PowerCycle(t *testing.T) {
	t.Parallel()

	type args struct {
		wait       time.Duration
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				wait:       time.Second * 5,
				statusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				wait:       time.Second * 5,
				statusCode: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/power_cycle" {
					t.Errorf("expected /power_cycle path, got %s", r.URL.Path)
				}

				values := r.URL.Query()
				if values.Get("time") != strconv.Itoa(int(tt.args.wait.Seconds())) {
					t.Errorf("expected time=%d query parameter, got %s", int(tt.args.wait.Seconds()), values.Get("time"))
				}

				w.WriteHeader(tt.args.statusCode)
			}))
			defer ts.Close()

			baseURL, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := mystrom.NewClient()
			s := client.NewSwitch(baseURL)

			err = s.PowerCycle(context.Background(), tt.args.wait)
			if (err != nil) != tt.wantErr {
				t.Errorf("Switch.PowerCycle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSwitch_Timer(t *testing.T) {
	t.Parallel()

	type args struct {
		mode     mystrom.SwitchTimerMode
		duration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success on",
			args: args{
				mode:     mystrom.SwitchTimerModeOn,
				duration: time.Second * 5,
			},
			wantErr: false,
		},
		{
			name: "success off",
			args: args{
				mode:     mystrom.SwitchTimerModeOff,
				duration: time.Second * 5,
			},
			wantErr: false,
		},
		{
			name: "success toggle",
			args: args{
				mode:     mystrom.SwitchTimerModeToggle,
				duration: time.Second * 5,
			},
			wantErr: false,
		},
		{
			name: "success none",
			args: args{
				mode:     mystrom.SwitchTimerModeNone,
				duration: time.Second * 5,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				mode:     mystrom.SwitchTimerModeOn,
				duration: time.Second * 5,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/timer" {
					t.Errorf("expected /timer path, got %s", r.URL.Path)
				}

				values := r.URL.Query()
				if values.Get("time") != strconv.Itoa(int(tt.args.duration.Seconds())) {
					t.Errorf("expected time=%d query parameter, got %s", int(tt.args.duration.Seconds()), values.Get("time"))
				}

				if values.Get("mode") != string(tt.args.mode) {
					t.Errorf("expected mode=%s query parameter, got %s", tt.args.mode, values.Get("mode"))
				}

				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			baseURL, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := mystrom.NewClient()
			s := client.NewSwitch(baseURL)

			err = s.Timer(context.Background(), tt.args.mode, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Switch.Timer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
