/*	Scrivete un programma in Go che conta il numero di volte in cui un determinato
	carattere "x" compare in una stringa. Il programma deve utilizzare la concorrenza,
	avviando una goroutine per ogni carattere nella stringa e verificando se il carattere
	corrisponde al carattere cercato.
	Esempio: se la stringa è "aaaaaaaaaaaaabbbbbbbbcccccddddccccccfff" e il carattere
	da cercare è 'c', il programma dovrebbe avviare una goroutine per ogni carattere della
	stringa e utilizzare un meccanismo di sincronizzazione (come un WaitGroup) e un
	channel per tenere traccia del conteggio totale dei caratteri corrispondenti.
	Inizializzare nel main una stringa di test e il carattere da cercare, e.g.:
	stringa := "aaaaaaaaaaaaabbbbbbbbcccccddddccccccfff"
	carattere := 'c'
	Alla fine del processo, il programma deve stampare il conteggio finale dei caratteri
	corrispondenti. Nel nostro esempio, il conteggio finale è 11, poiché il carattere 'c'
	compare 11 volte nella stringa.
*/

package main

import (
	"fmt"
	"sync"
)

func rune_counter(a, b rune, ic chan int, wg *sync.WaitGroup) {
	if a == b {
		temp_num := <-ic //	possibili data races, evitate dalla capienza 1 del channel
		temp_num++       //	definito nel main
		ic <- temp_num
	}

	wg.Done()
}

func main() {

	//	definisco una carattere costante che non deve poter
	//	essere modificato essendo il criterio di conteggio
	const carattere_ricercato rune = 'm'

	//	inizializzo una stringa da analizzare
	var frase_da_analizzare string = "Se, nel momento culminante, proprio quando la marionetta che rappresenta" +
		" Oreste è per vendicare la morte del padre sopra Egisto e la madre" +
		" si facesse uno strappo nel cielo di carta del teatrino, che avverrebbe?"

	fmt.Printf("Conto le ricorrenze di [%c] all'interno della frase \"%s\"\n", carattere_ricercato, frase_da_analizzare)

	//	definisco il contatore di capienza 1, tale capienza evita data races
	contatore := make(chan int, 1)
	contatore <- 0

	var wg sync.WaitGroup

	//	avvio una go routine per ciascuna rune della stringa
	for _, c := range frase_da_analizzare {
		wg.Add(1)
		go rune_counter(carattere_ricercato, c, contatore, &wg)
	}

	wg.Wait()
	close(contatore)

	ricorrenze_totali := <-contatore

	fmt.Printf("Ho trovato il carattere [%c] %d volte \n", carattere_ricercato, ricorrenze_totali)

}
