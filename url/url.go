package url

import (
	"math/rand"
	"net/url"
	"time"
)

const (
	tamanho  = 5
	simbolos = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)

var (
	urls map[string]Url
)

func init() {
	urls = make(map[string]Url)
	rand.Seed(time.Now().UnixNano())
}

type Url struct {
	Id      string
	Criacao time.Time
	Destino string
}

func generateId() string {
	id := ""
	for {
		if len(id) == tamanho {
			break
		}
		id += string(simbolos[rand.Intn(len(simbolos))])
	}

	return id
}

/*
BuscarOuCriarNovaUrl identifica se a url já existe ou cria uma url encurtada caso não
exista,devolvendo a url original e a encurtada.
*/
func BuscarOuCriarNovaUrl(destino string) (novaUrl Url, nova bool, err error) {
	if _, err = url.ParseRequestURI(destino); err != nil {
		return Url{}, false, err
	}

	if value, ok := urls[destino]; ok != true {
		id := generateId()
		novaUrl = Url{id, time.Now(), destino}
		nova = true

		urls[destino] = novaUrl
	} else {
		novaUrl = value
		nova = false
	}

	return novaUrl, nova, err
}

/*
Buscar pesquisa na base a url que possui o id igual ao informado e a retorna.
*/
func Buscar(id string) (url Url, ok bool) {
	ok = false

	for _, value := range urls {
		if value.Id == id {
			url = value
			ok = true
			break
		}
	}

	return url, ok
}
