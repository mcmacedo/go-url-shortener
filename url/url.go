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
	repo Repositorio
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
Repositorio define a interface para a implementação de repositórios de Urls
*/
type Repositorio interface {
	IdExiste(id string) bool
	BuscarPorId(id string) *Url
	BuscarPorUrl(url string) *Url
	Salvar(url Url) error
	RegistrarClick(id string)
	BuscarClicks(id string) int
}

type Url struct {
	Id      string
	Criacao time.Time
	Destino string
}

func generateId() string {

	criaId := func() string {
		id := ""
		for {
			if len(id) == tamanho {
				break
			}
			id += string(simbolos[rand.Intn(len(simbolos))])
		}
		return id
	}

	for {
		if id := criaId(); !repo.IdExiste(id) {
			return id
		}
	}

}

/*
ConfigurarReposotirio inicializa a variavel 'repo' do pacote com a implementação
de um repositório
*/
func ConfigurarReposotirio(repositorio Repositorio) {
	repo = repositorio
}

/*
BuscarOuCriarNovaUrl identifica se a url já existe ou cria uma url encurtada caso não
exista,devolvendo a url original e a encurtada.
*/
func BuscarOuCriarNovaUrl(destino string) (novaUrl *Url, nova bool, err error) {
	if _, err = url.ParseRequestURI(destino); err != nil {
		return nil, false, err
	}

	if novaUrl = repo.BuscarPorUrl(destino); novaUrl == nil {
		novaUrl := Url{generateId(), time.Now(), destino}
		repo.Salvar(novaUrl)
		return &novaUrl, true, err
	}

	return novaUrl, false, err
}

/*
Buscar pesquisa na base a url que possui o id igual ao informado e a retorna.
*/
func Buscar(id string) (url *Url, ok bool) {
	ok = false

	url = repo.BuscarPorId(id)
	if url != nil {
		ok = true
	}

	return url, ok
}

/*
RegistrarClick registra uma 'click' para o id informado no servidor.
*/
func RegistrarClick(id string) {
	repo.RegistrarClick(id)
}
