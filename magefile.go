// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Available aliases can be viewed via 'mage -l'
var (
	Aliases = map[string]interface{}{
		"bi": Build.Image,
		"bl": Build.Local,
		"bp": Build.Production,
		"cb": Clean.Binary,
	}
)

// Build specific variables
var (
	Version = struct {
		SemVer string
		Branch string
		SHA    string
	}{
		SemVer: "0.1.0",
		Branch: "",
		SHA:    "",
	}
	App = struct {
		Name      string
		Supported map[string][]string
	}{
		Name: "battleship",
		Supported: map[string][]string{
			"darwin":  []string{"amd64"},
			"linux":   []string{"amd64"},
			"windows": []string{"amd64"},
		},
	}
	Paths = struct {
		Binary    string
		Templates string
	}{
		Binary:    path.Join(".", "publish"),
		Templates: path.Join(".", "templates"),
	}
	GoEnv = []string{
		"GO111MODULE=on",
		"CGO_ENABLED=0",
	}
	LdFlags = map[string]string{}
)

type (
	Build  mg.Namespace
	Clean  mg.Namespace
	Deps   mg.Namespace
	Deploy mg.Namespace
)

func init() {
	Version.Branch, _ = sh.OutCmd("git", "rev-parse", "--abbrev-ref", "HEAD")()
	Version.SHA, _ = sh.OutCmd("git", "rev-parse", "HEAD")()

	LdFlags["SHA"] = Version.SHA
	LdFlags["SemVer"] = Version.SemVer
	LdFlags["Branch"] = Version.Branch
}

// Build the application via a local installation of GO
func (Build) Local() error {
	mg.Deps(Deps.Install)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", filepath.Join(Paths.Binary, App.Name), ".")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, GoEnv...)
	return cmd.Run()
}

// Build the binary into a docker image
func (Build) Image() error {
	fmt.Println("Building image...")
	defer elapsed()()
	output, err := sh.OutCmd("docker", "build", "-t", App.Name+":"+Version.SemVer, "--no-cache", "--rm", ".")()
	fmt.Println(output)
	return err
}

func (Build) Production() error {
	mg.Deps(Deps.Install)
	for goos, arches := range App.Supported {
		for _, arch := range arches {
			goos = strings.ToLower(goos)
			arch = strings.ToLower(arch)
			binName := []string{
				App.Name,
				goos,
				arch,
			}

			flags := []string{}
			for k, v := range LdFlags {
				flags = append(flags, fmt.Sprintf("-X \"github.com/arkticman/go-armored-warship/cmd.Version.%s=%s\"", k, v))
			}

			name := strings.Join(binName, "-")
			cmd := exec.Command("go", "build", "-o", filepath.Join(Paths.Binary, name), ".")
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", goos), fmt.Sprintf("GOARCH=%s", arch))
			cmd.Env = append(cmd.Env, GoEnv...)
			fmt.Printf("[Building] %7s:%7s ... ", goos, arch)
			start := time.Now()
			if err := cmd.Run(); err != nil {
				return err
			}
			fmt.Printf("done(%v)\n", timed(start))
		}
	}
	versionInfo()

	return nil
}

// Run the command "docker rmi $(docker images -q -f dangling=true)"
func (Clean) Images() error {
	fmt.Println("Cleaning 'untagged/dangling' (<none>) images...")
	output, err := sh.OutCmd("docker", "images", "-q", "-f", "dangling=true")()
	for _, id := range strings.Split(output, "\n") {
		_, err := sh.OutCmd("docker", "rmi", id)()
		if err != nil {
			return err
		}
	}
	fmt.Println(output)
	return err
}

// Clean up after yourself
func (Clean) Binary() {
	fmt.Println("Cleaning binaries...")
	os.RemoveAll(Paths.Binary)
}

// Manage your deps, or running package managers.
func (Deps) Install() error {
	fmt.Printf("[Installing] Dependencies ... ")
	cmd := exec.Command("go", "get")
	fmt.Println("done")
	return cmd.Run()
}

func (Deploy) Image() error {
	fmt.Printf("[Releasing] Image %s:%s ...\n", App.Name, Version.SemVer)
	cmd, err := sh.OutCmd("docker", "tag", "battleship:"+Version.SemVer, "arkticman/go-armored-warship:"+Version.SemVer)()
	if err != nil {
		return err
	}
	cmd, err := sh.OutCmd("docker", "tag", "battleship:"+Version.SemVer, "arkticman/go-armored-warship:"+Version.SHA)()
	if err != nil {
		return err
	}
	cmd, err := sh.OutCmd("docker", "tag", "battleship:"+Version.SemVer, "arkticman/go-armored-warship:"+Version.Branch)()
	if err != nil {
		return err
	}
	cmd, err = sh.OutCmd("docker", "push", "arkticman/go-armored-warship:"+Version.SemVer)()
	if err != nil {
		return err
	}
	fmt.Println(cmd)
	return nil
}

func versionInfo() {
	repeatCount := 80
	fmt.Println(strings.Repeat("=", repeatCount))
	fmt.Printf("Application: %s\nSemVer: %s\nBranch: %s\nSHA: %s\n", App.Name, Version.SemVer, Version.Branch, Version.SHA)
	fmt.Println(strings.Repeat("=", repeatCount))
}

func elapsed() func() {
	start := time.Now()
	return func() {
		fmt.Printf("Completed in %v\n", time.Since(start).Truncate(time.Second*1))
	}
}
func timed(start time.Time) time.Duration {
	return time.Since(start).Truncate(time.Millisecond * 1)
}
