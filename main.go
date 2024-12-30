package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	cep          = "01153000"
	urlBrasilAPI = "https://brasilapi.com.br/api/cep/v1/"
	urlViaCEP    = "http://viacep.com.br/ws/"
	timeout      = 1 * time.Second
)

func main() {
	brasilAPIChannel := make(chan *http.Response)
	viaCEPChannel := make(chan *http.Response)

	go requisitaAPI(urlBrasilAPI+cep, brasilAPIChannel)
	go requisitaAPI(urlViaCEP+cep+"/json/", viaCEPChannel)

	select {
	case response := <-brasilAPIChannel:
		fmt.Println("Recebido do BrasilAPI:")
		pegaResposta(response)
	case response := <-viaCEPChannel:
		fmt.Println("Recebido do ViaCEP:")
		pegaResposta(response)
	case <-time.After(timeout):
		fmt.Println("Timeout ao aguardar respostas das APIs.")
	}
}

func requisitaAPI(url string, c chan<- *http.Response) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro ao fazer requisição para %s: %v\n", url, err)
		c <- nil
		return
	}
	c <- resp
}

func pegaResposta(resp *http.Response) {
	if resp == nil {
		fmt.Println("Erro: resposta vazia ou falha na requisição.")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler o corpo da resposta: %v\n", err)
		return
	}
	fmt.Println(string(body))
}
