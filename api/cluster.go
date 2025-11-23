package api

import (
	"net/http"

	"github.com/alireza0/s-ui/cluster"
	"github.com/alireza0/s-ui/service"

	"github.com/gin-gonic/gin"
)

type ClusterAPI struct {
	clusterService *service.ClusterService
}

func NewClusterAPI(g *gin.RouterGroup, cs *service.ClusterService) {
	if cs == nil {
		return
	}
	h := &ClusterAPI{clusterService: cs}
	h.initRouter(g)
}

func (h *ClusterAPI) initRouter(g *gin.RouterGroup) {
	g.POST("/register", h.guard(h.register))
	g.POST("/heartbeat", h.guard(h.register))
	g.GET("/nodes", h.nodes)
}

func (h *ClusterAPI) guard(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Cluster-Token")
		if !h.clusterService.ValidateToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			c.Abort()
			return
		}
		next(c)
	}
}

func (h *ClusterAPI) register(c *gin.Context) {
	var node cluster.NodeInfo
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	node.Role = "client"
	node = h.clusterService.Register(node)
	c.JSON(http.StatusOK, node)
}

func (h *ClusterAPI) nodes(c *gin.Context) {
	nodes := h.clusterService.List()
	c.JSON(http.StatusOK, nodes)
}
