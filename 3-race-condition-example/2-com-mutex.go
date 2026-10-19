//go:build ignore

// A linha acima tira este arquivo do build normal do pacote. Sem ela, este
// arquivo e o 1-sem-mutex.go teriam duas funções main e duas variáveis number
// na mesma pasta, e nada compilaria. Rode este arquivo sozinho:
//   go run 2-com-mutex.go
//   go run -race 2-com-mutex.go

package main

import (
	"net/http"
	"strconv"
	"sync"
)

// Mesmo exemplo do 1-sem-mutex.go, agora COM lock (sync.Mutex).
//
// Cada requisição continua rodando numa goroutine própria. A diferença é que
// o acesso à variável number fica protegido por um mutex (mutual exclusion,
// exclusão mútua).
//
// Como o mutex funciona:
//   - mutex.Lock()   -> "pega a chave". Se outra goroutine já estiver com ela,
//                       esta goroutine fica bloqueada esperando a vez.
//   - mutex.Unlock() -> "devolve a chave", liberando a próxima goroutine que
//                       estava esperando.
//
// Assim, só UMA goroutine por vez executa o trecho entre Lock e Unlock (a
// "seção crítica"). Os 3 passos do number++ (ler, somar 1, gravar) acontecem
// sem outra goroutine interferir, e nenhum incremento é perdido.
//
// Teste:
//   hey -n 10000 -c 100 http://localhost:8080/   -> 10.000 requisições concorrentes
//   curl.exe http://localhost:8080/              -> requisição de número 10.001
//
// Resultado esperado: "Acessos: 10001" (sempre, em toda execução).
// No 1-sem-mutex.go o mesmo teste deu "Acessos: 9223" (778 perdidos).
var (
	number uint64 = 0
	mutex  sync.Mutex // protege number: só acesse number com o mutex travado
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		// Início da seção crítica: só uma goroutine por vez chega aqui.
		number++
		// A leitura também precisa estar protegida. Se lêssemos number depois
		// do Unlock, outra goroutine poderia alterá-lo ao mesmo tempo (race).
		// Por isso copiamos o valor para uma variável local ainda com o lock.
		current := number
		mutex.Unlock()
		// Fim da seção crítica.

		// Escrever a resposta fica FORA do lock: é a parte lenta (rede) e não
		// mexe em number. Segurar o lock só pelo mínimo necessário evita que as
		// outras requisições fiquem esperando à toa.
		w.Write([]byte("Acessos: " + strconv.FormatUint(current, 10)))
	})
	http.ListenAndServe(":8080", nil)
}
