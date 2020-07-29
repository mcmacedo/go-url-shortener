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
	id := extrairPathID(request)

	rotaFindOrNotFound(response, request, id, func(url *url.Url) {
		http.Redirect(response, request, url.Destino, http.StatusMovedPermanently)

		red.stats <- id
	})
}

func main() {
	url.ConfigurarReposotirio(url.NovoRepositorioMemoria())

	stats := make(chan string)
	defer close(stats)
	go registrarEstatisticas(stats)

	http.HandleFunc("/api/encurtar", Encurtador)
	http.Handle("/r/", &Redirecionador{stats})
	http.HandleFunc("/api/stats/", Visualizador)

	logar("Iniciando servidor na porta %d...", porta)
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
	logar("URL %s encurtada com sucesso para %s.", urlNova.Destino, urlCurta)
}

/*
Visualizador recupera os Stats de uma url e os retorna se for encontrada.
*/
func Visualizador(response http.ResponseWriter, request *http.Request) {
	id := extrairPathID(request)

	rotaFindOrNotFound(response, request, id, func(url *url.Url) {
		json, err := json.Marshal(url.Stats())

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		responderComJSON(response, string(json))
	})
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

func rotaFindOrNotFound(
	response http.ResponseWriter,
	request *http.Request,
	id string,
	manager func(*url.Url),
) {
	if url, ok := url.Buscar(id); ok {
		manager(url)
	} else {
		http.NotFound(response, request)
	}
}

func extrairPathID(request *http.Request) string {
	caminho := strings.Split(request.URL.Path, "/")
	return caminho[len(caminho)-1]
}

func registrarEstatisticas(ids <-chan string) {
	for id := range ids {
		url.RegistrarClick(id)
		logar("Click registrado com sucesso para %s.", id)
	}
}

func logar(formato string, valores ...interface{}) {
	log.Printf(fmt.Sprintf("%s\n", formato), valores...)
}
