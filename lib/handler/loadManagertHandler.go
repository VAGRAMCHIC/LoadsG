package handler

import (
	"github.com/gin-gonic/gin"

	"loadsg/lib/dto"
	"loadsg/lib/service"
)

type LoadManagerHandler struct {
	service service.LoadManagerService
}

func NewLoadManagerHandler(s service.LoadManagerService) *LoadManagerHandler {
	return &LoadManagerHandler{service: s}
}

func (h *LoadManagerHandler) CreateFixedHttp(c *gin.Context) {
	var req dto.CreateFixedHTTPLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(400)
		return
	}

	loadJob, err := h.service.CreateFixedHTTPLoadJob(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, gin.H{
		"job_id": loadJob.Id,
	})
}

func (h *LoadManagerHandler) CreateConstantHttp(c *gin.Context) {
	var req dto.CreateConstantHTTPLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(400)
		return
	}

	loadJob, err := h.service.CreateConstantHTTPLoadJob(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, gin.H{
		"job_id": loadJob.Id,
	})
}

func (h *LoadManagerHandler) CreateRampUpHttp(c *gin.Context) {
	var req dto.CreateRampUpHTTPLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(400)
		return
	}

	loadJob, err := h.service.CreateRampUpHTTPLoadJob(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, gin.H{
		"job_id": loadJob.Id,
	})
}

func (h *LoadManagerHandler) DeleteLoadJob(c *gin.Context) {
	Id := c.Param("id")
	if Id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
	}
	var resp dto.DeleteLoadJobResponse
	var err error
	resp.JobID, err = h.service.DeleteLoadJob(
		c.Request.Context(),
		Id,
	)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	c.JSON(204, resp)
}
