/* Author : Shravan Aras. */
/* Module which simulates running in LTE only. */

package modules

import ("lib")

type Module_stats struct {
	state int;
	ltecount int;
	datalte float64;
	data3g float64;
	time3g int;
	totaltime int;
	count int;
	idletime []int;
	idle int;
}

func (m *Module_stats) Init() bool {
	m.state = lib.CLTE;
	m.ltecount = 0;
	m.datalte = 0;
	m.data3g = 0;
	m.time3g = 0;
	m.totaltime = 0;
	m.count = 0;
	m.idle = 0;

	return true;
}

func (m *Module_stats) GetState() (int, bool) {
	return m.state, false;
}

func (m *Module_stats) HandleEvent(event lib.Event,
								data float64,
								spiketime float64,
								islte bool) bool {
	m.count ++;

	m.state = lib.CLTE;

	if islte {
		m.ltecount ++;
	}

	return true;
}

func (m *Module_stats) SetState(state int) bool {
	m.state = state;

	return true;
}

func (m *Module_stats) ServeData(data float64) bool {

	m.totaltime ++;

	if lib.GetBand(data) == lib.CLTE {
		m.datalte += data;
	} else {
		m.data3g += data;
		m.time3g ++;
	}

	if data == 0 {
		m.idle ++;
	} else {
		if m.idle > 0 {
			m.idletime = append(m.idletime, m.idle);
			m.idle = 0;
		}
	}

	return true;
}

func (m *Module_stats) GetSwitchingEnergy() float64 {
	return 0.0;
}
func (m *Module_stats) GetDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_stats) GetDelayTransition() float64 {
	return 0.0;
}

func (m *Module_stats) GetAvgDelayTransition() float64 {
	return 0.0;
}

func (m *Module_stats) GetAvgDelayTransmission() float64 {
	return 0.0;
}

func (m *Module_stats) GetCorrect() int {
	return 0;
}

func (m *Module_stats) GetMissed() int {
	return 0;
}

func (m *Module_stats) GetUnnecesary() int {
	return 0;
}

func (m *Module_stats) GetServe3G() int {
	return 0;
}

func (m *Module_stats) GetTotal() int {
	return m.count;
}

func (m *Module_stats) GetTotalLTE() int {
	return m.ltecount;
}

func (m *Module_stats) GetDataLTE() float64 {
	return m.datalte;
}

func (m *Module_stats) GetData3G() float64 {
	return m.data3g;
}

func (m *Module_stats) GetTime3G() int {
	return m.time3g;
}

func (m *Module_stats) GetTotalTime() int {
	return m.totaltime;
}

func (m *Module_stats) Reset() bool {

	return true;
}

func (m *Module_stats) GetIdleTime() []int {
	return m.idletime;
}

/* Total trace time in hrs. */
func (m *Module_stats) GetTotalTimeHR() float64 {
	return float64(m.GetTotalTime())/3600.0;
}

func (m *Module_stats) GetFirstAvgDelayTransition() float64 {
	panic("unimplemented method");

	return 0.0;
}
