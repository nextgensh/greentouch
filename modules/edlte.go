/* Author : Shravan Aras. */
/* Module which simulates running in LTE only. */

package modules

import ("lib")

type Module_edlte struct {
	state int;
	ltecount int;
	firstevent bool;
	bcount float64;
}

func (m *Module_edlte) Init() bool {
	m.state = lib.CLTE;
	m.ltecount = 0;
	m.firstevent = true;
	m.bcount = 0.0;

	return true;
}

func (m *Module_edlte) GetState() (int, bool) {
	return m.state, false;
}

func (m *Module_edlte) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {
	if islte {
		m.ltecount ++;
	}

	return true;
}

func (m *Module_edlte) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_edlte) ServeData(data float64) bool {
	m.bcount ++;

	if lib.GetBand(data) == lib.CLTE {
		if m.firstevent {
			m.firstevent = false;
		}
	}

	return m.firstevent;
}

func (m *Module_edlte) GetSwitchingEnergy() float64 {
	return 0.0;
}
func (m *Module_edlte) GetDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_edlte) GetDelayTransition() float64 {
	return 0.0;
}

func (m *Module_edlte) GetAvgDelayTransition() float64 {
	return 0.0;
}

func (m *Module_edlte) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_edlte) GetCorrect() int {
	return 0;
}

func (m *Module_edlte) GetMissed() int {
	return 0;
}

func (m *Module_edlte) GetUnnecesary() int {
	return 0;
}

func (m *Module_edlte) GetServe3G() int {
	return 0;
}

func (m *Module_edlte) GetTotal() int {
	return 0;
}

func (m *Module_edlte) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_edlte) Reset() bool {

	return true;
}

func (m *Module_edlte) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}
