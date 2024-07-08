package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var tmpl *template.Template

func main() {
	tmpl = template.Must(template.ParseGlob("template/*.html"))

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", Home)
	http.HandleFunc("/ascii-art", DisplayASCIIArt)
	// http.HandleFunc("/download", DownloadASCIIArt)
	fmt.Println("Server is running on http://localhost:8081/")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusNotFound)
		return
	}
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing templates: %v\n", err)
	}
}

func DisplayASCIIArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "400", http.StatusBadRequest)
		return
	}
	text := r.FormValue("text")
	banner := r.FormValue("banner")
	asciiArt := DisplayAscii(text, banner)

	filename := fmt.Sprintf("ascii_art_%d.txt", time.Now().UnixNano())
	err := os.WriteFile(filename, []byte(asciiArt), 0o644)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error writing ASCII art to file: %v\n", err)
		return
	}
	data := struct {
		AsciiArt string
		Filename string
	}{
		AsciiArt: asciiArt,
		Filename: filename,
	}

	err = tmpl.ExecuteTemplate(w, "ascii.html", data)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v\n", err)
	}
}

func DisplayAscii(text, banner string) string {
	var filecontent []byte
	filecontent, err := os.ReadFile("./src/" + banner + ".txt")
	if err != nil {
		return "Error reading file"
	}
	filestring := string(filecontent)
	filestring = strings.ReplaceAll(filestring, "\r\n", "\n")
	lines := strings.Split(filestring, "\n")

	asciiArt := ""
	InputSplit := strings.Split(text, "\\n")
	fmt.Println(InputSplit)

	for _, words := range InputSplit {
		outputLines := make([]string, 8)
		fmt.Println(outputLines)
		for _, char := range words {
			charIndex := (int(char) - 32) * 9
			if charIndex < 0 || charIndex+8 >= len(lines) {
				for i := 0; i < 8; i++ {
					outputLines[i] += "        "
				}
				continue
			}
			for j := 0; j < 8; j++ {
				outputLines[j] += lines[charIndex+j]
			}
		}
		asciiArt += strings.Join(outputLines, "\n") + "\n"
	}
	fmt.Println(asciiArt)
	return asciiArt
}

// func DownloadASCIIArt(w http.ResponseWriter, r *http.Request) {
// 	filename := r.URL.Query().Get("file")
// 	if filename == "" {
// 		http.Error(w, "400 Bad Request", http.StatusNotFound)
// 		return
// 	}
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		http.Error(w, "404 Not Found", http.StatusNotFound)
// 		return
// 	}
// 	defer file.Close()

// 	w.Header().Set("Content-Disposition", "attachement; filename="+filename)
// 	w.Header().Set("Content-Type", "text/plain")
// 	http.ServeFile(w, r, filename)
// }
