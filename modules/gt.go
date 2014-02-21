/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib")

type EventSwitch struct {
    data_consumed float64;
    arr_data_consumed []float64
    count int;
    jump bool;
    spike_time float64;
    arr_spike_time []float64;
}

var eventtable map[string]EventSwitch;

type Module_GT struct {
	state int;
	delay_transition float64;
	delay_transmission float64;
	switching_energy float64;
	switch_flag bool;
	ltecount int;
}

func (m *Module_GT) Init() bool {
	m.state = lib.C3G;
	m.delay_transition = 0.0;
	m.delay_transmission = 0.0;
	m.switch_flag = false;
	m.ltecount = 0;
	eventtable = make(map[string]EventSwitch);

	return true;
}

func (m *Module_GT) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

/* Looking at the data and spike time we decide whether
 * we need to switch or not. */
func (m *Module_GT) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {

	element, found := eventtable[event.Name];

	if islte {
		m.ltecount ++;
	}

	if found {
		if element.jump {
			m.jumpAssist(m.state, lib.CLTE, spiketime, true);
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

func (m *Module_GT) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_GT) ServeData(data float64) bool {
	if m.state == lib.C3G {
		if lib.GetBand(data) == lib.CLTE {
			m.jumpAssist(m.state, lib.CLTE, 0.0, false);
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
			return true;
		}
	}

	return true;
}

/* 
 * Helper method which fills in the transition values, energy and 
 * keeps track of misc. flags. 
 */
func (m *Module_GT) jumpAssist(prev_state int, current_state int,
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

func (m *Module_GT) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_GT) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_GT) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_GT) GetAvgDelayTransition() float64 {
	return m.GetDelayTransition() / float64(m.ltecount);
}

func (m *Module_GT) GetAvgDelayTransmission() float64 {
	return 0.0;
}
