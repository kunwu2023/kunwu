package main

import (
	"errors"
	"flag"
	"fmt"
	progress "github.com/jony-lee/go-progress-bar"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"kunwu_cli/lib/art"
	"kunwu_cli/lib/scan"
	"os"
	"path/filepath"
	"strings"
)

func GetAllFilePaths(dirPath string) ([]string, error) {
	var paths []string
	var mu sync.Mutex // 互斥锁保护共享资源
	var wg sync.WaitGroup

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // 跳过目录，但继续递归遍历
		}

		if info.Size() == 0 {
			fmt.Printf("Skipping empty file: %s\n", path)
			return nil
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			mu.Lock()
			paths = append(paths, path)
			mu.Unlock()
		}(path)

		return nil
	})

	if err != nil {
		return nil, err
	}

	wg.Wait()
	return paths, nil
}

func main() {
	art.PrintArt()
	var filePath string
	var cloudScan bool
	var filter bool

	flag.StringVar(&filePath, "file", "", "file path")
	flag.BoolVar(&cloudScan, "cloud", true, "cloud scan enabled")
	flag.BoolVar(&filter, "filter", true, "filter normal files")
	flag.Parse()
	// 判断是否输出帮助信息
	if flag.NFlag() == 0 || flag.Arg(0) == "-help" {
		fmt.Println("Load local engine")
	}

	if filePath == "" { // TODO 多文件检测
		fmt.Println("File path is required")
		return
	}

	if !cloudScan {
		fmt.Println("Cloud scan disabled")
	}

	fmt.Printf("File path: %s\n", filePath)
	fmt.Printf("Cloud scan: %v\n", cloudScan)
	fmt.Printf("filter normal files: %v\n", filter)
	fmt.Printf("--------------------------start----------------------------\n")
	// TODO 下面要开始检测了
	fmt.Println("Local engine scanning...")
	var cloudScanPath []string

	// 开始云端检测
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"文件路径", "检出引擎", "结果"})
	table.SetAutoWrapText(false)
	fmt.Println("cloud-based engine scan...")
	cloudScanPath, err := GetAllFilePaths(filePath)
	if err != nil {
		fmt.Println("文件夹内文件路径容获取失败", err.Error())
	}
	cloudScanPathLen := int64(len(cloudScanPath))
	if cloudScanPathLen != 0 {
		bar := progress.New(cloudScanPathLen)
		for _, path := range cloudScanPath {
			scanner := scan.NewScanner(path, cloudScan)
			lastResults, err := scanner.CloudScan()
			if err != nil {
				if strings.Contains(err.Error(), "timeout") {
					fmt.Println("Cloud connection timeout, stop cloud scanning")
					break
				} else {
					scannerL := scan.NewScanner(path, cloudScan)
					lastResultsL, _ := scannerL.CloudScan()
					table.Append([]string{path, "云端检出", lastResultsL})
				}
			} else {
				if lastResults == "" {
					lastResults = "检测超时"
				}
				table.Append([]string{path, "云端检出", lastResults})
			}
			bar.Done(1)
		}
		bar.Finish()
		table.Render()
	}

	fmt.Printf("--------------------------end----------------------------\n")
	return
}
