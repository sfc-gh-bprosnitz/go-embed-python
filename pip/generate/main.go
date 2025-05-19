package main

import (
	"fmt"
	"github.com/sfc-gh-bprosnitz/go-embed-python/pip"
	"github.com/sfc-gh-bprosnitz/go-embed-python/python"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	targetDir := "./pip/internal/data"

	// ensure we have a stable extract path for the python distribution (otherwise shebangs won't be stable)
	tmpDir := filepath.Join("/tmp", fmt.Sprintf("python-pip-bootstrap"))
	ep, err := python.NewEmbeddedPythonWithTmpDir(tmpDir, false)
	if err != nil {
		panic(err)
	}
	defer ep.Cleanup()

	bootstrapPip(ep)

	err = pip.CreateEmbeddedPipPackages2(ep, "./pip/internal/requirements.txt", "", "", nil, targetDir)
	if err != nil {
		panic(err)
	}
}

func bootstrapPip(ep *python.EmbeddedPython) {
	getPip := downloadGetPip()
	defer os.Remove(getPip)

	cmd, err := ep.PythonCmd(getPip)
	if err != nil {
		panic(err)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func downloadGetPip() string {
	resp, err := http.Get("https://bootstrap.pypa.io/pip/3.8/get-pip.py")
	if err != nil {
		panic(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		panic("failed to download get-pip.py: " + resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "get-pip.py")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		panic(err)
	}

	return tmpFile.Name()
}
