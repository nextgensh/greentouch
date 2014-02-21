/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib")

type Module_O struct {
	state int;
	delay_transition float64;
	delay_transmission float64;
	switch_flag bool;
	switching_energy float64;
	ltecount int;
	serve_in_3g bool;
	count3g int;
	correct int;
	missed int;
	unnecesary int;
}

func (m *Module_O) Init() bool {
	m.state = lib.C3G;
	m.delay_transition = 0.0;
	m.delay_transmission = 0.0;
	m.switch_flag = false;
	m.switching_energy = 0.0;
	m.ltecount = 0;
	m.serve_in_3g = false;
	m.count3g = 0;

	return true;
}

func (m *Module_O) GetState() (int, bool) {
	tempflag := m.switch_flag;
	m.switch_flag = false;
	return m.state, tempflag;
}

/* Looking at the data and spike time we decide whether
 * we need to switch or not. */
func (m *Module_O) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {

	if !islte {
		return false;
	}

	m.ltecount ++;

	action := lib.ShouldISwitch(data, lib.Switch_3gtolte);

	if action {
		m.jumpAssist(m.state, lib.CLTE, lib.Switch_3gtolte, true);
		if !lib.ShouldISwitch(data, lib.Switch_3gtolte) {
			m.unnecesary ++;
		}
	} else {
		m.serve_in_3g = true;
		m.delay_transmission += (data / lib.Bandwidth_3g);
		m.count3g ++;
	}

	return false;
}

func (m *Module_O) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_O) ServeData(data float64) bool {
	if m.state == lib.C3G {
		if lib.GetBand(data) == lib.CLTE && !m.serve_in_3g {
			m.jumpAssist(m.state, lib.CLTE, 0.0, false);
		}

		if lib.GetBand(data) == lib.C3G {
			m.serve_in_3g = false;
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
			m.jumpAssist(m.state, lib.C3G, 0.0, false);
			return true;
		}
	}

	return true;
}

/* 
 * Helper method which fills in the transition values, energy and 
 * keeps track of misc. flags. 
 */
func (m *Module_O) jumpAssist(prev_state int, current_state int,
								spiketime float64, fromevent bool) {

	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}

	if prev_state == lib.C3G && current_state == lib.CLTE {
		if fromevent {
			m.delay_transition += temp;
		}
		m.switching_energy += lib.Energy_3gtolte;
		m.state = lib.CLTE;
		m.switch_flag = true;
	} else if prev_state == lib.CLTE && current_state == lib.C3G {
		m.switching_energy += lib.Energy_lteto3g;
		m.state = lib.C3G;
		m.switch_flag = true;
	}
}

func (m *Module_O) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_O) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_O) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_O) GetAvgDelayTransition() float64 {
	return 0.0;
}

func (m *Module_O) GetAvgDelayTransmission() float64 {
	return m.GetDelayTransmission() / float64(m.count3g);
}

func (m *Module_O) GetCorrect() int {
	return m.correct - (m.GetMissed() + m.GetUnnecesary());
}

func (m *Module_O) GetMissed() int {
	return m.missed;
}

func (m *Module_O) GetUnnecesary() int {
	return m.unnecesary;
}

func (m *Module_O) GetServe3G() int {
	return 0;
}

func (m *Module_O) GetTotal() int {
	return 0;
}

func (m *Module_O) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_O) Reset() bool {
	m.switch_flag = true;

	return true;
}

func (m *Module_O) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}
