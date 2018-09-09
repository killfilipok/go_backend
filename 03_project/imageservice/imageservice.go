//imageService is a packege to upload and download images
package imageservice

import (
	"context"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/killfilipok/backend_stuff/03_project/database"
)

const maxUploadSize = 2 * 1000000 // 2 mb
// UploadPath is the path to folder where images are stored
const UploadPath = "./profile_images"

// UploadFileHandler 1.check size of file 2.check type of file 3.save it to UploadPath
func UploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		fileName := r.Context().Value("user").(string)
		// parse and validate file and post parameters
		// fileType := r.PostFormValue("type")
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "imageurl", SaveImg(fileBytes, w, fileName))
		r.WithContext(ctx)
	})
}

// SaveImg 1.check type of file 2.delete all files with uid of user so the are no duplicates 3.save it to UploadPath 4.return url
func SaveImg(fileBytes []byte, w http.ResponseWriter, fileName string) string {
	// check file type, detectcontenttype only needs the first 512 bytes
	if len(fileBytes) > maxUploadSize {
		renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return ""
	}
	filetype := http.DetectContentType(fileBytes)
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/png":
		break
	default:
		renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		return ""
	}

	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		return ""
	}

	for _, ending := range [3]string{".jpeg", ".jpg", ".png"} {
		os.Remove(filepath.Join(UploadPath, fileName+ending))
	}

	newPath := filepath.Join(UploadPath, fileName+fileEndings[0])

	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return ""
	}
	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return ""
	}

	imageUrl := ("http://127.0.0.1:3000/images/" + fileName + fileEndings[0])
	_, err = database.DBCon.Exec("UPDATE users SET imageurl = $1 WHERE uid=$2",
		imageUrl, fileName)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return ""
	}
	return imageUrl
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}
