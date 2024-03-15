package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "Hello World")
	})

	request := httptest.NewRequest("GET", FullLocalhost, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	bytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Hello World", string(bytes))
}

func TestRouterWithParams(t *testing.T) {
	router := httprouter.New()
	router.GET("/product/:id",func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		getId := p.ByName("id")
		fmt.Fprint(w, "Id produknya adalah " + getId)
	})

	request := httptest.NewRequest("GET", FullLocalhost + "product/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	bytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Id produknya adalah 1", string(bytes))
}

func TestRouterWithNamedParams(t *testing.T) {
	router := httprouter.New()
	router.GET("/product/:id/name/:name", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		getId := p.ByName("id")
		getName := p.ByName("name")
		fmt.Fprint(w, "Produk dengan id : " + getId + " dan namanya : "  + getName)
	})
	
	request := httptest.NewRequest("GET", FullLocalhost + "product/1/name/mangga", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Produk dengan id : 1 dan namanya : mangga", string(body))
}

func TestRouterCatchAllParams(t *testing.T) {
	router := httprouter.New()
	router.GET("/product/image/*images", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		getImages := p.ByName("images")
		fmt.Fprint(w, "Produk memiliki gambar di direktori : " + getImages)
	})
	
	request := httptest.NewRequest("GET", FullLocalhost + "product/image/name/mangga.jpg", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Produk memiliki gambar di direktori : /name/mangga.jpg", string(body))
}

//go:embed resources
var resources embed.FS

func TestServeFile(t *testing.T) {
	router := httprouter.New()
	directory, _ := fs.Sub(resources, "resources")
	router.ServeFiles("/files/*filepath", http.FS(directory))

	request := httptest.NewRequest("GET", FullLocalhost + "files/hello.txt", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Hai kawan. Mari ngoding bersama saya.", string(body))
}

func TestPanicHandler(t *testing.T) {
	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		fmt.Fprint(w, "Panic : ", i)
	}
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		panic("ada error")
	})

	// request := httptest.NewRequest("GET", FullLocalhost , nil)
	// recorder := httptest.NewRecorder()

	// router.ServeHTTP(recorder, request)

	// response := recorder.Result()
	// body, _ := io.ReadAll(response.Body)
	// assert.Equal(t, "Panic : ada error", string(body))

	server := http.Server{
		Handler: router,
		Addr: Localhost,
	}

	server.ListenAndServe()
}

/*
 secara default ketika hendak mengakses route yang tidak ada, maka akan diteruskan
 ke http.NotFound. namun kita bisa mengubah router.NotFound = http.Handler
*/

func TestRouterNotFound(t *testing.T) {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Halaman yang anda cari tidak ditemukan")
	})
	
	request := httptest.NewRequest("GET", FullLocalhost , nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Halaman yang anda cari tidak ditemukan", string(body))

}

func TestMethodNotAllowed(t *testing.T) {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Method salah")
	})

	router.POST("/post", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "Berhasil memposting postingan")
	})

	request := httptest.NewRequest("GET", FullLocalhost + "post", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Method salah", string(body))
}


// menggunakan middleware buatan sendiri karena di bawaan httprouter tidak ada middleware

type ReqMiddleware struct {
	Handler http.Handler
}

func (middleware *ReqMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request diterima/masuk")
	middleware.Handler.ServeHTTP(w, r)
}

func TestMiddleware(t *testing.T) {
	router := httprouter.New()
	router.GET("/post", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "Berhasil membuat postingan")
	})

	middleware := ReqMiddleware{Handler: router} // middleware

	request := httptest.NewRequest("GET", FullLocalhost + "post", nil)
	recorder := httptest.NewRecorder()

	middleware.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	fmt.Println(string(body))
	// output : request diterima/masuk
	// output : Berhasil membuat postingan
}