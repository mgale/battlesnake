package main

import (
	"strings"
)

type Strategy interface {
	Initialize(GameRequest)
	GetMove(GameRequest) string
}

func createStrategy(strategyName string, gameRequest GameRequest) Strategy {
	if strings.ToLower(strategyName) == "solo" {
		solo := &Solo{}
		solo.Initialize(gameRequest)
		return solo
	} else {
		solo := &Solo{}
		solo.Initialize(gameRequest)
		return solo
	}
}
