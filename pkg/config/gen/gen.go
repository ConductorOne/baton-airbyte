package main

import (
	cfg "github.com/conductorone/baton-airbyte/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("airbyte", cfg.Config)
}
