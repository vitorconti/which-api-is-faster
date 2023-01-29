package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

/*
 Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.

As duas requisições serão feitas simultaneamente para as seguintes APIs:

https://cdn.apicep.com/file/apicep/" + cep + ".json

http://viacep.com.br/ws/" + cep + "/json/

Os requisitos para este desafio são:

- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.

- O resultado da request deverá ser exibido no command line, bem como qual API a enviou.

- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.



*/

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCep struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	cep := "15501-340"
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	var parseCep1 ApiCep
	var parseCep2 ViaCep
	chan1 := make(chan bool)
	chan2 := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		resp1, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
		if err != nil {
			fmt.Println("failed to get cdn api cep")
		}
		json.NewDecoder(resp1.Body).Decode(&parseCep1)
		waitGroup.Done()

	}()
	go func() {

		resp2, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
		if err != nil {
			fmt.Println("failed to get via cep")
		}

		json.NewDecoder(resp2.Body).Decode(&parseCep2)
		waitGroup.Done()
		chan2 <- true
	}()

	waitGroup.Wait()
	select {
	case <-chan1:
		fmt.Println("first response receive from api cep", parseCep1)
	case <-chan2:
		fmt.Println("first response receive from via cep", parseCep2)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}
