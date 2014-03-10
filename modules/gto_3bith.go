/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib"
		"fmt"
		"os")

type Module_GTO3H struct {
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
}

func (m *Module_GTO3H) Init() bool {
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
	eventtable2h = make(map[string]EventSwitch2H);

	return true;
}

func (m *Module_GTO3H) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

/* Looking at the data and spike time we decide whether
 * we need to switch or not. */
func (m *Module_GTO3H) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {
	m.count ++;

	element, found := eventtable2h[event.Name];

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
			if lib.IsWasteful(m.state) {
				m.unnecesary ++;
			}
			m.jumpAssist(m.state, lib.CLTE, spiketime, true);
		} else {
			if islte && m.state == lib.C3G {
				m.serve_in_3g = true;
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
		e := EventSwitch2H{0, make([]int, 8)};
		for i:=0; i < len(e.predict_table); i++ {
			e.predict_table[i] = 7;
		}
		eventtable2h[event.Name] = e;
		if m.state == lib.C3G {
			m.avg_count ++;
			m.delay_transition += lib.Switch_3gtolte;
			if m.isfirstlteevent {
				m.first_delay_transition += lib.Switch_3gtolte;
			}
		}
	}

	if element, present := eventtable2h[event.Name]; present {

		if islte {
			element.history = lib.BitAddN(element.history, lib.CLTE, 7);
			if lib.ShouldISwitch(data, spiketime) {
				element.predict_table[element.history] = lib.BitInc(
							element.predict_table[element.history]);
			} else {
				element.predict_table[element.history] = lib.BitDec(
							element.predict_table[element.history]);
			}
		} else {
			element.history = lib.BitAddN(element.history, lib.C3G, 7);
			element.predict_table[element.history] = lib.BitDec(
							element.predict_table[element.history]);
		}

		eventtable2h[event.Name] = element;
	}

	if islte && m.isfirstlteevent {
		m.isfirstlteevent = false;
		m.first_avg_count ++;
	}

	return true;
}

func (m *Module_GTO3H) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_GTO3H) ServeData(data float64) bool {
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
func (m *Module_GTO3H) jumpAssist(prev_state int, current_state int,
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

func (m *Module_GTO3H) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_GTO3H) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_GTO3H) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_GTO3H) GetAvgDelayTransition() float64 {
	return m.GetDelayTransition() / float64(m.avg_count);
}

/* We go on-demand if we miss to serve something in lte. 
	Hence we cannot have transmission delays. */
func (m *Module_GTO3H) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_GTO3H) GetCorrect() int {
	return m.ltecount - (m.GetMissed() + m.GetUnnecesary());
}

func (m *Module_GTO3H) GetMissed() int {
	return m.missed;
}

func (m *Module_GTO3H) GetUnnecesary() int {
	return m.unnecesary;
}

func (m *Module_GTO3H) GetServe3G() int {
	return m.serve3g;
}

func (m *Module_GTO3H) GetTotal() int {
	return m.count;
}

func (m *Module_GTO3H) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_GTO3H) Reset() bool {
	m.switch_flag = true;
	m.isfirstlteevent = true;

	return true;
}

func (m *Module_GTO3H) GetFirstAvgDelayTransition() float64 {
	return m.first_delay_transition / float64(m.first_avg_count);
}

func (m *Module_GTO3H) GetHiddenLatency(spiketime float64) float64 {
	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}

	return temp;
}
