package webdav

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/studio-b12/gowebdav"
	"github.com/tces1/file_watcher/pkg"
)

// Watch watches a WebDAV server for changes
func Watch(watch pkg.Watch, config pkg.Config) {
	serverURL := watch.URL
	username := watch.Username
	password := watch.Password

	webdavClient := gowebdav.NewClient(serverURL, username, password)
	webdavClient.Connect()

	for {
		pathSlice := strings.Split(watch.Src, "/")
		var paths string
		for _, path := range pathSlice {
			if path != "" {
				paths = filepath.Join(paths, path)
				webdavClient.ReadDir(paths)
			}
		}
		isRefresh := downloadDirectory(webdavClient, watch.Src, watch.Dst, false)
		time.Sleep(30 * time.Second)
		if isRefresh {
			fmt.Println("Refresh Emby library...")
			pkg.RefreshEmbyLibrary(config.Emby.URL, config.Emby.APIKey)
		}
	}
}

func downloadDirectory(webdavClient *gowebdav.Client, remoteDir, localDir string, refresh bool) bool {
	fmt.Println("Scanning directory", remoteDir)
	files, err := webdavClient.ReadDir(remoteDir)

	if err != nil {
		log.Printf("Failed to list files on WebDAV server: %v", err)
		return false
	}

	for _, file := range files {
		fmt.Println("Processing file", file.Name())
		remoteFilePath := filepath.Join(remoteDir, file.Name())
		localFilePath := filepath.Join(localDir, file.Name())

		if file.IsDir() {
			// check this file if is a empty directory
			subLocalDir := filepath.Join(localDir, file.Name())
			err := os.MkdirAll(subLocalDir, 0755)
			if err != nil {
				log.Printf("Failed to create local directory: %v", err)
				continue
			}
			files, err := webdavClient.ReadDir(remoteFilePath)
			if err != nil {
				log.Printf("Failed to list files on WebDAV server: %v", err)
				continue
			}
			if len(files) != 0 {
				refresh = downloadDirectory(webdavClient, remoteFilePath, subLocalDir, refresh)
			} else {
				fmt.Printf("Empty directory %s, remove source directory\n", remoteFilePath)
				webdavClient.Remove(remoteFilePath)
			}
		} else {
			localFileInfo, err := os.Stat(localFilePath)
			if err == nil && localFileInfo.Size() == file.Size() {
				fmt.Printf("File %s skip, remove source file\n", remoteFilePath)
				webdavClient.Remove(remoteFilePath)
			} else {
				os.RemoveAll(localFilePath)
				reader, _ := webdavClient.ReadStream(remoteFilePath)
				fileSteam, _ := os.Create(localFilePath)
				defer fileSteam.Close()
				io.Copy(fileSteam, reader)
				fmt.Printf("File %s downloaded\n", file.Name())
				refresh = true
			}
		}
	}
	return refresh
}
