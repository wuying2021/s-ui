package cluster

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Agent struct {
	id       string
	address  string
	master   string
	token    string
	interval time.Duration
	client   *http.Client
	cancel   context.CancelFunc
}

type registrationPayload struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Role     string `json:"role"`
	Hostname string `json:"hostname"`
}

func NewAgent(id, address, master, token string) *Agent {
	return &Agent{
		id:       id,
		address:  address,
		master:   master,
		token:    token,
		interval: time.Second * 30,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *Agent) Start() error {
	if a.master == "" {
		return errors.New("master endpoint is not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	go func() {
		ticker := time.NewTicker(a.interval)
		defer ticker.Stop()

		a.send("register")

		for {
			select {
			case <-ticker.C:
				a.send("heartbeat")
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (a *Agent) Stop() {
	if a.cancel != nil {
		a.cancel()
	}
}

func (a *Agent) send(action string) {
	body, err := json.Marshal(&registrationPayload{ID: a.id, Address: a.address, Role: "client", Hostname: a.id})
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/cluster/%s", a.master, action), bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if a.token != "" {
		req.Header.Set("X-Cluster-Token", a.token)
	}

	_, _ = a.client.Do(req)
}
