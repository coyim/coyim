package tor

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type tor struct {
	path string
}

var singleton tor

func init() {
	wp, _ := os.Getwd()

	singleton.path = path.Join(wp, "bin/tor")
}

func Exec() *exec.Cmd {
	cmd := exec.Command(singleton.path)
	log.Printf("[Tor manager] INFO Starting Tor daemon...")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("[Tor manager] ERROR %s", err)
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("[Tor manager] ERROR %s", err)
	}
	r := bufio.NewReader(stdout)
	for {
		line, _, _ := r.ReadLine()

		if strings.Contains(string(line), "Address already in use.") {
			log.Printf("[Tor manager] WARN %s", string(line))
			return cmd
		}
		if strings.Contains(string(line), "Bootstrapped 100%: Done") {
			log.Printf("[Tor manager] INFO %s", "Tor Bootstrap Done.")
			return cmd
		}
	}
	return cmd
}
