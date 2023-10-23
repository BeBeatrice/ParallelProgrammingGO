/*
Scrivere un programma in Go che simuli un'attività di trading di valute in un mercato
fittizio.
Il programma deve simulare usando la concorrenza tre coppie di valute: EUR/USD,
GBP/USD e JPY/USD, e simulare le operazioni di acquisto e vendita in parallelo.
Creare una funzione "simulateMarketData" che simuli il prezzo delle coppie di valute e
invii i dati simulati su un canale. In particolare:
• Il prezzo della coppia EUR/USD varia casualmente tra 1.0 e 1.5.
• Il prezzo della coppia GBP/USD varia casualmente tra 1.0 e 1.5.
• Il prezzo della coppia JPY/USD varia casualmente tra 0.006 e 0.009.
I prezzi vengono generati e inviati sul canale corrispondente ad intervalli regolari, in
particolare ogni secondo.
Creare una funzione "selectPair" che utilizza una "select" per gestire le operazioni di
vendita e acquisto in base alle condizioni specificate. In particolare:
• Se il prezzo di EUR/USD supera 1.20, deve vendere EUR/USD. Simulare la
vendita con un tempo di 4 secondi, cioè inserire un delay di 4 secondi prima di
confermare la vendita.
• Se il prezzo di GBP/USD scende sotto 1.35, deve acquistare GBP/USD. Simulare
l'acquisto con un tempo di 3 secondi, cioè inserire un delay di 3 secondi prima
di confermare l'acquisto.
• Se il prezzo di JPY/USD scende sotto 0.0085, deve acquistare JPY/USD.
Simulare l'acquisto con un tempo di 3 secondi, cioè inserire un delay di 3
secondi prima di confermare l'acquisto.
Il programma deve eseguire il ciclo di trading per un minuto e alla fine del ciclo deve
terminare
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Coppia struct {
	valuta string
	prezzo float32
} //	rappresenta la coppia valuta-prezzo

func random_prezzo_eur(coppia *Coppia) {
	eur_prezzo := rand.Float32()*(1.5-1) + 1
	coppia.prezzo = eur_prezzo
	fmt.Printf("prezzo EUR/USD: %f \n", coppia.prezzo)
}

func random_prezzo_gdb(coppia *Coppia) {
	gdb_prezzo := rand.Float32()*(1.5-1) + 1
	coppia.prezzo = gdb_prezzo
	fmt.Printf("prezzo GDB/USD: %f \n", coppia.prezzo)
}

func random_prezzo_jpy(coppia *Coppia) {
	jpy_prezzo := rand.Float32()*(0.009-0.006) + 0.006
	coppia.prezzo = jpy_prezzo
	fmt.Printf("prezzo JPY/USD: %f \n", coppia.prezzo)
}

func simulate_market_data(ch_to1, ch_to2, ch_to3 chan *Coppia) {

	time.Sleep(1 * time.Second)

	//	di default i prezzi sono settati di default
	//	al limite inferiore dell'intervallo per semplicità
	coppia1 := Coppia{"EUR/USD", 1.0}
	coppia2 := Coppia{"GDB/USD", 1.0}
	coppia3 := Coppia{"JPY/USD", 0.006}

	random_prezzo_eur(&coppia1)
	random_prezzo_gdb(&coppia2)
	random_prezzo_jpy(&coppia3)

	ch_to1 <- &coppia1
	ch_to2 <- &coppia2
	ch_to3 <- &coppia3

}

func select_pair(ch_eur, ch_gdb, ch_jpy chan *Coppia, wg *sync.WaitGroup) {

	//	la select in base alla provenienza della coppia che arriva
	//	esegue il confronto con i prezzi stabiliti
	select {

	case coppia := <-ch_eur:
		if coppia.prezzo > 1.20 {
			// vendo
			time.Sleep(4 * time.Second)
			fmt.Printf("Vendita di EUR/USD effettuata con prezzo %f\n", coppia.prezzo)
		}
	case coppia := <-ch_gdb:
		if coppia.prezzo < 1.35 {
			//acquisto
			time.Sleep(3 * time.Second)
			fmt.Printf("Acquisto di GDB/USD effettuato con prezzo %f\n", coppia.prezzo)
		}
	case coppia := <-ch_jpy:
		if coppia.prezzo < 0.085 {
			//acquisto
			time.Sleep(3 * time.Second)
			fmt.Printf("Acquisto di JPY/USD effettuato con prezzo %f\n", coppia.prezzo)
		}
	default:
		fmt.Println("Nessun input")
	}

	wg.Done()

}

func main() {

	prezzi_eur := make(chan *Coppia, 1)
	prezzi_gdb := make(chan *Coppia, 1)
	prezzi_jpy := make(chan *Coppia, 1)

	var wg sync.WaitGroup

	start_time := time.Now()
	inizio := start_time.Unix()

	for time.Since(start_time) < 60*time.Second {
		wg.Add(3)
		simulate_market_data(prezzi_eur, prezzi_gdb, prezzi_jpy) //	simulate genera tre valute al secondo
		go select_pair(prezzi_eur, prezzi_gdb, prezzi_jpy, &wg)  //	quindi dovrò fare 3 controlli al secondo
		go select_pair(prezzi_eur, prezzi_gdb, prezzi_jpy, &wg)  //	per procedure di vendita/acquisto
		go select_pair(prezzi_eur, prezzi_gdb, prezzi_jpy, &wg)
	}

	wg.Wait()

	fine := time.Now().Unix()
	fmt.Println("Ciclo di trading finito")

	close(prezzi_eur)
	close(prezzi_gdb)
	close(prezzi_jpy)

	durata := fine - inizio
	fmt.Printf("Il tempo impiegato e' di %d secondi", durata)

}
