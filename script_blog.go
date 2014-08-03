package main

import (
	"appengine"
	"appengine/urlfetch"
	"fmt"
	"github.com/google/go-github/github"
	rss "github.com/maiconio/go-pkg-rss"
	"io/ioutil"
	"net/http"
	"net/url"
	sa "serviceaccount"
	"strings"
	"sync"
	"time"
)

func init() {
	http.HandleFunc("/script_blog", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	atualiza(w, r)
}

type Configuracao struct {
	UsuarioCommit string
	EmailCommit   string
	Proprietario  string
	Repositorio   string
	GithubAPIKey  string
	ListaBlogs    string
	UltimaLeitura string
}

func atualiza(w http.ResponseWriter, r *http.Request) {
	conf := lerConfiguracao(w, r)
	listaBlogs := listaBlogs(conf, w, r)

	var wg sync.WaitGroup
	wg.Add(len(listaBlogs))

	for i := 0; i < len(listaBlogs); i++ {
		wLocal := w
		rLocal := r
		go gravaUltimasEntradas(listaBlogs[i], &wg, conf, wLocal, rLocal)
	}

	wg.Wait()
	fmt.Fprint(w, "pronto")

}

func gravaUltimasEntradas(urlBlog string, wg *sync.WaitGroup, conf Configuracao, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "gravaUltimasEntradas\n")
	ultimaLeitura := conf.UltimaLeitura

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	if len(urlBlog) > 0 {
		fmt.Fprint(w, "Processando o blog: "+urlBlog+" "+ultimaLeitura+"\n")

		feed := rss.New(5, true, nil, nil)
		feed.FetchClient(urlBlog, client, nil)

		if len(feed.Channels) > 0 {
			c := feed.Channels[0]
			for i := 0; i < len(c.Items); i++ {
				dataPublicacao, _ := c.Items[i].ParsedPubDate()
				d := dataPublicacao.UTC().Format(time.RFC3339Nano)
				d = d[0:4] + d[5:7] + d[8:10] + d[11:13] + d[14:16]

				if d > ultimaLeitura {
					entrada := c.Items[i]
					ehGolang := false
					for j := 0; j < len(entrada.Categories); j++ {
						if strings.Contains(strings.ToLower(entrada.Categories[j].Text), "golang") {
							ehGolang = true
						}
					}

					if ehGolang {
						nomeArquivo := dataPublicacao.UTC().Format(time.RFC3339Nano)[0:10] + "-" + url.QueryEscape(entrada.Title) + ".md"
						conteudo := "---\n"
						conteudo = conteudo + "layout: default\n"
						conteudo = conteudo + "title: " + entrada.Title + "\n"
						conteudo = conteudo + "---\n"
						conteudo = conteudo + entrada.Content.Text

						comitaArquivo(nomeArquivo, conteudo, conf, w, r)
					}
				}
			}
		}

	}
	wg.Done()
}

func listaBlogs(conf Configuracao, w http.ResponseWriter, r *http.Request) []string {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	resp, err := client.Get(conf.ListaBlogs)
	if err == nil {
		defer resp.Body.Close()
		bytes, _ := ioutil.ReadAll(resp.Body)
		blogs := strings.Split(string(bytes), "\n")
		return blogs
	}

	return nil
}

func lerConfiguracao(w http.ResponseWriter, r *http.Request) Configuracao {
	configuracao := Configuracao{
		UsuarioCommit: "maiconio",
		EmailCommit:   "maiconscosta@gmail.com",
		Proprietario:  "maiconio",
		Repositorio:   "blog.golangbr.org",
		GithubAPIKey:  "API_KEY",
		ListaBlogs:    "https://raw.githubusercontent.com/maiconio/blog.golangbr.org/master/_BLOGS",
	}

	c := appengine.NewContext(r)
	client2, _ := sa.NewClient(c, configuracao.GithubAPIKey)
	client := github.NewClient(client2)

	opt := &github.CommitsListOptions{}
	commits, _, _ := client.Repositories.ListCommits(configuracao.Proprietario, configuracao.Repositorio, opt)

	ultimaLeitura := ""
	
	if len(commits) > 0 {
		for i := 0; i < len(commits); i++ {
			if len(*commits[i].Commit.Message) > 18 {
				if (*commits[i].Commit.Message)[0:18] == "Adicionando post: " {
					n := (commits[i].Commit.Author.Date).Format(time.RFC3339)
					ultimaLeitura = n[0:4] + n[5:7] + n[8:10] + n[11:13] + n[14:16]
					break
				}
			}
		}
	}

	configuracao.UltimaLeitura = ultimaLeitura
	fmt.Fprint(w, ultimaLeitura+"\n")

	return configuracao
}

func comitaArquivo(nomeArquivo, conteudo string, conf Configuracao, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client2, _ := sa.NewClient(c, conf.GithubAPIKey)
	client := github.NewClient(client2)

	arquivoPost, _, _, _ := client.Repositories.GetContents(conf.Proprietario, conf.Repositorio, "_posts/"+nomeArquivo, &github.RepositoryContentGetOptions{})

	if arquivoPost == nil {
		opt := &github.CommitsListOptions{}
		commits, _, err := client.Repositories.ListCommits(conf.Proprietario, conf.Repositorio, opt)

		if err == nil {
			menssagem := "Adicionando post: " + nomeArquivo
			bConteudo := []byte(conteudo)

			repositoryContentsOptions := &github.RepositoryContentFileOptions{
				Message: &menssagem,
				Content: bConteudo,
				SHA:     commits[0].SHA, //
				Committer: &github.CommitAuthor{
					Name: github.String(conf.UsuarioCommit), Email: github.String(conf.EmailCommit)},
			}

			_, _, err = client.Repositories.CreateFile(conf.Proprietario, conf.Repositorio, "_posts/"+nomeArquivo, repositoryContentsOptions)
			fmt.Fprint(w, repositoryContentsOptions)
		}
	}
}
