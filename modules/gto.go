/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib")

type Module_GTO struct {
	state int;
	delay_transition float64;
	delay_transmission float64;
	switching_energy float64;
	switch_flag bool;
	ltecount int;
	count3g int;
	serve_in_3g bool;
	correct int;
	missed int;
	unnecesary int;
	serve3g int;
	count int;
	avg_count int;
}

func (m *Module_GTO) Init() bool {
	m.state = lib.C3G;
	m.delay_transition = 0.0;
	m.delay_transmission = 0.0;
	m.switch_flag = false;
	m.ltecount = 0;
	m.count3g = 0;
	m.count = 0;
	m.serve_in_3g = false;
	eventtable = make(map[string]EventSwitch);

	return true;
}

func (m *Module_GTO) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

/* Looking at the data and spike time we decide whether
 * we need to switch or not. */
func (m *Module_GTO) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {

	m.count ++;

	element, found := eventtable[event.Name];

	if islte {
		m.ltecount ++;
	}

	if found {
		if element.jump {
			m.jumpAssist(m.state, lib.CLTE, spiketime, true);
			if !lib.ShouldISwitch(data, spiketime) {
				m.unnecesary ++;
			}
		} else {
			if islte && m.state == lib.C3G {
				m.serve_in_3g = true;
				m.delay_transmission += (data / lib.Bandwidth_3g);
				if lib.ShouldISwitch(data, spiketime) {
					m.missed ++;
				} else {
					m.serve3g ++;
				}
				m.count3g ++;
			}
		}
	} else if islte {
		eventtable[event.Name] = EventSwitch{0.0,
									[]float64{0.0, 0.0, 0.0, 0.0, 0.0},
									0, false,
									0.0,
									[]float64{0.0, 0.0, 0.0, 0.0, 0.0}};
	}

	if element, present := eventtable[event.Name]; present {
		element.count ++;
		element.data_consumed, element.arr_data_consumed = lib.GetAverage(
											element.arr_data_consumed,
											data,
											element.count);
		element.spike_time, element.arr_spike_time = lib.GetAverage(
											element.arr_spike_time,
											spiketime,
											element.count);

		element.jump = lib.ShouldISwitch(element.data_consumed, element.spike_time);

		eventtable[event.Name] = element;
	}

	return true;
}

func (m *Module_GTO) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_GTO) ServeData(data float64) bool {
	if m.state == lib.C3G {
		if lib.GetBand(data) == lib.CLTE {
			m.jumpAssist(m.state, lib.CLTE, 0.0, m.serve_in_3g);
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
					return true;
				}
			}
			m.jumpAssist(m.state, lib.C3G, 0.0, false);
			m.serve_in_3g = false;
			return true;
		}
	}

	return true;
}

/* 
 * Helper method which fills in the transition values, energy and 
 * keeps track of misc. flags. 
 */
func (m *Module_GTO) jumpAssist(prev_state int, current_state int,
								spiketime float64, fromevent bool) {

	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}

	if prev_state == lib.C3G && current_state == lib.CLTE {
		if fromevent {
			m.delay_transition += temp;
			m.avg_count ++;
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

func (m *Module_GTO) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_GTO) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_GTO) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_GTO) GetAvgDelayTransition() float64 {
	return m.GetDelayTransition() / float64(m.avg_count);
}

func (m *Module_GTO) GetAvgDelayTransmission() float64 {
	return m.GetDelayTransmission() / float64(m.count3g);
}

func (m *Module_GTO) GetCorrect() int {
	return m.count - (m.GetMissed() + m.GetUnnecesary());
}

func (m *Module_GTO) GetMissed() int {
	return m.missed;
}

func (m *Module_GTO) GetUnnecesary() int {
	return m.unnecesary;
}

func (m *Module_GTO) GetServe3G() int {
	return m.serve3g;
}

func (m *Module_GTO) GetTotal() int {
	return m.count;
}

func (m *Module_GTO) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_GTO) Reset() bool {

	panic("unimplemented method");

	return true;
}

func (m *Module_GTO) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}
