package local

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tces1/file_watcher/pkg"
)

// Watch 监控目录
func Watch(watch pkg.Watch, config pkg.Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer watcher.Close()

	err = filepath.Walk(watch.Src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					fmt.Println("New dir:", event.Name)
					// copyDir(event.Name, watch.Dst)
					// os.RemoveAll(event.Name)
				} else {
					fmt.Println("New file:", event.Name)
					// copyFile(event.Name, watch.Dst)
					// 可以选择删除原文件
				}
			}
		case err := <-watcher.Errors:
			fmt.Println("Error:", err)
		}
	}
}

func copyFile(src, dest string) (err error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return
}

func copyDir(src, dest string) (err error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return
	}
	if err = os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return
	}

	directory, err := os.Open(src)
	if err != nil {
		return
	}
	defer directory.Close()

	objects, err := directory.Readdir(-1)
	for _, obj := range objects {
		srcPath := filepath.Join(src, obj.Name())
		destPath := filepath.Join(dest, obj.Name())

		if obj.IsDir() {
			err = copyDir(srcPath, destPath)
			if err != nil {
				return
			}
		} else {
			err = copyFile(srcPath, destPath)
			if err != nil {
				return
			}
		}
	}
	return
}
