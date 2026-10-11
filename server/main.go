package main

import (
	"AgentToolGateway/server/gateway"
	"AgentToolGateway/server/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := gateway.NewToolRepo()
	am := gateway.NewApprovalManager()
	repo.Register(gateway.Tool{
		Name:    "read_note",
		Risk:    gateway.RiskLevelLow,
		Execute: tools.ReadNote,
	})
	repo.Register(gateway.Tool{
		Name:    "send_email",
		Risk:    gateway.RiskLevelMedium,
		Execute: tools.SendEmail,
	})
	server := gateway.NewGatewayServer(repo, am)

	r := gin.Default()
	r.POST("/toolCalls",
		func(c *gin.Context) {
			var call gateway.ToolCall
			if err := c.BindJSON(&call); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			result, err := server.HandleToolCall(call)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{
				"result": result,
			})
		})
	r.POST("/approvals/:id", func(c *gin.Context) {
		id := c.Param("id")
		approval := am.Get(id)
		if approval == nil {
			return
		}
		approval.Status = "pending"
		result, err := server.ExecuteApproval(approval.Call)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"result": result,
		})
	})
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}
