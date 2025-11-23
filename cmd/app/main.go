package main

import (
	"eduVix_backend/internal/common"
	"eduVix_backend/internal/common/validate"
	"eduVix_backend/internal/config"
	"eduVix_backend/internal/server"
)

func main() {
	config.Init()
	common.InitDB()
	validate.InitValidator()
	server.Start()
}
