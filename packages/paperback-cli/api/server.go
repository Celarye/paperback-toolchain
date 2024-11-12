package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type ServerContext struct {
	Server *http.Server
	StopWG *sync.WaitGroup
}

type ServerConfiguration struct {
	Directory string
}

var mimeTypes map[string]string = map[string]string{
	".css":  "text/css",
	".eot":  "application/vnd.ms-fontobject",
	".gif":  "image/gif",
	".html": "text/html",
	".jpg":  "image/jpg",
	".js":   "text/javascript",
	".json": "application/json",
	".mp4":  "video/mp4",
	".otf":  "application/font-otf",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".ttf":  "application/font-ttf",
	".wasm": "application/wasm",
	".wav":  "audio/wav",
	".woff": "application/font-woff",
}

func RunServer(config ServerConfiguration) *ServerContext {
	if config.Directory == "" {
		log.Panicln("invalid server directory")
	}

	srv := &http.Server{Addr: ":8080"}
	wg := &sync.WaitGroup{}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("server: %v %v", req.Method, req.URL.Path)

		filePath := filepath.Join(config.Directory, filepath.Clean(req.URL.Path))
		if req.URL.Path == "/" {
			filePath = filepath.Join(config.Directory, "index.html")
		}

		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		mime := mimeTypes[filepath.Ext(filePath)]
		if len(mime) > 0 {
			w.Header().Add("content-type", mime)
		}

		w.Write(fileBytes)
	})

	wg.Add(1)
	go func() {
		defer wg.Done() // let main know we are done cleaning up

		log.Printf("server: starting HTTP server")
		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return &ServerContext{
		Server: srv,
		StopWG: wg,
	}
}

func ShutdownServer(ctx *ServerContext) {
	srv := ctx.Server
	wg := ctx.StopWG

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	wg.Wait()
	log.Printf("server: shutdown HTTP server")
}
