package main

import (
	"AgentToolGateway/server/gateway"
	"AgentToolGateway/server/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := gateway.NewToolRepo()
	repo.Register(gateway.Tool{
		Name:    "read_note",
		Execute: tools.ReadNote,
	})
	server := gateway.NewGatewayServer(repo)

	r := gin.Default()
	r.POST("/toolCalls",
		func(c *gin.Context) {
			var call gateway.ToolCall
			if err := c.BindJSON(&call); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			result, err := server.Execute(call)
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
