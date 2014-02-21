/* Author : Shravan Aras. */
/* Module which simulates running in 3G only. */

package modules 

import ("lib")

type Module_3g struct {
	state int;
	delay_transmission float64;
	avg_delay_transmission float64;
	count3g int;
	switch_flag bool;
}

func (m *Module_3g) Init() bool {
	m.state = lib.C3G;
	m.delay_transmission = 0.0;
	m.avg_delay_transmission = 0.0;
	m.count3g = 0;
	m.switch_flag = false;

	return true;
}

func (m *Module_3g) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

func (m *Module_3g) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {

	if islte {
		m.avg_delay_transmission += (data / lib.Bandwidth_3g);
		m.count3g ++;
	}

	m.state = lib.C3G;

	return false;
}

func (m *Module_3g) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_3g) ServeData(data float64) bool {
	if lib.GetBand(data) == lib.CLTE {
		m.delay_transmission += (data / lib.Bandwidth_3g);
	}

	return true;
}

func (m *Module_3g) GetSwitchingEnergy() float64 {
	return 0.0;
}

func (m *Module_3g) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_3g) GetDelayTransition() float64 {
	return 0.0;
}

func (m *Module_3g) GetAvgDelayTransition() float64 {
	return 0.0;
}

func (m *Module_3g) GetAvgDelayTransmission() float64 {
	return m.avg_delay_transmission / float64(m.count3g); 
}

func (m *Module_3g) GetCorrect() int {
	return 0;
}

func (m *Module_3g) GetMissed() int {
	return 0;
}

func (m *Module_3g) GetUnnecesary() int {
	return 0;
}

func (m *Module_3g) GetServe3G() int {
	return 0;
}

func (m *Module_3g) GetTotal() int {
	return 0;
}

func (m *Module_3g) GetTotalLTE() int {
	return 0;
}

func (m *Module_3g) Reset() bool {
	m.switch_flag = true;

	return true;
}

func (m *Module_3g) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}
