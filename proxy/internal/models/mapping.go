package models

import "time"

// PortMapping 端口映射记录
type PortMapping struct {
	ID           string    `json:"id"`
	EnvID        string    `json:"env_id"`
	ServiceType  string    `json:"service_type"`
	ExternalPort int       `json:"external_port"`
	InternalPort int       `json:"internal_port"`
	TargetHost   string    `json:"target_host"`
	TargetPort   int       `json:"target_port"`
	Protocol     string    `json:"protocol"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// MappingRequest 创建映射请求
type MappingRequest struct {
	EnvID       string `json:"env_id" binding:"required"`
	ServiceType string `json:"service_type" binding:"required"`
	TargetHost  string `json:"target_host" binding:"required"`
	TargetPort  int    `json:"target_port" binding:"required"`
	Protocol    string `json:"protocol"`
}

// MappingResponse 映射响应
type MappingResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
