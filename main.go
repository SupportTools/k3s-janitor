package main

import (
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/unix"
)

var gitCommit string
var gitBranch string

func printVersion() {
	log.Printf("Current build version: %s", gitCommit)
	log.Printf("Current build branch: %s", gitBranch)
}

type DiskStatus struct {
	All     uint64  `json:"all"`
	Used    uint64  `json:"used"`
	Free    uint64  `json:"free"`
	Avail   uint64  `json:"avail"`
	Percent float64 `json:"percent"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Avail = fs.Bavail * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	disk.Percent = float64(disk.Used) / float64(disk.All)
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func main() {
	log.Println("Starting k3s-janitor")
	printVersion()

	log.Println("Starting health check")
	r := mux.NewRouter()
	r.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(gitCommit))
		w.Write([]byte("\n"))
		w.Write([]byte(gitBranch))
		w.Write([]byte("\n"))
	})
	r.HandleFunc("/disk", func(w http.ResponseWriter, r *http.Request) {
		disk := DiskUsage("/")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatUint(disk.All, 10)))
	})
	r.HandleFunc("/disk/used", func(w http.ResponseWriter, r *http.Request) {
		disk := DiskUsage("/")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatUint(disk.Used, 10)))
	})
	r.HandleFunc("/disk/free", func(w http.ResponseWriter, r *http.Request) {
		disk := DiskUsage("/")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatUint(disk.Free, 10)))
	})
	r.HandleFunc("/disk/avail", func(w http.ResponseWriter, r *http.Request) {
		disk := DiskUsage("/")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatUint(disk.Avail, 10)))
	})
	r.HandleFunc("/disk/percent", func(w http.ResponseWriter, r *http.Request) {
		disk := DiskUsage("/")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatFloat(disk.Percent, 'f', 2, 64)))
	})
	go http.ListenAndServe(":8080", r)

	percentThesholdEnv := os.Getenv("PERCENT_THRESHOLD")
	if percentThesholdEnv == "" {
		percentThesholdEnv = "90"
	}
	percentThreshold, err := strconv.ParseFloat(percentThesholdEnv, 64)
	if err != nil {
		log.Println("Error parsing percent threshold:", err)
		os.Exit(1)
	}
	log.Println("Percent threshold:", percentThreshold)

	sleepBackgroundEnv := os.Getenv("SLEEP_BACKGROUND")
	if sleepBackgroundEnv == "" {
		sleepBackgroundEnv = "15"
	}
	sleepBackground, err := strconv.ParseInt(sleepBackgroundEnv, 10, 64)
	if err != nil {
		log.Println("Error parsing sleep:", err)
		os.Exit(1)
	}
	log.Println("Background Sleep:", sleepBackground)

	sleepForegroundEnv := os.Getenv("SLEEP_FOREGROUND")
	if sleepForegroundEnv == "" {
		sleepForegroundEnv = "5"
	}
	sleepForeground, err := strconv.ParseInt(sleepForegroundEnv, 10, 64)
	if err != nil {
		log.Println("Error parsing sleep:", err)
		os.Exit(1)
	}

	// Verify that crictl is available
	cmd := exec.Command("/bin/sh", "-c", "/var/lib/rancher/k3s/data/current/bin/crictl --version")
	stdout, err := cmd.Output()
	if err != nil {
		log.Println("crictl not available")
		log.Println("Error:", err)
		//os.Exit(1)
	}
	log.Print(string(stdout))

	// Dump crictl info
	cmd = exec.Command("/bin/sh", "-c", "/var/lib/rancher/k3s/data/current/bin/crictl info")
	stdout, err = cmd.Output()
	if err != nil {
		log.Println("crictl info not available")
		log.Println("Error:", err)
		//os.Exit(1)
	}
	log.Print(string(stdout))

	for {
		log.Println("Checking filesystem usage")
		disk := DiskUsage("/")
		log.Printf("Disk usage: %.2f%%\n", disk.Percent*100)
		log.Printf("Percent threshold: %.2f%%\n", percentThreshold)
		if disk.Percent*100 > float64(percentThreshold) {
			log.Println("Disk usage is above threshold, starting cleaning up")
			log.Println("Cleaning up unused containers")
			cmd = exec.Command("/bin/sh", "-c", "for id in `/var/lib/rancher/k3s/data/current/bin/crictl ps -a | grep -i exited | awk '{print $1}'`; do /var/lib/rancher/k3s/data/current/bin/crictl rm $id ; done")
			err = cmd.Run()
			if err != nil {
				log.Println("Error cleaning up unused containers:", err)
				os.Exit(2)
			}
			log.Println("Cleaning up unused images")
			cmd := exec.Command("/bin/sh", "-c", "/var/lib/rancher/k3s/data/current/bin/crictl", "rmi", "prune")
			err := cmd.Run()
			if err != nil {
				log.Println("Error cleaning up:", err)
				os.Exit(2)
			}
			log.Println("Successfully cleaned up, sleeping...")
			time.Sleep(time.Duration(sleepForeground) * time.Minute)
		} else {
			log.Println("Disk usage is below threshold, sleeping...")
			time.Sleep(time.Duration(sleepBackground) * time.Minute)
		}
	}
}
