/*	Scrivete un programma in GO che simuli un’agenzia di noleggi d’auto che deve gestire
	le prenotazioni di 10 clienti. Ogni cliente noleggia un veicolo tra quelli disponibili:
	Berlina, SUV o Station Wagon.
	• Creare la struttura Cliente con il campo "nome"
	• Creare la struttura Veicolo con il campo "tipo"
	• Creare la function "noleggia" che prende come input un cliente e che prenota
	uno a caso tra i veicoli. Questa function deve anche stampare che il cliente x ha
	noleggiato il veicolo y.
	• Creare una function "stampa" che, alla fine del processo, stampa il numero di
	Berline, SUV e Station Wagon noleggiati.
	• Ogni cliente può noleggiare un veicolo contemporaneamente ad altri.
	Si noti che si possono creare ulteriori funzioni per risolvere il problema, oltre alle due
	obbligatorie, descritte sopra.
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Cliente struct {
	nome string
}

type Veicolo struct {
	tipo string
}

type Prenotazione struct {
	persona  Cliente
	macchina Veicolo
} //	struttura prenotazione: rappresenta la coppia Cliente-Veicolo

func stampa_prenotazione_corrente(z Prenotazione) {
	m := z.macchina
	c := z.persona
	fmt.Println("Prenotazione confermata ( Cliente: [" + c.nome + "], Veicolo: [" + m.tipo + "] )")
}

func (a Veicolo) equals_type(b Veicolo) bool {
	return a.tipo == b.tipo
}

func noleggia(x Cliente, c chan Prenotazione, wg *sync.WaitGroup) {

	//	determino in modo randomico un indice nell'intervallo [0, 3) per scegliere il tipo
	//	di veicolo, dove: Berlina = 0, SUV = 1, Station Wagon = 2
	index := rand.Intn(3)

	//	dichiaro p che poi definisco in modi diversi a seconda del veicolo
	var p Prenotazione

	switch index {
	case 0:
		p = Prenotazione{persona: x, macchina: Veicolo{"Berlina"}}
	case 1:
		p = Prenotazione{persona: x, macchina: Veicolo{"SUV"}}
	case 2:
		p = Prenotazione{persona: x, macchina: Veicolo{"Station Wagon"}}
	default:
		fmt.Print("Attenzione: prenotazione non andata a buon fine\n")
	}

	stampa_prenotazione_corrente(p)
	c <- p
	wg.Done()

} //	crea un oggetto prenotazione con il cliente passato come argomento e una macchina di tipologia casuale

func stampa(ic chan Prenotazione) {

	//	tre variabili dichiarate per contenere il numero di macchine noleggiate
	var numero_Berline int = 0
	var numero_SUV int = 0
	var numero_ST int = 0

	for val := range ic {
		if (val.macchina).equals_type(Veicolo{"Berlina"}) {
			numero_Berline++
		} else if (val.macchina).equals_type(Veicolo{"SUV"}) {
			numero_SUV++
		} else if (val.macchina).equals_type(Veicolo{"Station Wagon"}) {
			numero_ST++
		}
	}

	fmt.Printf("Numero di Berline noleggiate: %d \n", numero_Berline)
	fmt.Printf("Numero di SUV noleggiati: %d \n", numero_SUV)
	fmt.Printf("Numero di Station Wagon noleggiate: %d \n", numero_ST)

} //	stampa il numero di macchine, per tipo, prenotate

func main() {

	sl := []Cliente{Cliente{"Alberto"}, Cliente{"Barbara"}, Cliente{"Carla"}, Cliente{"Davide"}, Cliente{"Elisa"},
		Cliente{"Federico"}, Cliente{"Giulio"}, Cliente{"Ilaria"}, Cliente{"Lorenzo"}, Cliente{"Maria"}}

	//	i clienti sono 10
	channel := make(chan Prenotazione, 10)

	var wg sync.WaitGroup

	//	inizializzo il seed prima delle chiamate a funzione, le quali avvengono in parallelo, altrimenti il seed
	//	sarebbe lo stesso e avrei per tutti la stessa macchina
	rand.Seed(time.Now().UnixNano())

	//	una go routine per cliente
	for _, c := range sl {
		wg.Add(1)
		go noleggia(c, channel, &wg)
	}

	wg.Wait()

	close(channel)

	stampa(channel)

}
