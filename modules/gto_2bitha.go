/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib"
		"fmt"
		"os")

type EventSwitch2HA struct {
	history int;
	predict_table []int;	/* Jump prediction. */
	fall_bit int;
	reactive_fall_bit int;
}

var eventtable2ha map[string]EventSwitch2HA;

type Module_GTO2HA struct {
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
	isfirstlteevent bool;
	first_delay_transition float64;
	first_avg_count int;
	reg_timeout int;
	event_seen bool;
	reg_subtimeout int;
	lasteventfallbit int;
	lastreactivefallbit int;
	lasteventname string;
}

func (m *Module_GTO2HA) Init() bool {
	m.state = lib.C3G;
	m.delay_transition = 0.0;
	m.delay_transmission = 0.0;
	m.switch_flag = false;
	m.ltecount = 0;
	m.count3g = 0;
	m.serve_in_3g = false;
	m.count = 0;
	m.isfirstlteevent = true;
	m.first_delay_transition = 0.0;
	m.first_avg_count = 0.0;
	m.reg_timeout = 0;
	m.event_seen = false;
	m.reg_subtimeout = 0;
	m.lasteventfallbit = 1;
	eventtable2ha = make(map[string]EventSwitch2HA);

	return true;
}

func (m *Module_GTO2HA) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

/* Looking at the data and spike time we decide whether
 * we need to switch or not. */
func (m *Module_GTO2HA) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {
	m.count ++;

	element, found := eventtable2ha[event.Name];

	if islte {
		m.ltecount ++;
	}

	if found {
		if lib.BitJump(element.predict_table[element.history]) {
			if m.isfirstlteevent  {
				if m.state == lib.C3G {
					m.first_delay_transition += m.GetHiddenLatency(spiketime);
				}
			}
			m.jumpAssist(m.state, lib.CLTE, spiketime, true);
			m.event_seen = true;
			m.lasteventname = event.Name;
			m.lasteventfallbit = element.fall_bit;
			if !lib.ShouldISwitch(data, spiketime) {
				m.unnecesary ++;
			}
		} else {
			if islte && m.state == lib.C3G {
				m.serve_in_3g = true;
				m.event_seen = true;
				m.lasteventname = event.Name;
				m.lasteventfallbit = element.fall_bit;
				m.delay_transmission += (data / lib.Bandwidth_3g);
				if lib.ShouldISwitch(data, spiketime) {
					m.missed ++;
				} else {
				fmt.Fprintln(os.Stderr, "+", event.Timestamp, event.Name, data);
					m.serve3g ++;
				}
				m.count3g ++;
			}
		}
	} else if islte {
		eventtable2ha[event.Name] = EventSwitch2HA{0, []int{7, 7, 7, 7},
													1, 1};
		m.lasteventfallbit = 1;
		if m.state == lib.C3G {
			m.avg_count ++;
			m.delay_transition += lib.Switch_3gtolte;
			if m.isfirstlteevent {
				m.first_delay_transition += lib.Switch_3gtolte;
			}
		}
	}

	if element, present := eventtable2ha[event.Name]; present {

		if islte {
			element.history = lib.BitAdd(element.history, lib.CLTE);
			if lib.ShouldISwitch(data, spiketime) {
				element.predict_table[element.history] = lib.BitInc(
							element.predict_table[element.history]);
			} else {
				element.predict_table[element.history] = lib.BitDec(
							element.predict_table[element.history]);
			}
		} else {
			element.history = lib.BitAdd(element.history, lib.C3G);
			element.predict_table[element.history] = lib.BitDec(
							element.predict_table[element.history]);
		}

		/* Decide if I should drop down or not. */
		if lib.EnergyDelayFall(islte, data) {
			element.fall_bit = lib.BitDec(element.fall_bit);
		} else {
			element.fall_bit = lib.BitInc(element.fall_bit);
		}

		eventtable2ha[event.Name] = element;
	}

	if islte && m.isfirstlteevent {
		m.isfirstlteevent = false;
		m.first_avg_count ++;
	}

	return true;
}

func (m *Module_GTO2HA) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_GTO2HA) ServeData(data float64) bool {
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
				m.reg_subtimeout ++;
				if m.reg_subtimeout >= lib.SpikeWait {
					m.event_seen = false;
					if lib.BitFall(m.lasteventfallbit) {
						m.reg_timeout = lib.Timeout + 1;
					}
				}
			} else {
				/* Came here as a result of on-demand. */
				e, ok := eventtable2ha[m.lasteventname];
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
					eventtable2ha[m.lasteventname] = e;
				}
			}
			if m.reg_timeout >= lib.Timeout {
				m.jumpAssist(m.state, lib.C3G, 0.0, false);
				m.serve_in_3g = false;
				m.reg_timeout = 0;
			}
		} else {
			m.reg_timeout = 0;
			m.reg_subtimeout = 0;
		}
	}

	return true;
}

/* 
 * Helper method which fills in the transition values, energy and 
 * keeps track of misc. flags. 
 */
func (m *Module_GTO2HA) jumpAssist(prev_state int, current_state int,
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

func (m *Module_GTO2HA) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_GTO2HA) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_GTO2HA) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_GTO2HA) GetAvgDelayTransition() float64 {
	return m.GetDelayTransition() / float64(m.avg_count);
}

/* We go on-demand if we miss to serve something in lte. 
	Hence we cannot have transmission delays. */
func (m *Module_GTO2HA) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_GTO2HA) GetCorrect() int {
	return m.ltecount - (m.GetMissed() + m.GetUnnecesary());
}

func (m *Module_GTO2HA) GetMissed() int {
	return m.missed;
}

func (m *Module_GTO2HA) GetUnnecesary() int {
	return m.unnecesary;
}

func (m *Module_GTO2HA) GetServe3G() int {
	return m.serve3g;
}

func (m *Module_GTO2HA) GetTotal() int {
	return m.count;
}

func (m *Module_GTO2HA) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_GTO2HA) Reset() bool {
	m.switch_flag = true;
	m.isfirstlteevent = true;

	return true;
}

func (m *Module_GTO2HA) GetFirstAvgDelayTransition() float64 {
	return m.first_delay_transition / float64(m.first_avg_count);
}

func (m *Module_GTO2HA) GetHiddenLatency(spiketime float64) float64 {
	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}

	return temp;
}
