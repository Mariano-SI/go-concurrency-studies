package main

import (
	"net/http"
	"strconv"
)

// Exemplo de race condition (condição de corrida) SEM lock.
//
// O servidor HTTP do Go atende cada requisição em uma goroutine própria.
// Com várias requisições chegando ao mesmo tempo, várias goroutines executam
// o handler em paralelo e acessam a mesma variável global `number`.
//
// `number++` não é atômico: na prática são 3 passos (ler o valor, somar 1 e
// gravar de volta). Se duas goroutines leem o mesmo valor (ex: 500) antes de
// qualquer uma gravar, as duas gravam 501 e um incremento é perdido.
//
// Teste feito:
//   hey -n 10000 -c 100 http://localhost:8080/   -> 10.000 requisições concorrentes
//   curl.exe http://localhost:8080/              -> requisição de número 10.001
//
// Resposta esperada: "Acessos: 10001"
// Resposta obtida:   "Acessos: 9223"
// Perdidos:          10001 - 9223 = 778 incrementos (~7,8%) por acesso concorrente sem lock
//
// O valor muda a cada execução (é imprevisível). Para detectar o problema,
// rode com: go run -race .
var number uint64 = 0

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Ponto da race condition: leitura + escrita concorrente sem proteção.
		number++
		w.Write([]byte("Acessos: " + strconv.FormatUint(number, 10)))

	})
	http.ListenAndServe(":8080", nil)
}
