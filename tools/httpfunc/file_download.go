package httpfunc

import (
	"fmt"
	diylog "github.com/Ho-Go-Music/GoServer/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func FileDownload(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//filePath := strings.TrimPrefix(path, "/static")
	filePath := strings.Replace(path, "/static", "/public", 1)
	//Retrieve  information about the current stack frame by passing 0 as a parameter
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Failed to get current file path")
		return
	}
	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Println("Failed to get absolute path:", err)
		return
	}
	// third-level parent directory
	for i := 0; i < 3; i++ {
		absolutePath = filepath.Dir(absolutePath)
	}
	// splicing target file path
	filePath = absolutePath + filePath
	//diylog.Sugar.Infoln(filePath)
	//open file
	file, err := os.Open(filePath)
	if err != nil {
		diylog.Sugar.Errorln(err)
		http.Error(w, "no such file", http.StatusNotFound)
		return
	}
	defer file.Close()
	// Retrieve file information
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get file information", http.StatusInternalServerError)
		return
	}
	// The browser will trigger the download behavior
	w.Header().Set("Content-Disposition", "attachment;filename="+fileInfo.Name())
	// Unknown file type
	w.Header().Set("Content-Type", "application/octet-stream")
	// file size
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	// write
	_, err = io.Copy(w, file)
	if err != nil {
		// 处理文件写入响应错误
		http.Error(w, "Failed to write file to response", http.StatusInternalServerError)
		return
	}
}
