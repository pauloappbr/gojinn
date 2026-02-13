package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Iniciando ataque de flood de memória...")

	// Cria um bloco de texto de 1 Megabyte de tamanho
	megabyteChunk := strings.Repeat("A", 1024*1024)

	// Tenta imprimir 10 Megabytes no console (o limite do Gojinn agora é 5MB)
	for i := 1; i <= 10; i++ {
		fmt.Print(megabyteChunk)
		fmt.Printf("\n--- %d MB IMPRESSOS ---\n", i)
	}

	fmt.Println("Se você está lendo isso, o ataque funcionou e o servidor está vulnerável!")
}
