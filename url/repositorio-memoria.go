package url

type repositorioMemoria struct {
	urls map[string]*Url
}

/*
NovoRepositorioMemoria cria uma instância de um repositório em memória
e retorna o seu ponteiro.
*/
func NovoRepositorioMemoria() *repositorioMemoria {
	return &repositorioMemoria{make(map[string]*Url)}
}

/*
IdExiste retorna true se a Url existir e false para quando não existir
a partir do id informado
*/
func (rep *repositorioMemoria) IdExiste(id string) bool {
	_, existe := rep.urls[id]

	return existe
}

/*
BurcarPorId retorna uma url a partir do id informado ou retorna nulo
caso não exista.
*/
func (rep *repositorioMemoria) BuscarPorId(id string) *Url {
	return rep.urls[id]
}

/*
BuscarPorUrl retorna uma url a partir da url informada ou retorna nulo
caso não exista.
*/
func (rep *repositorioMemoria) BuscarPorUrl(url string) *Url {
	for _, u := range rep.urls {
		if u.Destino == url {
			return u
		}
	}

	return nil
}

/*
Salvar persiste no repositório uma referência a uma url
*/
func (rep *repositorioMemoria) Salvar(url Url) error {
	rep.urls[url.Id] = &url

	return nil
}
