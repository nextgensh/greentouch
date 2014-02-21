/* Author : Shravan Aras. */
/* Module which simulates running in LTE only. */

package modules 

import ("lib")

type Module_lte struct {
	state int;
	ltecount int;
	switch_flag bool;
}

func (m *Module_lte) Init() bool {
	m.state = lib.CLTE;
	m.ltecount = 0;
	m.switch_flag = false;

	return true;
}

func (m *Module_lte) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

func (m *Module_lte) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {
	m.state = lib.CLTE;

	if islte {
		m.ltecount ++;
	}

	return true;
}

func (m *Module_lte) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_lte) ServeData(data float64) bool {
	return true;
}

func (m *Module_lte) GetSwitchingEnergy() float64 {
	return 0.0;
}
func (m *Module_lte) GetDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_lte) GetDelayTransition() float64 {
	return 0.0;
}

func (m *Module_lte) GetAvgDelayTransition() float64 {
	return 0.0;
}

func (m *Module_lte) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_lte) GetCorrect() int {
	return 0;
}

func (m *Module_lte) GetMissed() int {
	return 0;
}

func (m *Module_lte) GetUnnecesary() int {
	return 0;
}

func (m *Module_lte) GetServe3G() int {
	return 0;
}

func (m *Module_lte) GetTotal() int {
	return 0;
}

func (m *Module_lte) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_lte) Reset() bool {
	m.switch_flag = true;

	return true;
}

func (m *Module_lte) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}
