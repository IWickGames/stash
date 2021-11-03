package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		w.WriteHeader(404)
		w.Write([]byte("404: Invalid endpoint"))
		return
	}

	fmt.Println("[API] Handled (from " + fmt.Sprintf(getIP(req)) + ") /")

	file, handle, err := req.FormFile("file")
	if err != nil {
		w.WriteHeader(200)

		w.Write([]byte("<!DOCTYPE html>"))

		w.Write([]byte("<form action=\"/\" method=\"POST\" enctype=\"multipart/form-data\">"))
		w.Write([]byte("<input type=\"file\" id=\"fileUpload\" name=\"file\">"))
		w.Write([]byte("<input value=\"Submit\" type=\"submit\">"))
		w.Write([]byte("</form>"))

		w.Write([]byte("</html>"))
		return
	}
	defer file.Close()

	fmt.Println("[API] Accepted file [" + handle.Filename + "]")
	fmt.Println("[API]               | " + fmt.Sprint(handle.Size) + "/bytes")

	digest := randSeq(20)
	name := strings.ReplaceAll(digest+"-"+handle.Filename, " ", "_")

	if !ifExists("uploads") {
		os.Mkdir("uploads", 0755)
	}

	fileW, err := os.Create("uploads/" + name)
	if err != nil {
		w.Write([]byte("Upload failed: " + err.Error()))
		fmt.Println("[API] Upload (from " + fmt.Sprintf(getIP(req)) + ") failed")
		fmt.Println("           | " + err.Error())
		return
	}
	defer fileW.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.Write([]byte("Upload failed: " + err.Error()))
		fmt.Println("[API] Upload (from " + fmt.Sprintf(getIP(req)) + ") failed")
		fmt.Println("           | " + err.Error())
		return
	}
	fileW.Write(bytes)

	w.WriteHeader(200)
	w.Write([]byte("<!DOCTYPE html>"))
	w.Write([]byte("Upload completed: <a href=\"/download/" + name + "\">" + name + "</a>"))
	w.Write([]byte("</html>"))
}

func Download(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[API] Handled (from " + fmt.Sprintf(getIP(req)) + ") /download")

	id := strings.TrimPrefix(req.URL.Path, "/download/")
	filePath := "uploads/" + id

	if !ifExists(filePath) {
		w.WriteHeader(404)
		w.Write([]byte("File not found by that ID"))
		return
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to read file from disk"))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bytes)
}
