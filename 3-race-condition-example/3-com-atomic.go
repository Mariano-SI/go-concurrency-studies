//go:build ignore

// Assim como o 2-com-mutex.go, este arquivo fica fora do build normal do pacote
// (por causa da linha acima) e deve ser rodado sozinho:
//   go run 3-com-atomic.go
//   go run -race 3-com-atomic.go

package main

import (
	"net/http"
	"strconv"
	"sync/atomic"
)

// Mesmo exemplo do contador, agora usando operação ATÔMICA em vez de lock.
//
// Cada requisição continua rodando numa goroutine própria, mas o incremento
// é feito com atomic.AddUint64.
//
// O que é uma operação atômica:
//   - É uma operação que acontece "de uma vez só", sem poder ser interrompida
//     no meio. Não existe um momento em que outra goroutine veja o valor pela
//     metade ou leia o valor antigo antes da escrita terminar.
//   - O number++ do 1-sem-mutex.go são 3 passos separados (ler, somar 1,
//     gravar). O atomic.AddUint64 faz esses 3 passos como uma única instrução
//     do processador. Por isso nenhum incremento é perdido.
//
// Diferença para o mutex (2-com-mutex.go):
//   - O mutex bloqueia: as goroutines fazem fila e esperam a vez de entrar na
//     seção crítica. Serve para proteger qualquer trecho de código, inclusive
//     várias operações juntas.
//   - O atomic não bloqueia: o próprio processador garante a operação, sem
//     fila. É mais leve e rápido, mas só serve para operações simples em UMA
//     variável (somar, ler, gravar, trocar valor).
//   - Regra prática: contador simples -> atomic. Mais de uma variável ou
//     lógica entre leitura e escrita -> mutex.
//
// Teste:
//   hey -n 10000 -c 100 http://localhost:8080/   -> 10.000 requisições concorrentes
//   curl.exe http://localhost:8080/              -> requisição de número 10.001
//
// Resultado esperado: "Acessos: 10001", igual ao do mutex.
// No 1-sem-mutex.go o mesmo teste deu "Acessos: 9223" (778 perdidos).
var number uint64 = 0

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Recebe o ENDEREÇO da variável (&number) para alterar o valor original,
		// e devolve o novo valor já incrementado.
		//
		// Usamos esse retorno na resposta em vez de ler number de novo. Se
		// fizéssemos number++ com atomic e depois lêssemos number normalmente,
		// essa leitura comum seria uma race. Para ler, use atomic.LoadUint64.
		current := atomic.AddUint64(&number, 1)

		w.Write([]byte("Acessos: " + strconv.FormatUint(current, 10)))
	})
	http.ListenAndServe(":8080", nil)
}

// Desde o Go 1.19 também existe o tipo atomic.Uint64, que é a forma mais
// recomendada hoje porque impede acessos não atômicos por engano:
//
//   var number atomic.Uint64
//   current := number.Add(1)
