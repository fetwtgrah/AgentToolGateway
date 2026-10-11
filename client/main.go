package main

import "AgentToolGateway/client/controller"

func main() {
	client := controller.NewClient("http://localhost:8080")
	call := map[string]any{
		"call_id":   "001",
		"tool_name": "send_email",
		"args": map[string]string{
			"to":      "test@qq.com",
			"subject": "test",
			"body":    "噢噢噢噢",
		},
	}
	if err := client.CallTool(call); err != nil {
		panic(err)
	}

}
