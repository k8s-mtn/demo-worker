package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
)

func serve(addr string) (func(context.Context) error, error) {

	http.HandleFunc("/resize", resizeHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/", notFoundHandler)

	s := http.Server{
		Addr:    addr,
		Handler: nil,
	}

	// we're ignoring the error here, which is bad
	go s.ListenAndServe()

	return s.Shutdown, nil
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Not found:", r.URL.String())

	http.NotFound(w, r)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received ping request")
	defer log.Println("Done with ping request")

	w.Write([]byte("pong"))
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Received resize request")
	defer log.Println("Done with resize request")

	r.ParseForm()

	maxX, err := strconv.Atoi(r.Form.Get("x"))
	if err != nil {
		log.Println("invalid x value")
		http.Error(w, "invalid x value", http.StatusBadRequest)
		return
	}

	maxY, err := strconv.Atoi(r.Form.Get("y"))
	if err != nil {
		log.Println("invalid y value")
		http.Error(w, "invalid y value", http.StatusBadRequest)
		return
	}

	new, err := resizeImage(r.Body, maxX, maxY)
	if err != nil {
		log.Printf("unable to send resized image to client: %s\n", err)
		http.Error(w, "unable to resize image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// if we can't copy the output, there's nothing we
	// can do to tell the client
	_, err = io.Copy(w, new)
	if err != nil {
		log.Printf("unable to send resized image to client: %s\n", err)
		return
	}
}
