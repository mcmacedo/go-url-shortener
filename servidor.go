package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mcmacedo/go-url-shortener/url"
)

var (
	porta     int
	urlBase   string
	mockedUrl string
)

func init() {
	porta = 8888
	urlBase = fmt.Sprintf("http://locahost:%d", porta)
	mockedUrl = "mocked6ZYm9"
}

/*
Headers é um tipo 'map[string]string' para representar um cabeçalho http
*/
type Headers map[string]string

func main() {
	http.HandleFunc("/api/encurtar", Encurtador)
	http.HandleFunc("/r/", Redirecionador)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", porta), nil))
}

/*
Encurtador extrai uma url da requisição, realiza o encurtamento e responde a url encurtada.
*/
func Encurtador(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		responderCom(response, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})

		return
	}

	urlNova, nova, err := url.BuscarOuCriarNovaUrl(extrairUrl(request))

	if err != nil {
		responderCom(response, http.StatusBadRequest, nil)
		return
	}

	urlCurta := fmt.Sprintf("%s/r/%s", urlBase, urlNova.Id)
	var status int

	if nova {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	responderCom(response, status, Headers{"Location": urlCurta})
}

/*
Redirecionador recupera a url original a partir do hash e realiza o redirect
*/
func Redirecionador(response http.ResponseWriter, request *http.Request) {
	caminho := strings.Split(request.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url := buscar(id); url != nil {
		urlDestino := url.(string)
		http.Redirect(response, request, urlDestino, http.StatusMovedPermanently)
	} else {
		http.NotFound(response, request)
	}
}

func buscar(id string) interface{} {
	if id != mockedUrl {
		return nil
	}

	return "https://www.github.com/mcmacedo"
}

func extrairUrl(request *http.Request) string {
	url := make([]byte, request.ContentLength, request.ContentLength)
	request.Body.Read(url)

	return string(url)
}

func responderCom(
	response http.ResponseWriter,
	status int,
	headers Headers) {

	for chave, valor := range headers {
		response.Header().Set(chave, valor)
	}

	response.WriteHeader(status)
}
