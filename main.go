package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	cep := "79052564"
	go func() {
		url := "https://cdn.apicep.com/file/apicep/"
		re := regexp.MustCompile(`(\d{5})(\d{3})`)

		// Usa a expressão regular para formatar o CEP
		formattedCEP := re.ReplaceAllString(cep, "$1-$2")
		res, err := http.Get(url + formattedCEP + ".json")
		if err != nil {
			c1 <- fmt.Sprintf("Erro ao fazer a requisição para %s: %v", url, err)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			c1 <- fmt.Sprintf("Erro ao ler a resposta de %s: %v", url, err)
			return
		}
		c1 <- fmt.Sprintf("Resposta da API: %s\n%s", url, body)

	}()
	go func() {
		url := "https://viacep.com.br/ws/"
		res, err := http.Get(url + cep + "/json")
		if err != nil {
			c2 <- fmt.Sprintf("Erro ao fazer a requisição para %s: %v", url, err)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			c2 <- fmt.Sprintf("Erro ao ler a resposta de %s: %v", url, err)
			return
		}
		c1 <- fmt.Sprintf("Resposta da API: %s\n%s", url, body)
	}()

	select {
	case msg1 := <-c1:
		println("Resposta mais rápida", msg1)
	case msg2 := <-c2:
		println("Resposta mais rápida", msg2)
	case <-time.After(time.Second * 1):
		println("Timeout")
	}

}
