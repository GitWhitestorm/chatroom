package server

import (
	"chatroom/logic"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterHandle() {
	inferRootDir()
	go logic.Broadcaster.Start()
	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}

var rootDir string

// 推断出项目根目录
func inferRootDir() {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		if exists(d + "/template") {
			return d
		}
		return infer(filepath.Dir(d))
	}
	rootDir = infer(cwd)
}
func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
