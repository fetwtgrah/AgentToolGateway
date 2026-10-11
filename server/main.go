package main

import (
	"AgentToolGateway/server/gateway"
	"AgentToolGateway/server/tools"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := gateway.NewToolRepo()
	am := gateway.NewApprovalManager()
	p := gateway.NewPolicy()
	repo.Register(gateway.Tool{
		Name:    "read_note",
		Risk:    gateway.RiskLevelLow,
		Execute: tools.ReadNote,
	})
	repo.Register(gateway.Tool{
		Name:    "send_email",
		Risk:    gateway.RiskLevelLow,
		Execute: tools.SendEmail,
	})
	repo.Register(gateway.Tool{
		Name:    "delete_note",
		Risk:    gateway.RiskLevelLow,
		Execute: tools.DeleteNote,
	})
	server := gateway.NewGatewayServer(repo, am, p)

	r := gin.Default()
	r.POST("/toolCalls",
		func(c *gin.Context) {
			var call gateway.ToolCall
			if err := c.ShouldBindJSON(&call); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			fmt.Println("tool:", call.ToolName)
			fmt.Println("args:", string(call.Args))
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
