/* Author : Shravan Aras. */
/* Module which simulates running in LTE only. */

package modules

import ("lib")

type Module_ed3g struct {
	state int;
	ltecount int;
	firstevent bool;
	bcount float64;
}

func (m *Module_ed3g) Init() bool {
	m.state = lib.C3G;
	m.ltecount = 0;
	m.firstevent = true;
	m.bcount = 0.0;

	return true;
}

func (m *Module_ed3g) GetState() (int, bool) {
	return m.state, false;
}

func (m *Module_ed3g) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {
	if islte {
		m.ltecount ++;
	}

	return true;
}

func (m *Module_ed3g) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_ed3g) ServeData(data float64) bool {
	m.bcount ++;

	if lib.GetBand(data) == lib.CLTE {
		if m.firstevent {
			m.firstevent = false;
		}
	}

	return m.firstevent;
}

func (m *Module_ed3g) GetSwitchingEnergy() float64 {
	return 0.0;
}
func (m *Module_ed3g) GetDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_ed3g) GetDelayTransition() float64 {
	return 0.0;
}

func (m *Module_ed3g) GetAvgDelayTransition() float64 {
	return 0.0;
}

func (m *Module_ed3g) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_ed3g) GetCorrect() int {
	return 0;
}

func (m *Module_ed3g) GetMissed() int {
	return 0;
}

func (m *Module_ed3g) GetUnnecesary() int {
	return 0;
}

func (m *Module_ed3g) GetServe3G() int {
	return 0;
}

func (m *Module_ed3g) GetTotal() int {
	return 0;
}

func (m *Module_ed3g) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_ed3g) Reset() bool {

	return true;
}

func (m *Module_ed3g) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}
