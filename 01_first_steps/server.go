package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println(r.Form)
	fmt.Println("path: ", r.URL.Path)
	fmt.Println("sheme:", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key: ", k)
		fmt.Println("val: ", v)
	}

	fmt.Fprintf(w, "Hello philip!")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		// t, _ := template.ParseFiles("login.gtpl")
		// t.Execute(w, nil)

		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseForm()

		checkFor("fruits", []string{"apple", "pear", "banana"}, r)
		checkFor("gender", []string{"1", "2"}, r)
		slice := []string{"football", "basketball", "tennis"}

		a := sliceDiff(r.Form["interest"], slice)

		if a == nil {
			fmt.Println("interest's not exist's")
		} else {
			fmt.Println("interest's: ", a)
		}
		t, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		t.ExecuteTemplate(w, "T", template.HTML("<script>alert('you have been pwned')</script>"))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // print at server side
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) // responded to clients
	}
}

func sliceDiff(a []string, b []string) []string {
	s := []string{}

	for _, av := range a {
		for _, bv := range b {
			if av == bv {
				s = append(s, av)
			}
		}
	}
	return s
}

func checkFor(target string, s []string, r *http.Request) {

	for _, v := range s {
		if v == r.Form.Get(target) {
			fmt.Println(target, ": ", v)
		} else {
			fmt.Println(target, " doesnt exist")
		}
	}

}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)

	if r.Method == "GET" {
		curTime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curTime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println("r:", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("os Err:", err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	// go func() {
	// 	time.Sleep(time.Second * 5)
	// 	postFile("246x0w.jpg", "http://localhost:9000/upload")
	// }()
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("Server Error: ", err)
	}
}

func postFile(filename string, targetURL string) error {
	fmt.Println("error writing to buffedr")
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetURL, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	return nil
}
