package main_test

import (
	"fmt"
	"net"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTerraformProviderCSBPG(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TerraformProviderCSBPG Suite")
}

var _ = BeforeSuite(func() {
	createVolume("ssl_postgres")
})

func freePort() int {
	listener, err := net.Listen("tcp", "localhost:0")
	Expect(err).NotTo(HaveOccurred())
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func createVolume(fixtureName string) {
	fixturePath := path.Join(getPWD(), "testfixtures", fixtureName)
	mustRun("docker", "volume", "create", fixtureName)
	for _, folder := range []string{"certs", "keys", "pgconf"} {
		mustRun("docker", "run",
			"-v", fixturePath+":/fixture",
			"--mount", fmt.Sprintf("source=%s,destination=/mnt", fixtureName),
			"postgres", "rm", "-rf", "/mnt/"+folder)
		mustRun("docker", "run",
			"-v", fixturePath+":/fixture",
			"--mount", fmt.Sprintf("source=%s,destination=/mnt", fixtureName),
			"postgres", "cp", "-r", "/fixture/"+folder, "/mnt")
	}
	mustRun("docker", "run",
		"-v", fixturePath+":/fixture",
		"--mount", fmt.Sprintf("source=%s,destination=/mnt", fixtureName),
		"postgres", "chmod", "-R", "0600", "/mnt/keys/server.key")
	mustRun("docker", "run",
		"-v", fixturePath+":/fixture",
		"--mount", fmt.Sprintf("source=%s,destination=/mnt", fixtureName),
		"postgres", "chown", "-R", "postgres:postgres", "/mnt/keys/server.key")
}

func mustRun(command ...string) {
	GinkgoWriter.Printf("running: %s\n", strings.Join(command, " "))
	start, err := gexec.Start(exec.Command(
		command[0], command[1:]...,
	), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(start).WithTimeout(time.Minute).WithPolling(time.Second).Should(gexec.Exit(0))
}

func getPWD() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}
