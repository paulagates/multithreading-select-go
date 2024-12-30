package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	cep := "01153000"
	urlBrasilAPI := "https://brasilapi.com.br/api/cep/v1/" + cep
	urlViaCEP := "http://viacep.com.br/ws/" + cep + "/json/"
	go buscaAPI(urlBrasilAPI, c1)
	go buscaAPI(urlViaCEP, c2)
	select {
	case r1 := <-c1:
		println("Recebido do brasilapi:", r1)
	case r2 := <-c2:
		println("Recebido do viacep:", r2)
	case <-time.After(1 * time.Second):
		println("timeout")
	}
}

func buscaAPI(url string, c chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		c <- fmt.Sprintf("Erro na requisição: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c <- fmt.Sprintf("Erro ao ler a resposta: %v", err)
		return
	}

	c <- string(body)
}
