# file_watcher

这个工具是用来同步WebDAV指定目录到本地
我的应用场景是将群晖中aliyundrive webdav的文件同步到DSM本地，方便emby播放时可以快速拖拽进度不用转圈

## Getting started

Running it then should be as simple as:

```console
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./file_watcher main.go
$ ./file_watcher -f ./file_watcher.yaml
```

Docker

```
docker run -d -v /xxx:/xxx -v /yyy/file_watcher.yaml:/root/file_watcher.yaml tces1/file_wather
```
