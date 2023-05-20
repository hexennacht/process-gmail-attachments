package main

import (
	"github.com/hexennacht/process-gmail-attachments/config"
	"github.com/hexennacht/process-gmail-attachments/server"
)

func main() {
	server.Serve(config.Read())
}
