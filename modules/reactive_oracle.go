/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib")

type Module_RO struct {
	state int;
	delay_transition float64;
	switching_energy float64;
	switch_flag bool;
	ltecount int;
}

func (m *Module_RO) Init() bool {
	m.state = lib.C3G;
	m.delay_transition = 0.0;
	m.switching_energy = 0.0;
	m.switch_flag = false;
	m.ltecount = 0;

	return true;
}

func (m *Module_RO) GetState() (int, bool) {
	tempflag := m.switch_flag;
	m.switch_flag = false;
	return m.state, tempflag;
}

func (m *Module_RO) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {

	if islte {
		m.ltecount ++;
		m.delay_transition += lib.Switch_3gtolte;
	}

	return false;
}

func (m *Module_RO) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_RO) ServeData(data float64) bool {
	if m.state == lib.C3G {
		if lib.GetBand(data) == lib.CLTE {
			m.jumpAssist(m.state, lib.CLTE);
		}

		return true;
	}


	/* Look ahead for the timeout period and drop down to 
	 * 3G if we don't see any LTE activity. */
	if m.state == lib.CLTE {
		if lib.GetBand(data) != lib.CLTE {
			_, bandwidth := lib.LookAhead(lib.Timeout);
			for a:=0; a < len(bandwidth); a++ {
				if lib.GetBand(bandwidth[a].Bandwidth) == lib.CLTE {
					m.state = lib.CLTE;
					return true;
				}
			}
			m.jumpAssist(m.state, lib.C3G);
			return true;
		}
	}

	return true;
}

/* 
 * Helper method which fills in the transition values, energy and 
 * keeps track of misc. flags. 
 */
func (m *Module_RO) jumpAssist(prev_state int, current_state int) {
	if prev_state == lib.C3G && current_state == lib.CLTE {
		m.switching_energy += lib.Energy_3gtolte;
		m.state = lib.CLTE;
		m.switch_flag = true;
	} else if prev_state == lib.CLTE && current_state == lib.C3G {
		m.switching_energy += lib.Energy_lteto3g;
		m.state = lib.C3G;
		m.switch_flag = true;
	}
}

func (m *Module_RO) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_RO) GetDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_RO) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_RO) GetAvgDelayTransition() float64 {
	return m.GetDelayTransition() / float64(m.ltecount);
}

func (m *Module_RO) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_RO) GetCorrect() int {
	return 0;
}

func (m *Module_RO) GetMissed() int {
	return 0;
}

func (m *Module_RO) GetUnnecesary() int {
	return 0;
}

func (m *Module_RO) GetServe3G() int {
	return 0;
}

func (m *Module_RO) GetTotal() int {
	return 0;
}

func (m *Module_RO) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_RO) Reset() bool {
	m.switch_flag = true;

	return true;
}

func (m *Module_RO) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}
