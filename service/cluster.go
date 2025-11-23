package service

import "github.com/alireza0/s-ui/cluster"

type ClusterService struct {
	manager *cluster.Manager
}

func NewClusterService(manager *cluster.Manager) *ClusterService {
	if manager == nil {
		return &ClusterService{}
	}
	return &ClusterService{manager: manager}
}

func (s *ClusterService) Register(node cluster.NodeInfo) cluster.NodeInfo {
	if s.manager == nil {
		return node
	}
	return s.manager.UpdateNode(node)
}

func (s *ClusterService) List() []cluster.NodeInfo {
	if s.manager == nil {
		return []cluster.NodeInfo{}
	}
	return s.manager.List()
}

func (s *ClusterService) ValidateToken(token string) bool {
	if s.manager == nil {
		return true
	}
	return s.manager.ValidateToken(token)
}
