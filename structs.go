package main

type Args struct {
	Path       string  `arg:"positional, required"`
	Age        float32 `arg:"-a" default:"-1"`
	Coins      float32 `arg:"-c" default:"-1"`
	Reputation int32   `arg:"-r" default:"-1"`
}
