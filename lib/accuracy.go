package lib

type Accuracy struct {
	Missed int;
	Correct int;
	Unnecesary int;
}

func (a *Accuracy) Init() {
	a.Missed = 0;
	a.Correct = 0;
	a.Unnecesary = 0;
}

func (a *Accuracy) GetMissed() int {
	return a.Missed;
}

func (a *Accuracy) GetCorrect() int {
	return a.Correct;
}

func (a *Accuracy) GetUnnecesary() int {
	return a.Unnecesary;
}

func (a *Accuracy) IncMissed() {
	a.Missed ++;
}

func (a *Accuracy) IncCorrect() {
	a.Correct ++;
}

func (a *Accuracy) IncUnnecesary() {
	a.Unnecesary ++;
}
