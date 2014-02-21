/* Author: Shravan Aras. */
/* Code to manage / calculate energy. */

package lib

type Energy struct {
	TimeActive3G float64;
	TimeActivelte float64;
	TimeIdle3G float64;
	TimeIdlelte float64;
	EnergySwitching float64;
}

func (e Energy) Get3GEnergy() float64 {
	return (e.TimeActive3G * Power_active_3g);
}

func (e Energy) GetLTEEnergy() float64 {
	return (e.TimeActivelte * Power_active_lte);
}

func (e Energy) GetIdleEnergy() float64 {
	return (e.TimeIdle3G * Power_idle_3g) +
			(e.TimeIdlelte * Power_idle_lte);
}

func (e Energy) GetSwitchingEnergy() float64 {
	return e.EnergySwitching;
}

func (e Energy) TotalEnergy() float64 {
	return e.Get3GEnergy() + e.GetLTEEnergy() +
			e.GetIdleEnergy() +
			e.GetSwitchingEnergy();
}

func (e Energy) PerGet3GEnergy() float64 {
	return (e.TimeActive3G * Power_active_3g) /
			e.TotalEnergy();
}

func (e Energy) PerGetLTEEnergy() float64 {
	return (e.TimeActivelte * Power_active_lte) /
			e.TotalEnergy();
}

func (e Energy) PerGetIdleEnergy() float64 {
	return ((e.TimeIdle3G * Power_idle_3g) +
			(e.TimeIdlelte * Power_idle_lte)) /
				e.TotalEnergy();
}

func (e Energy) PerGetSwitchingEnergy() float64 {
	return e.EnergySwitching / e.TotalEnergy();
}


