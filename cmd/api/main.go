package main

import (
	"cmyk/internal/clients/env"
	"cmyk/internal/clients/k8s"
	"cmyk/internal/clients/mock"
	"cmyk/internal/clients/socks5"
	"cmyk/internal/handlers"

	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

var (
	port = flag.String("port", ":4000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	flag.Parse()

	envClient := env.New()
	log.Printf("created env client: %v", envClient)

	var socks5Client *socks5.Client
	if envClient.Socks5ProxyMode() {
		var err error
		socks5Client, err = socks5.New(envClient.Socks5ProxyEnv())
		if err != nil {
			log.Fatalf("failed creating socks5 client: %v", err)
		}
		log.Printf("created socks5 client")
	}

	// The default location for the kubeconfig file is in the user's home directory.
	var kubeconfig string
	if home := os.Getenv("HOME"); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	k8sClient, err := k8s.New(envClient, socks5Client, kubeconfig)
	if err != nil {
		log.Printf("failed creating k8s client: %v", err)
	} else {
		log.Printf("created k8s client")
		count, err := k8sClient.PodCountInDefaultNamespace()
		if err != nil {
			log.Printf("failed determining pod count: %v", err)
		} else {
			log.Printf("there are %d pods in the default namespace", count)
		}
	}

	mockClient, err := mock.New()
	if err != nil {
		log.Printf("failed creating mock client: %v", err)
	} else {
		log.Printf("created mock client")
	}

	handlerClient := handlers.NewHandlers(fiber.New(), envClient, k8sClient, mockClient)

	// Create channel for idle connections
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals
		<-sigint

		log.Print("received an interrupt signal, shutting down...")

		if err := handlerClient.App.Shutdown(); err != nil {
			log.Printf("server is not shutting down due to error from closing listeners, or context timeout! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	go func() {
		log.Fatal(handlerClient.App.Listen(*port))
	}()

	<-idleConnsClosed
}
