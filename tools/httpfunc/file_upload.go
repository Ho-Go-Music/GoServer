package httpfunc

import (
	"fmt"
	"github.com/Ho-Go-Music/GoServer/tools"
	"io"
	"net/http"
	"os"
)

func FileUpload(w http.ResponseWriter, r *http.Request) {
	// 解析 HTTP 请求中提交的多段表单数据
	// 该方法解析时可分配的最大内存大小（1 MB）
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Retrieve uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	dst, err := os.Create(tools.Conf.StaticFilePath.Path + header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write success message to response
	fmt.Fprintln(w, "upload file successfully！")
}
