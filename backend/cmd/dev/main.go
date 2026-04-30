package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	backendCmd := exec.Command("go", "run", "./backend/cmd/api")
	backendCmd.Stdout = os.Stdout
	backendCmd.Stderr = os.Stderr
	backendCmd.Stdin = os.Stdin

	frontendCmd := exec.Command("go", "run", "./backend/cmd/web")
	frontendCmd.Stdout = os.Stdout
	frontendCmd.Stderr = os.Stderr
	frontendCmd.Stdin = os.Stdin

	log.Println("Starting backend server on http://localhost:8080")
	if err := backendCmd.Start(); err != nil {
		log.Fatalf("failed to start backend: %v", err)
	}

	log.Println("Starting frontend server on http://localhost:3000")
	if err := frontendCmd.Start(); err != nil {
		_ = backendCmd.Process.Kill()
		log.Fatalf("failed to start frontend: %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	backendErr := waitInBackground(backendCmd)
	frontendErr := waitInBackground(frontendCmd)

	backendRunning := true
	frontendRunning := true

	for backendRunning || frontendRunning {
		select {
		case err := <-backendErr:
			backendRunning = false
			if err != nil {
				log.Printf("backend exited with error: %v", err)
			} else {
				log.Println("backend exited")
			}
			if frontendRunning {
				log.Println("frontend server still running on http://localhost:3000")
			}
		case err := <-frontendErr:
			frontendRunning = false
			if err != nil {
				log.Printf("frontend exited with error: %v", err)
			} else {
				log.Println("frontend exited")
			}
			if backendRunning {
				log.Println("backend server still running on http://localhost:8080")
			}
		case <-signals:
			log.Println("Stopping servers...")
			stopProcess(backendCmd)
			stopProcess(frontendCmd)
		}
	}
}

func waitInBackground(cmd *exec.Cmd) chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- cmd.Wait()
	}()
	return ch
}

func stopProcess(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		_ = cmd.Process.Kill()
	}
}
