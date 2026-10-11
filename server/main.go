package main

import (
	"AgentToolGateway/server/gateway"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := gateway.NewToolRepo()
	repo.Register(gateway.Tool{
		Name: "read_note",
		//Execute: tools.ReadNote,
	})
	r := gin.Default()
	r.POST("/toolCalls",
		func(c *gin.Context) {
			var call gateway.ToolCall
			if err := c.BindJSON(&call); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{
				"tool": call.ToolName,
			})
		})
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}
