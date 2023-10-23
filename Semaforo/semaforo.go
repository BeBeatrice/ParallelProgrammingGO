package Semaforo

type Semaforo struct {
	ch chan struct{}
}

func Crea_semaforo(capacity int) *Semaforo {
	return &Semaforo{
		ch: make(chan struct{}, capacity),
	}
}

func (sem *Semaforo) P() {
	sem.ch <- struct{}{}
} //	acquisizione: aggiungo un'unità al buffered channel

func (sem *Semaforo) V() {
	<-sem.ch
} //	rilascio: tolgo un'unità al buffered channel
