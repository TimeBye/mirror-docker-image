package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/go-cmd/cmd"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	green *color.Color
	c     Config
)

type Config struct {
	Batch  Batch
	Single []Single
}

type Single struct {
	Image          string
	TargetRegistry string
}

type Batch struct {
	Images                 []string
	TargetRegistry         string
	MaxConcurrentDownloads int
}

func printInfo(envCmd *cmd.Cmd) {
	ticker := time.NewTicker(1 * time.Second)
	count := 0
	for range ticker.C {
		status := envCmd.Status()
		n := len(status.Stdout)
		if n != count {
			fmt.Println(strings.Join(status.Stdout[count:n], "\n"))
		}
		count = n
	}
}

func checkErr(envCmd *cmd.Cmd) {
	status := envCmd.Status()
	if status.Exit != 0 {
		fmt.Println(strings.Join(status.Stderr, "\n"))
		os.Exit(1)
	}
}

func generateImageName(oldImageName, registry string) string {
	if !strings.HasSuffix(registry, "/") {
		registry = fmt.Sprintf("%s/", registry)
	}
	s := strings.Split(oldImageName, "/")
	tag := s[len(s)-1]
	newImageName := fmt.Sprintf("%s%s", registry, tag)
	return newImageName
}

func pull(image string, wg *sync.WaitGroup) {
	green.Printf("docker pull %s\n", image)
	envCmd := cmd.NewCmd("docker", "pull", image)
	go printInfo(envCmd)
	<-envCmd.Start()
	checkErr(envCmd)
	if wg != nil {
		wg.Done()
	}
}

func batchPull() {
	wg := &sync.WaitGroup{}
	countImages := len(c.Batch.Images)
	for i := 0; i < countImages; i += c.Batch.MaxConcurrentDownloads {
		for j := 0; j < c.Batch.MaxConcurrentDownloads; j += 1 {
			if index := i + j; index < countImages {
				wg.Add(1)
				go pull(c.Batch.Images[i+j], wg)
			}
		}
		wg.Wait()
	}
}

func tag(oldImageName, registry string, wg *sync.WaitGroup) {
	newImageName := generateImageName(oldImageName, registry)
	green.Printf("docker tag %s\t%s\n", oldImageName, newImageName)
	envCmd := cmd.NewCmd("docker", "tag", oldImageName, newImageName)
	go printInfo(envCmd)
	<-envCmd.Start()
	checkErr(envCmd)
	if wg != nil {
		wg.Done()
	}
}

func batchTag() []string {
	wg := &sync.WaitGroup{}
	countImages := len(c.Batch.Images)
	images := make([]string, 0)
	for i := 0; i < countImages; i += c.Batch.MaxConcurrentDownloads {
		for j := 0; j < c.Batch.MaxConcurrentDownloads; j += 1 {
			if index := i + j; index < countImages {
				wg.Add(1)
				go tag(c.Batch.Images[i+j], c.Batch.TargetRegistry, wg)
			}
		}
		wg.Wait()
	}
	return images
}

func push(image string, wg *sync.WaitGroup) {
	green.Printf("docker push %s\n", image)
	envCmd := cmd.NewCmd("docker", "push", image)
	go printInfo(envCmd)
	<-envCmd.Start()
	checkErr(envCmd)
	if wg != nil {
		wg.Done()
	}
}

func batchPush() {
	wg := &sync.WaitGroup{}
	countImages := len(c.Batch.Images)
	for i := 0; i < countImages; i += c.Batch.MaxConcurrentDownloads {
		for j := 0; j < c.Batch.MaxConcurrentDownloads; j += 1 {
			if index := i + j; index < countImages {
				wg.Add(1)
				go push(generateImageName(c.Batch.Images[i+j], c.Batch.TargetRegistry), wg)
			}
		}
		wg.Wait()
	}
}

func batchMirror() {
	batchPull()
	batchTag()
	batchPush()
}

func singleMirror() {
	for _, v := range c.Single {
		pull(v.Image, nil)
		tag(v.Image, v.TargetRegistry, nil)
		push(generateImageName(v.Image, v.TargetRegistry), nil)
	}
}

func init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(data, &c)
	green = color.New(color.FgHiGreen)
}

func main() {
	batchMirror()
	singleMirror()
}
