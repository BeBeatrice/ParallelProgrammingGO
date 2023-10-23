/*
Scrivere un programma in Go che simuli la produzione di 5 torte da parte di 3
pasticceri. La produzione di ogni torta richiede 3 fasi che devono avvenire in ordine:
prima la torta viene cucinata, poi guarnita e infine decorata.
Il primo pasticcere si occupa solo di cucinare le torte e ci mette 1 secondo per ogni
torta. Questo pasticcere ha a disposizione 2 spazi per appoggiare le torte una volta
che ha finito di cucinarle. Se ci sono spazi liberi, può iniziare a cucinare la torta
successiva senza aspettare che il secondo pasticcere si liberi per guarnire quella
appena cucinata. Il secondo pasticcere si occupa solo di guarnire le torte e ci mette 4
secondi per ogni torta. Anche questo pasticcere ha a disposizione 2 spazi per
appoggiare le torte una volta che ha finito di guarnirle. Il terzo pasticcere si occupa
solo di decorare le torte e ci mette 8 secondi per ogni torta.
I tre pasticceri lavorano contemporaneamente.
*/

package main

import (
	"ProgettoGo/Semaforo"
	"fmt"
	"sync"
	"time"
)

type Torta struct {
	cucinata bool
	guarnita bool
	decorata bool
} //	struttura Torta: rappresentata da tre booleani che ne descrivono lo stato, settati per default a false

func cucina(torta Torta) {

	time.Sleep(1 * time.Second)
	torta.cucinata = true

} //	aspetta e poi cucina la torta

func guarnisci(torta Torta) {

	time.Sleep(4 * time.Second)
	torta.guarnita = true

} //	aspetta e poi guarnisce la torta

func decora(torta Torta) {

	time.Sleep(8 * time.Second)
	torta.decorata = true

} //	aspetta e poi decora

//	per ciascun elemento nel canale di input pasticcere compie la sua azione: una volta finito trasferisce l'oggetto nel
//	canale di output. Ciascuna funzione pasticcere chiude il suo canale input: il pasticcere vede quando le torte
//	finiscono di arrivare e sa quando smettere di lavorare

//	i pasticceri prendono inoltre come input 1 o 2 semafori (nel caso del pasticcere2), poiché come da consegna
//	non possono iniziare la loro azione se non hanno postazioni libere in cui posare le torte

func pasticcere1(ch_from, ch_to chan Torta, sem1 *Semaforo.Semaforo, wg *sync.WaitGroup) {

	numero_torta_cucinata := 1

	for t := range ch_from {
		sem1.P()
		cucina(t)
		ch_to <- t
		fmt.Printf("Torta numero %d cucinata\n", numero_torta_cucinata)
		numero_torta_cucinata++
	}

	close(ch_to)

	wg.Done()
}

func pasticcere2(ch_from, ch_to chan Torta, sem1, sem2 *Semaforo.Semaforo, wg *sync.WaitGroup) {

	numero_torta_guarnita := 1

	for t := range ch_from {
		sem2.P()
		sem1.V()
		guarnisci(t)
		ch_to <- t
		fmt.Printf("Torta numero %d guarnita\n", numero_torta_guarnita)
		numero_torta_guarnita++
	}

	close(ch_to)

	wg.Done()
}

func pasticcere3(ch_from chan Torta, ch_to chan Torta, sem2 *Semaforo.Semaforo, wg *sync.WaitGroup) {

	numero_torta_decorata := 1

	for t := range ch_from {
		sem2.V()
		decora(t)
		ch_to <- t
		fmt.Printf("Torta numero %d decorata\n", numero_torta_decorata)
		numero_torta_decorata++
	}

	close(ch_to)

	wg.Done()
}

func main() {

	var numero_torte int = 5
	var postazioni_disponibili int = 2 //	due posti per disporre le torte

	torte_da_cucinare := make(chan Torta, numero_torte)

	for i := 0; i < numero_torte; i++ {
		torte_da_cucinare <- Torta{cucinata: false, guarnita: false, decorata: false}
	}

	close(torte_da_cucinare)

	//	semaforo tra pasticcere 1 e pasticcere 2
	semaforo12 := Semaforo.Crea_semaforo(postazioni_disponibili)
	//	semaforo tra pasticcere 2 e pasticcere 3
	semaforo23 := Semaforo.Crea_semaforo(postazioni_disponibili)

	torte_cucinate := make(chan Torta, postazioni_disponibili)
	torte_guarnite := make(chan Torta, postazioni_disponibili)
	torte_decorate := make(chan Torta, numero_torte) //	cinque posti per disporre le torte

	var wg sync.WaitGroup

	inizio := time.Now().Local().UnixMilli()
	fmt.Println("*** Inizio della preparazione delle torte ***")

	wg.Add(3)

	go pasticcere1(torte_da_cucinare, torte_cucinate, semaforo12, &wg)
	go pasticcere2(torte_cucinate, torte_guarnite, semaforo12, semaforo23, &wg)
	go pasticcere3(torte_guarnite, torte_decorate, semaforo23, &wg)

	wg.Wait()

	fine := time.Now().UnixMilli()
	fmt.Println("*** Fine della preparazione delle torte ***")

	//	esprimo la durata in secondi, conversione millisecondi-secondi
	durata := (fine - inizio) / 1000

	fmt.Printf("Il tempo impiegato e' di %d secondi", durata)
}
