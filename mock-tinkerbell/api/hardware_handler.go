package api

import (
	"net/http"
	"mock-tinkerbell/services"
	"github.com/gin-gonic/gin"
)

type HardwareHandler struct {
	hardwareService *services.HardwareService
}

func NewHardwareHandler(hardwareService *services.HardwareService) *HardwareHandler {
	return &HardwareHandler{
		hardwareService: hardwareService,
	}
}

func (h *HardwareHandler) GetHardware(c *gin.Context) {
	hardware := h.hardwareService.GetAllHardware()
	c.JSON(http.StatusOK, gin.H{
		"data": hardware,
		"total": len(hardware),
	})
}

func (h *HardwareHandler) GetHardwareByID(c *gin.Context) {
	id := c.Param("id")
	hardware, exists := h.hardwareService.GetHardwareByID(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hardware not found"})
		return
	}
	c.JSON(http.StatusOK, hardware)
}

func (h *HardwareHandler) UpdateHardwareStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.hardwareService.UpdateHardwareStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (h *HardwareHandler) GetHardwareStats(c *gin.Context) {
	stats := h.hardwareService.GetHardwareStats()
	c.JSON(http.StatusOK, stats)
}

func (h *HardwareHandler) SimulateDiscovery(c *gin.Context) {
	h.hardwareService.SimulateHardwareDiscovery()
	c.JSON(http.StatusOK, gin.H{"message": "Hardware discovery simulation completed"})
}
