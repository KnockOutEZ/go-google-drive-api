package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env file to use the environment variable
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := mux.NewRouter()
	p := os.Getenv("PORT")

	r.HandleFunc("/download", driveDownloader).Methods("GET")
	r.HandleFunc("/", downloadSender).Methods("GET")

	c := cors.New(cors.Options{

		// AllowCredentials: true,
		AllowedMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedOrigins:     []string{"*"},
		AllowedHeaders:     []string{"Content-Type", "content-type", "Origin", "Accept", "Access-Control-Allow-Origin"},
		OptionsPassthrough: false,
		// Enable Debugging for testing, consider disabling in production
		// Debug: true,
	})
	handler := c.Handler(r)
	fmt.Println("Server is running on port: ", strings.Split(":"+p, ":")[1])
	log.Fatal(http.ListenAndServe(":"+p, handler))
}

func driveDownloader(w http.ResponseWriter, r *http.Request) {
	url := `http://docs.google.com/uc?export=download&id=1FxaIIBfjScRn5jn2gYDULXnWzJHRB5EE`
	fileName := "file.jpg"
	fmt.Println("Downloading file...")

	output, err := os.Create(fileName)
	defer output.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("req", req)
			fmt.Println("via", via)
			fmt.Println("headers", req.Header)
			return nil
		},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error while creating cookie jar", url, "-", err)
		return
	}
	client.Jar = jar

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("GET failed", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux i686; rv:17.0) Gecko/20100101 Firefox/17.0")
	response, err := client.Do(req)
	//response, err := client.Get(url)

	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Println(response.Status)
		return
	}

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded")
	fmt.Println(jar)

}

func downloadSender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicaiton/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=file.jpg")

	http.ServeFile(w, r, "./file.jpg")
}
