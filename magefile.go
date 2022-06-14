//go:build mage
// +build mage

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/nats-io/nats.go"
)

var (
	baseDir     = getMageDir()
	internalDir = fmt.Sprintf(filepath.Join(baseDir, "internal"))
	cmdDir      = fmt.Sprintf(filepath.Join(baseDir, "cmd"))
	apps        = []string{
		"points",
		"user",
		"payer",
	}
)

func getMageDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return dir
}

type Gen mg.Namespace

// generates open api boilerplate
func (Gen) Api() error {
	cfg := "oapi-cfg.yaml"
	yml := "openapi.json"

	for _, app := range apps {
		err := os.MkdirAll(filepath.Join(internalDir, app, "api"), os.ModePerm)
		if err != nil {
			return fmt.Errorf("Api: %s", err)
		}

		err = os.Chdir(filepath.Join(cmdDir, app))
		if err != nil {
			return fmt.Errorf("Api: %s", err)
		}

		err = sh.Run("oapi-codegen", "--config", cfg, yml)
		if err != nil {
			return fmt.Errorf("Api: %s", err)
		}

	}

	defer os.Chdir(baseDir)

	return nil
}

// runs unit tests
func Unit() error {
	coverage := "coverage.out"

	err := sh.Run("go", "test", "-race", "-v", "-tags=unit", "-covermode=atomic", fmt.Sprintf("-coverprofile=%s", coverage), filepath.Join(baseDir, "internal", "..."))
	if err != nil {
		return fmt.Errorf("Unit: %s", err)
	}

	return nil
}

// builds app binary
func Build() error {
	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
	}

	for _, app := range apps {
		err := os.Chdir(filepath.Join(baseDir, "cmd", app))
		if err != nil {
			return fmt.Errorf("Build: %s", err)
		}

		err = sh.RunWithV(env, "go", "build", "-o", "main")
		if err != nil {
			return fmt.Errorf("Build: %s", err)
		}
	}

	defer os.Chdir(baseDir)

	return nil
}

func Start() error {
	depsCompose := "docker-compose.deps.yaml"
	appsCompose := "docker-compose.apps.yaml"

	err := Build()
	if err != nil {
		return fmt.Errorf("Start: %s", err)
	}

	err = sh.Run("docker", "compose", "-f", filepath.Join(baseDir, depsCompose), "up", "-d")
	if err != nil {
		return fmt.Errorf("Start: %s", err)
	}

	err = setupNats()
	if err != nil {
		return fmt.Errorf("Start: %s", err)
	}

	err = sh.Run("docker", "compose", "-f", filepath.Join(baseDir, appsCompose), "up", "-d")
	if err != nil {
		return fmt.Errorf("Start: %s", err)
	}

	return nil
}

func setupNats() error {
	streams := []string{
		"PAYER",
		"USER",
	}

	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)

	// Creates JetStreamContext
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, stream := range streams {
		err := createStream(js, stream, fmt.Sprintf("%s.*", stream))
		if err != nil {
			return fmt.Errorf("setupNats: %s", err)
		}
	}

	return nil
}

// createStream creates a stream by using JetStreamContext
func createStream(js nats.JetStreamContext, name string, subjects string) error {
	stream, err := js.StreamInfo(name)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", name, subjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     name,
			Subjects: []string{subjects},
		})
		if err != nil {
			return err
		}

		// Create a Consumer
		js.AddConsumer(name, &nats.ConsumerConfig{
			Durable: "MONITOR",
		})
	}
	return nil
}
