package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Home(w http.ResponseWriter, req *http.Request) {
	// Validate that the path is a valid endpoint
	if req.URL.Path != "/" {
		w.WriteHeader(404)
		w.Write([]byte("404: Invalid endpoint"))
		return
	}

	fmt.Println("[API] Handled (from " + fmt.Sprintf(getIP(req)) + ") /")

	// Parce out if there is a file stored in the request
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
	defer file.Close() // Make sure to close file when function completes

	fmt.Println("[API] Accepted file [" + handle.Filename + "]")
	fmt.Println("[API]               | " + fmt.Sprint(handle.Size) + "/bytes")

	digest := randSeq(20)                                            // Generate 20 random bytes
	name := strings.ReplaceAll(digest+"-"+handle.Filename, " ", "_") // Generate name

	// Create uploads folder
	if !ifExists("uploads") {
		os.Mkdir("uploads", 0755)
	}

	// Create file on disk
	fileW, err := os.Create("uploads/" + name)
	if err != nil {
		w.Write([]byte("Upload failed: " + err.Error()))
		fmt.Println("[API] Upload (from " + fmt.Sprintf(getIP(req)) + ") failed")
		fmt.Println("           | " + err.Error())
		return
	}
	defer fileW.Close() // Ensure the file is closed properly when finished

	bytes, err := ioutil.ReadAll(file) // Read all data from request
	if err != nil {
		w.Write([]byte("Upload failed: " + err.Error()))
		fmt.Println("[API] Upload (from " + fmt.Sprintf(getIP(req)) + ") failed")
		fmt.Println("           | " + err.Error())
		return
	}
	fileW.Write(bytes) // Write bytes onto disk

	w.WriteHeader(200)
	w.Write([]byte("<!DOCTYPE html>"))
	w.Write([]byte("Upload completed: <a href=\"/download/" + name + "\">" + name + "</a>"))
	w.Write([]byte("</html>"))
}

func Download(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[API] Handled (from " + fmt.Sprintf(getIP(req)) + ") /download")

	// Parse out the file ID
	id := strings.TrimPrefix(req.URL.Path, "/download/")
	filePath := "uploads/" + id

	if !ifExists(filePath) { // Check if the file exists on disk
		w.WriteHeader(404)
		w.Write([]byte("File not found by that ID"))
		return
	}

	fi, err := os.Open(filePath) // Read file
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to read file from disk"))
		return
	}
	defer fi.Close()
	stat, err := fi.Stat()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to read file from disk"))
		return
	}

	fmt.Println("[API] Served file successfully (" + fmt.Sprintf(getIP(req)) + ")")
	fmt.Println("    | " + id)

	/*
		This section serves the file

		Content-Disposition and http.ServeContent help format the file correctly so that it
		removes the random code that is appended to the file.
	*/
	w.Header().Set("Content-Disposition", "attachment; filename=\""+strings.SplitN(id, "-", 2)[1]+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, req, strings.SplitN(id, "-", 2)[1], stat.ModTime(), fi)
}
