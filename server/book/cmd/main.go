package main

import (
	"log"
	"net/http"

	"github.com/basebandit/bookshop/server/book/internal/controller/book"
	metadatahttp "github.com/basebandit/bookshop/server/book/internal/gateway/metadata/http"
	ratinghttp "github.com/basebandit/bookshop/server/book/internal/gateway/rating/http"
	httphandler "github.com/basebandit/bookshop/server/book/internal/handler/http"
)

func main() {
	log.Println("starting the book service")
	metdataGateway := metadatahttp.New("localhost:8081")
	ratingGateway := ratinghttp.New("localhost:8082")
	ctrl := book.New(ratingGateway, metdataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/book", http.HandlerFunc(h.GetBook))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
