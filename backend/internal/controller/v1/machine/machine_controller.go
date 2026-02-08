package machine

import (
	"fmt"
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceAllocation "github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	serviceOps "github.com/YoungBoyGod/remotegpu/internal/service/ops"
	"github.com/gin-gonic/gin"
)

type MachineController struct {
	common.BaseController
	machineService    *serviceMachine.MachineService
	allocationService *serviceAllocation.AllocationService
	agentService      *serviceOps.AgentService
}

func NewMachineController(ms *serviceMachine.MachineService, as *serviceAllocation.AllocationService, agentSvc *serviceOps.AgentService) *MachineController {
	return &MachineController{
		machineService:    ms,
		allocationService: as,
		agentService:      agentSvc,
	}
}

// List 获取机器列表
// @Summary 获取机器列表
// @Description 获取所有机器的列表，支持分页和筛选
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "状态筛选 (idle, allocated, maintenance, offline)"
// @Param region query string false "区域筛选"
// @Param gpu_model query string false "GPU型号筛选"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines [get]
func (c *MachineController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}
	if deviceStatus := ctx.Query("device_status"); deviceStatus != "" {
		filters["device_status"] = deviceStatus
	}
	if allocationStatus := ctx.Query("allocation_status"); allocationStatus != "" {
		filters["allocation_status"] = allocationStatus
	}
	if region := ctx.Query("region"); region != "" {
		filters["region"] = region
	}
	if gpuModel := ctx.Query("gpu_model"); gpuModel != "" {
		filters["gpu_model"] = gpuModel
	}

	machines, total, err := c.machineService.ListMachines(ctx, page, pageSize, filters)
	if err != nil {
		c.Error(ctx, 500, "Failed to list machines")
		return
	}

	c.Success(ctx, gin.H{
		"list":      machines,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Detail 获取机器详情
// @Summary 获取机器详情
// @Description 获取单个机器的详细信息
// @Tags Admin - Machines
// @Produce json
// @Param id path string true "机器ID"
// @Security Bearer
// @Success 200 {object} entity.Host
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/machines/{id} [get]
func (c *MachineController) Detail(ctx *gin.Context) {
	hostID := ctx.Param("id")
	detail, err := c.machineService.GetMachineDetail(ctx, hostID)
	if err != nil {
		c.Error(ctx, 404, "Machine not found")
		return
	}
	c.Success(ctx, detail)
}

// Create 创建机器
// @Summary 创建机器
// @Description 添加新的机器到系统
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param request body v1.CreateMachineRequest true "机器信息"
// @Security Bearer
// @Success 200 {object} entity.Host
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines [post]
func (c *MachineController) Create(ctx *gin.Context) {
	var req apiV1.CreateMachineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	if req.IPAddress == "" && req.Hostname == "" {
		c.Error(ctx, 400, "Host address is required")
		return
	}

	address := req.IPAddress
	if address == "" {
		address = req.Hostname
	}
	region := req.Region
	if region == "" {
		region = "default"
	}

	host := entity.Host{
		ID:           req.Hostname,
		Name:         req.Name,
		Hostname:     req.Hostname,
		Region:       region,
		IPAddress:    address,
		PublicIP:     req.PublicIP,
		SSHHost:      req.SSHHost,
		SSHPort:      req.SSHPort,
		SSHUsername:  req.SSHUsername,
		SSHPassword:  req.SSHPassword,
		SSHKey:       req.SSHKey,
		JupyterURL:   req.JupyterURL,
		JupyterToken: req.JupyterToken,
		VNCURL:       req.VNCURL,
		VNCPassword:  req.VNCPassword,
		ExternalIP:          req.ExternalIP,
		ExternalSSHPort:     req.ExternalSSHPort,
		ExternalJupyterPort: req.ExternalJupyterPort,
		ExternalVNCPort:     req.ExternalVNCPort,
		NginxDomain:         req.NginxDomain,
		NginxConfigPath:     req.NginxConfigPath,
		Status:           "offline",
		DeviceStatus:     "offline",
		AllocationStatus: "idle",
	}
	if host.ID == "" {
		host.ID = address
	}

	if c.agentService != nil {
		if info, err := c.agentService.GetSystemInfo(ctx, host.ID, address); err == nil {
			applySystemInfo(&host, info)
		}
	}

	if err := c.machineService.CreateMachine(ctx, &host); err != nil {
		if err == serviceMachine.ErrHostDuplicateIP || err == serviceMachine.ErrHostDuplicateHostname {
			c.Error(ctx, 409, "Host already exists")
			return
		}
		c.Error(ctx, 500, "Failed to create machine")
		return
	}

	c.Success(ctx, host)
}

// CollectSpec 触发采集机器硬件信息
// @Summary 采集机器硬件信息
// @Description 触发 Agent 采集并更新硬件信息
// @Tags Admin - Machines
// @Param id path string true "机器ID"
// @Security Bearer
// @Success 200 {object} entity.Host
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id}/collect [post]
func (c *MachineController) CollectSpec(ctx *gin.Context) {
	hostID := ctx.Param("id")
	if hostID == "" {
		c.Error(ctx, 400, "Host ID is required")
		return
	}
	if c.agentService == nil {
		c.Error(ctx, 500, "Agent service not available")
		return
	}

	host, err := c.machineService.GetHost(ctx, hostID)
	if err != nil {
		c.Error(ctx, 404, "Host not found")
		return
	}

	info, err := c.agentService.GetSystemInfo(ctx, host.ID, host.IPAddress)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	if err := c.machineService.CollectHostSpec(ctx, host, &serviceMachine.SystemInfoSnapshot{
		Hostname:      info.Hostname,
		CPUCores:      info.CPUCores,
		MemoryTotalGB: info.MemoryTotalGB,
		DiskTotalGB:   info.DiskTotalGB,
		Collected:     info.Collected,
	}); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	c.Success(ctx, host)
}

func applySystemInfo(host *entity.Host, info *serviceOps.SystemInfoSnapshot) {
	if info == nil {
		return
	}
	if host.Hostname == "" {
		host.Hostname = info.Hostname
	}
	if host.Name == "" {
		host.Name = info.Hostname
	}
	if info.OSType != "" {
		host.OSType = info.OSType
	}
	if info.Kernel != "" {
		host.OSVersion = info.Kernel
	}
	if info.CPUCores > 0 {
		host.TotalCPU = info.CPUCores
		host.CPUInfo = fmt.Sprintf("%d cores", info.CPUCores)
	}
	if info.MemoryTotalGB > 0 {
		host.TotalMemoryGB = info.MemoryTotalGB
	}
	if info.DiskTotalGB > 0 {
		host.TotalDiskGB = info.DiskTotalGB
	}
	if info.Collected {
		host.Status = "idle"
		host.DeviceStatus = "online"
		host.AllocationStatus = "idle"
		host.HealthStatus = "healthy"
	}
}

// Allocate 分配机器
// @Summary 分配机器
// @Description 将机器分配给客户
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param id path string true "机器ID"
// @Param request body v1.AllocateRequest true "分配请求"
// @Security Bearer
// @Success 200 {object} entity.Allocation
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id}/allocate [post]
func (c *MachineController) Allocate(ctx *gin.Context) {
	// machineID from URL param
	hostID := ctx.Param("id")

	var req apiV1.AllocateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	// Override HostID from URL if needed, or validate consistency
	if req.HostID != "" && req.HostID != hostID {
		c.Error(ctx, 400, "Host ID mismatch")
		return
	}

	alloc, err := c.allocationService.AllocateMachine(ctx, req.CustomerID, hostID, req.DurationMonths, req.Remark)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	c.Success(ctx, alloc)
}

// Reclaim 回收机器
// @Summary 回收机器
// @Description 从客户处回收机器
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param id path string true "机器ID"
// @Param request body v1.ReclaimRequest false "回收请求"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id}/reclaim [post]
func (c *MachineController) Reclaim(ctx *gin.Context) {
	hostID := ctx.Param("id")

	// Optional: bind body for reason
	var req apiV1.ReclaimRequest
	ctx.ShouldBindJSON(&req)

	if err := c.allocationService.ReclaimMachine(ctx, hostID); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	c.Success(ctx, gin.H{"message": "Reclaim process started"})
}

// Import 批量导入机器
// @Summary 批量导入机器
// @Description 批量导入机器信息
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param request body v1.ImportMachineRequest true "导入请求"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/import [post]
func (c *MachineController) Import(ctx *gin.Context) {
	var req apiV1.ImportMachineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid import data")
		return
	}

	var hosts []entity.Host
	for _, m := range req.Machines {
		// 根据 gpu_count 创建 GPU 关联记录
		var gpus []entity.GPU
		for i := 0; i < m.GPUCount; i++ {
			gpus = append(gpus, entity.GPU{
				HostID: m.HostIP,
				Index:  i,
				Name:   m.GPUModel,
				Status: "available",
			})
		}

		hosts = append(hosts, entity.Host{
			ID:               m.HostIP,
			IPAddress:        m.HostIP,
			SSHPort:          m.SSHPort,
			Region:           m.Region,
			TotalCPU:         m.CPUCores,
			TotalMemoryGB:    int64(m.RAMSize),
			TotalDiskGB:      int64(m.DiskSize),
			Status:           "offline",
			DeviceStatus:     "offline",
			AllocationStatus: "idle",
			NeedsCollect:     true,
			GPUs:             gpus,
		})
	}

	if err := c.machineService.ImportMachines(ctx, hosts); err != nil {
		c.Error(ctx, 500, "Failed to import machines")
		return
	}

	c.Success(ctx, gin.H{
		"message": "Imported successfully",
		"count":   len(hosts),
	})
}

// Update 更新机器信息
// @Summary 更新机器信息
// @Description 根据机器 ID 更新机器的基本信息字段
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param id path string true "机器 ID"
// @Param request body v1.UpdateMachineRequest true "更新机器请求"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id} [put]
func (c *MachineController) Update(ctx *gin.Context) {
	hostID := ctx.Param("id")

	var req apiV1.UpdateMachineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	fields := make(map[string]interface{})
	if req.Name != "" {
		fields["name"] = req.Name
	}
	if req.Region != "" {
		fields["region"] = req.Region
	}
	if req.PublicIP != "" {
		fields["public_ip"] = req.PublicIP
	}
	if req.SSHHost != "" {
		fields["ssh_host"] = req.SSHHost
	}
	if req.SSHPort > 0 {
		fields["ssh_port"] = req.SSHPort
	}
	if req.SSHUsername != "" {
		fields["ssh_username"] = req.SSHUsername
	}
	if req.SSHPassword != "" {
		fields["ssh_password"] = req.SSHPassword
	}
	if req.SSHKey != "" {
		fields["ssh_key"] = req.SSHKey
	}
	if req.JupyterURL != "" {
		fields["jupyter_url"] = req.JupyterURL
	}
	if req.JupyterToken != "" {
		fields["jupyter_token"] = req.JupyterToken
	}
	if req.VNCURL != "" {
		fields["vnc_url"] = req.VNCURL
	}
	if req.VNCPassword != "" {
		fields["vnc_password"] = req.VNCPassword
	}
	// 外映射配置字段
	if req.ExternalIP != "" {
		fields["external_ip"] = req.ExternalIP
	}
	if req.ExternalSSHPort > 0 {
		fields["external_ssh_port"] = req.ExternalSSHPort
	}
	if req.ExternalJupyterPort > 0 {
		fields["external_jupyter_port"] = req.ExternalJupyterPort
	}
	if req.ExternalVNCPort > 0 {
		fields["external_vnc_port"] = req.ExternalVNCPort
	}
	if req.NginxDomain != "" {
		fields["nginx_domain"] = req.NginxDomain
	}
	if req.NginxConfigPath != "" {
		fields["nginx_config_path"] = req.NginxConfigPath
	}

	if len(fields) == 0 {
		c.Error(ctx, 400, "No fields to update")
		return
	}

	if err := c.machineService.UpdateMachine(ctx, hostID, fields); err != nil {
		c.Error(ctx, 500, "Failed to update machine")
		return
	}

	c.Success(ctx, gin.H{"message": "Machine updated"})
}

// Delete 删除机器
// @Summary 删除机器
// @Description 从系统中删除机器
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param id path string true "机器ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id} [delete]
func (c *MachineController) Delete(ctx *gin.Context) {
	hostID := ctx.Param("id")

	// 检查机器分配状态，已分配的机器不允许删除
	host, err := c.machineService.GetHost(ctx, hostID)
	if err != nil {
		c.Error(ctx, 404, "Machine not found")
		return
	}
	if host.AllocationStatus == "allocated" {
		c.Error(ctx, 400, "无法删除已分配的机器，请先回收后再删除")
		return
	}

	if err := c.machineService.DeleteMachine(ctx, hostID); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	c.Success(ctx, gin.H{"message": "Machine deleted"})
}

// SetMaintenance 设置机器维护状态
// @Summary 设置机器维护状态
// @Description 将机器设置为维护状态或取消维护状态
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param id path string true "机器ID"
// @Param request body map[string]bool true "维护状态"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id}/maintenance [post]
func (c *MachineController) SetMaintenance(ctx *gin.Context) {
	hostID := ctx.Param("id")

	var req struct {
		Maintenance bool `json:"maintenance"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if req.Maintenance {
		// 进入维护模式
		if err := c.machineService.UpdateAllocationStatus(ctx, hostID, "maintenance"); err != nil {
			c.Error(ctx, 500, err.Error())
			return
		}
	} else {
		// 取消维护：检查是否有活跃分配，有则恢复为 allocated，否则恢复为 idle
		restoreStatus, err := c.machineService.ResolvePostMaintenanceStatus(ctx, hostID)
		if err != nil {
			c.Error(ctx, 500, err.Error())
			return
		}
		if err := c.machineService.UpdateAllocationStatus(ctx, hostID, restoreStatus); err != nil {
			c.Error(ctx, 500, err.Error())
			return
		}
	}

	c.Success(ctx, gin.H{"message": "Status updated"})
}

// Usage 获取机器使用情况
// @Summary 获取机器使用情况
// @Description 获取指定机器的资源使用统计信息
// @Tags Admin - Machines
// @Produce json
// @Param id path string true "机器 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id}/usage [get]
func (c *MachineController) Usage(ctx *gin.Context) {
	hostID := ctx.Param("id")
	usage, err := c.machineService.GetMachineUsage(ctx, hostID)
	if err != nil {
		c.Error(ctx, 404, "Machine not found")
		return
	}
	c.Success(ctx, usage)
}

// BatchSetMaintenance 批量启用/禁用机器（设置维护状态）
// @Summary 批量设置机器维护状态
// @Description 批量启用或禁用指定机器的维护模式
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param request body v1.BatchSetMaintenanceRequest true "批量维护请求"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/batch/maintenance [post]
func (c *MachineController) BatchSetMaintenance(ctx *gin.Context) {
	var req apiV1.BatchSetMaintenanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	affected, err := c.machineService.BatchSetMaintenance(ctx, req.HostIDs, req.Maintenance)
	if err != nil {
		c.Error(ctx, 500, "批量操作失败: "+err.Error())
		return
	}

	c.Success(ctx, gin.H{
		"affected": affected,
		"total":    len(req.HostIDs),
	})
}

// BatchAllocate 批量分配机器给客户
// @Summary 批量分配机器
// @Description 将多台机器批量分配给指定客户
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param request body v1.BatchAllocateRequest true "批量分配请求"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/batch/allocate [post]
func (c *MachineController) BatchAllocate(ctx *gin.Context) {
	var req apiV1.BatchAllocateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	result, err := c.allocationService.BatchAllocate(ctx, req.HostIDs, req.CustomerID, req.DurationMonths, req.Remark)
	if err != nil {
		c.Error(ctx, 500, "批量分配失败: "+err.Error())
		return
	}

	c.Success(ctx, result)
}

// BatchReclaim 批量回收机器
// @Summary 批量回收机器
// @Description 批量回收已分配给客户的机器
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param request body v1.BatchReclaimRequest true "批量回收请求"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/batch/reclaim [post]
func (c *MachineController) BatchReclaim(ctx *gin.Context) {
	var req apiV1.BatchReclaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	result, err := c.allocationService.BatchReclaim(ctx, req.HostIDs)
	if err != nil {
		c.Error(ctx, 500, "批量回收失败: "+err.Error())
		return
	}

	c.Success(ctx, result)
}
