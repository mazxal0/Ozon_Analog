package main

import (
	"Market_backend/internal/common"
	"Market_backend/internal/common/validate"
	"Market_backend/internal/config"
	"Market_backend/internal/server"
)

func main() {
	config.Init()
	common.InitDB()
	validate.InitValidator()
	server.Start()
}
