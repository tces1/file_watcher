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
func Watch(watch pkg.Watch) {
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
		downloadDirectory(webdavClient, watch.Src, watch.Dst)
		time.Sleep(30 * time.Second)
	}
}

func downloadDirectory(webdavClient *gowebdav.Client, remoteDir, localDir string) {
	fmt.Println("Scanning directory", remoteDir)
	files, err := webdavClient.ReadDir(remoteDir)

	if err != nil {
		log.Printf("Failed to list files on WebDAV server: %v", err)
		return
	}

	for _, file := range files {
		fmt.Println("Processing file", file.Name())
		remoteFilePath := filepath.Join(remoteDir, file.Name())
		localFilePath := filepath.Join(localDir, file.Name())

		if file.IsDir() {
			subLocalDir := filepath.Join(localDir, file.Name())
			err := os.MkdirAll(subLocalDir, 0755)
			if err != nil {
				log.Printf("Failed to create local directory: %v", err)
				continue
			}
			downloadDirectory(webdavClient, remoteFilePath, subLocalDir)
		} else {
			localFileInfo, err := os.Stat(localFilePath)
			if err == nil && localFileInfo.Size() == file.Size() {
				fmt.Printf("File %s skip\n", remoteFilePath)
			} else {
				os.RemoveAll(localFilePath)
				reader, _ := webdavClient.ReadStream(remoteFilePath)
				fileSteam, _ := os.Create(localFilePath)
				defer fileSteam.Close()
				io.Copy(fileSteam, reader)
				fmt.Printf("File %s downloaded\n", file.Name())
			}
		}
	}
}
