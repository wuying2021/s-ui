package cluster

import (
	"sync"
	"time"
)

type NodeInfo struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	LastSeen  time.Time `json:"lastSeen"`
	Role      string    `json:"role"`
	Version   string    `json:"version"`
	Hostname  string    `json:"hostname"`
	Transport string    `json:"transport"`
}

type Manager struct {
	mu     sync.RWMutex
	nodes  map[string]NodeInfo
	token  string
	expiry time.Duration
}

func NewManager(token string) *Manager {
	return &Manager{
		nodes:  make(map[string]NodeInfo),
		token:  token,
		expiry: time.Minute * 5,
	}
}

func (m *Manager) UpdateNode(node NodeInfo) NodeInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	node.LastSeen = time.Now()
	m.nodes[node.ID] = node
	return node
}

func (m *Manager) List() []NodeInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, node := range m.nodes {
		if now.Sub(node.LastSeen) > m.expiry {
			delete(m.nodes, id)
		}
	}

	nodes := make([]NodeInfo, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (m *Manager) ValidateToken(token string) bool {
	if m.token == "" {
		return true
	}
	return token == m.token
}
