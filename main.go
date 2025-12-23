package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TheRootDaemon/watchman/alert"
	"github.com/TheRootDaemon/watchman/detector"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load environment: %v", err)
	}

	hook := os.Getenv("DISCORD_HOOK")
	if hook == "" {
		log.Fatal("failed to get discord hook")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targets := []detector.Prober{
		detector.HTTPTarget{
			Name: "google",
			URL:  "https://google.com",

			Timeout: 10 * time.Second,
		},
		detector.HTTPTarget{
			Name:    "github",
			URL:     "https://github.com",
			Timeout: 10 * time.Second,
		},
		detector.TCPTarget{
			Name:    "ssh (localhost)",
			Address: "http://localhost:22",
			Timeout: 10 * time.Second,
		},
	}

	reportCh := make(chan detector.Report, 100)
	for _, t := range targets {
		go detector.RunProbe(ctx, t, reportCh)
	}

	fmt.Println("Watchman on duty...")
	fmt.Printf("%-20s | %-8s | %-10s | %s\n", "TARGET", "STATUS", "LATENCY", "ERROR")
	fmt.Println("---------------------------------------------------------------------")
	for {
		select {

		case <-ctx.Done():
			fmt.Println("Monitoring session ended.")

			return

		case r := <-reportCh:
			status := "UP"

			if !r.Success {
				status = "DOWN"
			}

			errStr := ""
			if r.Error != nil {
				errStr = r.Error.Error()
			}

			fmt.Printf(
				"%-20s | %-8s | %-10v | %s\n",
				r.Target,
				status,

				r.Latency.Round(time.Millisecond),
				errStr,
			)

			if !r.Success {
				go func(rep detector.Report) {
					if err := alert.SendReport(hook, rep); err != nil {
						log.Println("failed to send discord alert:", err)
					}
				}(r)
			}
		}
	}
}
