/* Author : Shravan Aras. */
/* Module which simulates running in reactive oracle mode. */

package modules

import ("lib"
		"fmt"
		"os")

type EventSwitch2HN struct {
	history int;
	predict_table []int;
}

var eventtable2hn map[string]EventSwitch2HN;


type Module_GTO2HN struct {
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
	avg_spike_time float64;
	history_bits int;
	predict_bits int;
}

func (m *Module_GTO2HN) Init() bool {
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
	eventtable2hn = make(map[string]EventSwitch2HN);

	return true;
}

func (m *Module_GTO2HN) GetState() (int, bool) {
	temp_flag := m.switch_flag;
	m.switch_flag = false;
	return m.state, temp_flag;
}

/* Looking at the data and spike time we decide whether
 * we need to switch or not. */
func (m *Module_GTO2HN) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {
	m.count ++;

	element, found := eventtable2hn[event.Name];

	if islte {
		m.ltecount ++;
		m.avg_spike_time += spiketime;
	}


	if found {
		if lib.BitJumpN(element.predict_table[element.history],
							int(lib.BitValue(m.predict_bits)/2)) {
			if m.isfirstlteevent  {
				if m.state == lib.C3G {
					m.first_delay_transition += m.GetHiddenLatency(spiketime);
				}
			}
			if lib.IsWasteful(m.state) {
				m.unnecesary ++;
			}
			if islte {
				m.delayHelp(spiketime, m.state, lib.CLTE);
			}
			m.jumpAssist(m.state, lib.CLTE, spiketime, true);
		} else {
			if islte && m.state == lib.C3G {
				m.serve_in_3g = true;
				m.delay_transmission += (data / lib.Bandwidth_3g);
				m.avg_count ++;
				m.delay_transition += lib.Switch_3gtolte;
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
		e := EventSwitch2HN{0, make([]int,
				lib.BitValue(m.history_bits))};
		for i:=0; i < len(e.predict_table); i++ {
			e.predict_table[i] = lib.BitValue(m.predict_bits)-1;
		}
		eventtable2hn[event.Name] = e;
		if m.state == lib.C3G {
			//m.avg_count ++;
			//m.delay_transition += lib.Switch_3gtolte;
			if m.isfirstlteevent {
				m.first_delay_transition += lib.Switch_3gtolte;
			}
		}
	}

	if element, present := eventtable2hn[event.Name]; present {

		if islte {
			element.history = lib.BitAddN(element.history, lib.CLTE,
								lib.BitValue(m.history_bits)-1);
			if lib.ShouldISwitch(data, spiketime) {
				element.predict_table[element.history] = lib.BitIncN(
							element.predict_table[element.history],
								lib.BitValue(m.predict_bits)-1);
			} else {
				element.predict_table[element.history] = lib.BitDec(
							element.predict_table[element.history]);
			}
		} else {
			element.history = lib.BitAddN(element.history, lib.C3G,
								lib.BitValue(m.history_bits)-1);
			element.predict_table[element.history] = lib.BitDec(
							element.predict_table[element.history]);
		}

		eventtable2hn[event.Name] = element;
	}

	if islte && m.isfirstlteevent {
		m.isfirstlteevent = false;
		m.first_avg_count ++;
	}

	return true;
}

func (m *Module_GTO2HN) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_GTO2HN) ServeData(data float64) bool {
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

func (m *Module_GTO2HN) delayHelp (spiketime float64, prev_state int,
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

/* 
 * Helper method which fills in the transition values, energy and 
 * keeps track of misc. flags. 
 */
func (m *Module_GTO2HN) jumpAssist(prev_state int, current_state int,
								spiketime float64, fromevent bool) {

	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}

	if prev_state == lib.C3G && current_state == lib.CLTE {
		if fromevent {
			//m.delay_transition += temp;
			//m.avg_count ++;
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

func (m *Module_GTO2HN) GetSwitchingEnergy() float64 {
	return m.switching_energy;
}

func (m *Module_GTO2HN) GetDelayTransmission() float64 {
	return m.delay_transmission;
}

func (m *Module_GTO2HN) GetDelayTransition() float64 {
	return m.delay_transition;
}

func (m *Module_GTO2HN) GetAvgDelayTransition() float64 {
	return m.GetDelayTransition() / float64(m.ltecount);
}

/* We go on-demand if we miss to serve something in lte. 
	Hence we cannot have transmission delays. */
func (m *Module_GTO2HN) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_GTO2HN) GetCorrect() int {
	return m.ltecount - (m.GetMissed() + m.GetUnnecesary());
}

func (m *Module_GTO2HN) GetMissed() int {
	return m.missed;
}

func (m *Module_GTO2HN) GetUnnecesary() int {
	return m.unnecesary;
}

func (m *Module_GTO2HN) GetServe3G() int {
	return m.serve3g;
}

func (m *Module_GTO2HN) GetTotal() int {
	return m.count;
}

func (m *Module_GTO2HN) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_GTO2HN) Reset() bool {
	m.switch_flag = true;
	m.isfirstlteevent = true;

	return true;
}

func (m *Module_GTO2HN) GetFirstAvgDelayTransition() float64 {
	return m.first_delay_transition / float64(m.first_avg_count);
}

func (m *Module_GTO2HN) GetHiddenLatency(spiketime float64) float64 {
	temp := lib.Switch_3gtolte - spiketime;
	if temp < 0.0 {
		temp = 0.0;
	}

	return temp;
}

func (m *Module_GTO2HN) GetAvgSpikeTime() float64 {
	return m.avg_spike_time / float64(m.ltecount);
}

func (m *Module_GTO2HN) SetHistoryBit(size int) {
	m.history_bits = size;
}

func (m *Module_GTO2HN) SetPredictBit(size int) {
	m.predict_bits = size;
}
