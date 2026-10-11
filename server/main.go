package main

import (
	"AgentToolGateway/server/gateway"
	"AgentToolGateway/server/tools"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := gateway.NewToolRepo()
	am := gateway.NewApprovalManager()
	p := gateway.NewPolicy()
	audit := gateway.NewAuditLogger()
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
	repo.Register(gateway.Tool{
		Name:    "delete_note",
		Risk:    gateway.RiskLevelMedium,
		Execute: tools.DeleteNote,
	})
	server := gateway.NewGatewayServer(repo, am, p, audit)

	r := gin.Default()
	r.POST("/toolCalls",
		func(c *gin.Context) {
			var call gateway.ToolCall
			if err := c.ShouldBindJSON(&call); err != nil {
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
		err := am.Approve(id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err := server.ExecuteApproval(id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"result": result,
		})
	})
	r.GET("/logs", func(c *gin.Context) {
		events := audit.List()
		c.JSON(200, gin.H{
			"events": events,
		})

	})
	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
