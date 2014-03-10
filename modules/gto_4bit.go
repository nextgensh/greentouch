/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib")

type Module_GTO4 struct {
	state int;
	delay_transition float64;
	delay_transition_learn float64;
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
}

var event3gtable map[string]int;

func (m *Module_GTO4) Init() bool {
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
	m.delay_transition_learn = 0.0;
	eventtable2 = make(map[string]EventSwitch2);
	event3gtable = make(map[string]int);

	return true;
}

func (m *Module_GTO4) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

/* Looking at the data and spike time we decide whether
 * we need to switch or not. */
func (m *Module_GTO4) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {

	m.count ++;

	element, found := eventtable2[event.Name];

	if islte {
		m.avg_spike_time += spiketime;
		m.ltecount ++;
	} else {
		event3gtable[event.Name] = 10;
	}

	if found {
		if lib.BitJumpN(element.counter, 8) {
			if lib.IsWasteful(m.state) {
				m.unnecesary ++;
				m.delayHelpWaste(spiketime, m.state, lib.CLTE);
			}
			if islte {
				m.delayHelp(spiketime, m.state, lib.CLTE);
			}
			m.jumpAssist(m.state, lib.CLTE, spiketime, true);
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
		m.delay_transition += lib.Switch_3gtolte;
		eventtable2[event.Name] = EventSwitch2{15};
	}

	if element, present := eventtable2[event.Name]; present {

		if islte {
			if lib.ShouldISwitch(data, spiketime) {
				element.counter = lib.BitIncN(element.counter, 15);
			} else {
				element.counter = lib.BitDec(element.counter);
			}
		} else {
			element.counter = lib.BitDec(element.counter);
		}

		eventtable2[event.Name] = element;
	}

	return true;
}

func (m *Module_GTO4) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_GTO4) ServeData(data float64) bool {
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
func (m *Module_GTO4) jumpAssist(prev_state int, current_state int,
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

func (m *Module_GTO4) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_GTO4) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_GTO4) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_GTO4) GetAvgDelayTransition() float64 {
	return m.GetDelayTransition() / float64(m.ltecount);
}

func (m *Module_GTO4) GetAvgDelayLearnTransition() float64 {
	return m.delay_transition_learn / float64(m.ltecount);
}

func (m  *Module_GTO4) GetDelayWasteTransition() float64 {
	return m.delay_transition_waste;
}

func (m *Module_GTO4) GetAvgDelayWasteTransition() float64 {
	return m.delay_transition_waste / float64(m.unnecesary);
}

func (m *Module_GTO4) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_GTO4) GetCorrect() int {
	return m.ltecount - (m.GetMissed() + m.GetUnnecesary());
}

func (m *Module_GTO4) GetMissed() int {
	return m.missed;
}

func (m *Module_GTO4) GetUnnecesary() int {
	return m.unnecesary;
}

func (m *Module_GTO4) GetServe3G() int {
	return m.serve3g;
}

func (m *Module_GTO4) GetTotal() int {
	return m.count;
}

func (m *Module_GTO4) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_GTO4) GetAvgSpikeTime() float64 {
	return m.avg_spike_time / float64(m.ltecount);
}

func (m *Module_GTO4) Reset() bool {
	//panic ("unimplemented method");
	m.switch_flag = true;

	return true;
}

func (m *Module_GTO4) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}

func (m *Module_GTO4) delayHelp (spiketime float64, prev_state int,
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

func (m *Module_GTO4) delayHelpWaste (spiketime float64, prev_state int,
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

func (m *Module_GTO4) UniqueLTE() int {
	count := 0;

	for _, _ = range(eventtable2) {
		count ++;
	}

	return count;
}

func (m *Module_GTO4) Unique3G() int {
	count := 0;

	for _, _ = range(event3gtable) {
		count ++;
	}

	return count;

}
