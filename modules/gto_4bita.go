/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib")

type EventSwitch2A struct {
	counter int;
	fall_bit int;
	reactive_fall_bit int;
}

var eventtable2a map[string]EventSwitch2A;

type Module_GTO4A struct {
	state int;
	delay_transition float64;
	delay_transition_waste float64;
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
	avg_spike_time float64;
	reg_timeout int;
	lasteventfallbit int;
	event_seen bool;
	lasteventcorrect bool;
	seenlte bool;
	lasteventname string;
}

func (m *Module_GTO4A) Init() bool {
	m.state = lib.C3G;
	m.delay_transition = 0.0;
	m.delay_transmission = 0.0;
	m.switch_flag = false;
	m.ltecount = 0;
	m.count3g = 0;
	m.serve_in_3g = false;
	m.count = 0;
	m.delay_transition_waste = 0.0;
	m.avg_spike_time = 0.0;
	m.reg_timeout = 0;
	m.lasteventcorrect = false;
	m.event_seen = false;
	m.lasteventfallbit = 1;
	m.seenlte = false;
	eventtable2a = make(map[string]EventSwitch2A);

	return true;
}

func (m *Module_GTO4A) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

/* Looking at the data and spike time we decide whether
 * we need to switch or not. */
func (m *Module_GTO4A) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {

	m.count ++;

	element, found := eventtable2a[event.Name];

	if islte {
		m.avg_spike_time += spiketime;
		m.ltecount ++;
	}

	if found {
		if lib.BitJumpN(element.counter, 8) {
			if lib.IsWasteful(m.state) {
				m.unnecesary ++;
				m.delayHelpWaste(spiketime, m.state, lib.CLTE);
				m.lasteventcorrect = false;
			} else {
				m.seenlte = false;
				m.lasteventcorrect = true;
			}
			if islte {
				m.delayHelp(spiketime, m.state, lib.CLTE);
			}
			m.jumpAssist(m.state, lib.CLTE, spiketime, true);
			m.event_seen = true;
			m.lasteventname = event.Name;
			m.lasteventfallbit = element.fall_bit;
		} else {
			if islte && m.state == lib.C3G {
				m.serve_in_3g = true;
				m.delay_transmission += (data / lib.Bandwidth_3g);
				if lib.ShouldISwitch(data, spiketime) {
					m.delay_transition += lib.Switch_3gtolte;
					m.missed ++;
				} else {
					m.delay_transition += (data / lib.Bandwidth_3g);
					m.serve3g ++;
				}
				m.count3g ++;
			}
		}
	} else if islte {
		eventtable2a[event.Name] = EventSwitch2A{15, 1, 1};
	}

	if element, present := eventtable2a[event.Name]; present {

		if islte {
			if lib.ShouldISwitch(data, spiketime) {
				element.counter = lib.BitIncN(element.counter, 15);
			} else {
				element.counter = lib.BitDec(element.counter);
			}
		} else {
			element.counter = lib.BitDec(element.counter);
		}

		eventtable2a[event.Name] = element;
	}

	return true;
}

func (m *Module_GTO4A) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_GTO4A) ServeData(data float64) bool {
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
			m.reg_timeout ++;
			if m.event_seen {
				if m.lasteventcorrect &&
										m.seenlte {
					m.event_seen = false;
					m.seenlte = false;
					if lib.BitFall(m.lasteventfallbit) {
						m.reg_timeout = lib.Timeout + 1;
					}
				}
			} else {
				/* Came here as a result of on-demand. */
				e, ok := eventtable2a[m.lasteventname];
				if ok {
					/* Predict fall. */
					if lib.BitFall(e.reactive_fall_bit) {
						m.reg_timeout = lib.Timeout + 1;
					} else {
					}
					/* Learn fall. */
					if lib.EnergyDelayFall(false, 0) {
						e.reactive_fall_bit =
							lib.BitDec(e.reactive_fall_bit);
					} else {
						e.reactive_fall_bit =
							lib.BitInc1(e.reactive_fall_bit);
					}
					eventtable2a[m.lasteventname] = e;
				}
			}
			if m.reg_timeout >= lib.Timeout {
				m.jumpAssist(m.state, lib.C3G, 0.0, false);
				m.serve_in_3g = false;
				m.reg_timeout = 0;
			}
		} else {
			if m.event_seen && m.lasteventcorrect {
				m.seenlte = true;
			}
			m.reg_timeout = 0;
		}
	}

	return true;
}

/* 
 * Helper method which fills in the transition values, energy and 
 * keeps track of misc. flags. 
 */
func (m *Module_GTO4A) jumpAssist(prev_state int, current_state int,
								spiketime float64, fromevent bool) {

	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}

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

func (m *Module_GTO4A) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_GTO4A) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_GTO4A) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_GTO4A) GetAvgDelayTransition() float64 {
	return m.GetDelayTransition() / float64(m.ltecount);
}

func (m  *Module_GTO4A) GetDelayWasteTransition() float64 {
	return m.delay_transition_waste;
}

func (m *Module_GTO4A) GetAvgDelayWasteTransition() float64 {
	return m.delay_transition_waste / float64(m.unnecesary);
}

func (m *Module_GTO4A) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_GTO4A) GetCorrect() int {
	return m.ltecount - (m.GetMissed() + m.GetUnnecesary());
}

func (m *Module_GTO4A) GetMissed() int {
	return m.missed;
}

func (m *Module_GTO4A) GetUnnecesary() int {
	return m.unnecesary;
}

func (m *Module_GTO4A) GetServe3G() int {
	return m.serve3g;
}

func (m *Module_GTO4A) GetTotal() int {
	return m.count;
}

func (m *Module_GTO4A) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_GTO4A) GetAvgSpikeTime() float64 {
	return m.avg_spike_time / float64(m.ltecount);
}

func (m *Module_GTO4A) Reset() bool {
	//panic ("unimplemented method");
	m.switch_flag = true;

	return true;
}

func (m *Module_GTO4A) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}

func (m *Module_GTO4A) delayHelp (spiketime float64, prev_state int,
										current_state int) {
	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}
	if (prev_state == lib.C3G && current_state == lib.CLTE) {
		m.delay_transition += temp;
		m.avg_count ++;
	}
}

func (m *Module_GTO4A) delayHelpWaste (spiketime float64, prev_state int,
										current_state int) {
	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}
	if (prev_state == lib.C3G && current_state == lib.CLTE) {
		m.delay_transition_waste += temp;
		m.avg_count ++;
	}
}
