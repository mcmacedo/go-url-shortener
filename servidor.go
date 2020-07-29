package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mcmacedo/go-url-shortener/url"
)

var (
	porta   int
	urlBase string
)

func init() {
	porta = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", porta)
}

/*
Headers é um tipo 'map[string]string' para representar um cabeçalho http
*/
type Headers map[string]string

/*
Redirecionador recupera a url original a partir do hash e realiza o redirect
*/
type Redirecionador struct {
	stats chan string
}

/*
ServeHTTP implementa o método da interface type http.Handler
*/
func (red *Redirecionador) ServeHTTP(
	response http.ResponseWriter,
	request *http.Request,
) {
	caminho := strings.Split(request.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if urlEncontrada, ok := url.Buscar(id); ok {
		http.Redirect(response, request, urlEncontrada.Destino, http.StatusMovedPermanently)

		red.stats <- id
	} else {
		http.NotFound(response, request)
	}
}

func main() {
	url.ConfigurarReposotirio(url.NovoRepositorioMemoria())

	stats := make(chan string)
	defer close(stats)
	go registrarEstatisticas(stats)

	http.HandleFunc("/api/encurtar", Encurtador)
	http.Handle("/r/", &Redirecionador{stats})
	http.HandleFunc("/api/stats/", Visualizador)

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

	urlNova, nova, err := url.BuscarOuCriarNovaUrl(extrairURL(request))

	if err != nil {
		responderCom(response, http.StatusBadRequest, nil)
		return
	}

	urlCurta := fmt.Sprintf("%s/r/%s", urlBase, urlNova.Id)
	var status int
	headers := Headers{"Location": urlCurta}

	if nova {
		status = http.StatusCreated
		headers["Link"] = fmt.Sprintf(
			"<%s/api/stats/%s>; rel=\"stats\"", urlBase, urlNova.Id)

	} else {
		status = http.StatusOK
	}

	responderCom(response, status, headers)
}

/*
Visualizador recupera os Stats de uma url e os retorna se for encontrada.
*/
func Visualizador(response http.ResponseWriter, request *http.Request) {
	caminho := strings.Split(request.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url, ok := url.Buscar(id); ok {
		json, err := json.Marshal(url.Stats())

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		responderComJSON(response, string(json))

	} else {
		http.NotFound(response, request)
	}
}

func extrairURL(request *http.Request) string {
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

func responderComJSON(
	response http.ResponseWriter,
	json string) {

	responderCom(response, http.StatusOK, Headers{
		"Content-Type": "application/json",
	})

	fmt.Fprintf(response, json)
}

func registrarEstatisticas(ids <-chan string) {
	for id := range ids {
		url.RegistrarClick(id)
		fmt.Printf("Click registrado com sucesso para %s.\n", id)
	}
}
