package modules

import ("lib")

type Module interface {
	Init() bool
	GetState() (int, bool)	/* 3G / LTE */
	HandleEvent(event lib.Event, data float64, spiketime float64,
												islte bool) bool
	SetState(state int) bool
	ServeData(data float64) bool
	GetSwitchingEnergy() float64
	GetDelayTransition() float64
	GetAvgDelayTransition() float64
	GetDelayTransmission() float64
	GetAvgDelayTransmission() float64
	GetCorrect() int;
	GetMissed() int;
	GetUnnecesary() int;
	GetServe3G() int;
	GetTotal() int;
	GetTotalLTE() int;
	Reset() bool;
	GetFirstAvgDelayTransition() float64;
}
