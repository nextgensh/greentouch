/* Author : Shravan Aras */
/* A wrapper which does some energydelay calculations. */

package lib

import ("radio")

func EnergyDelay(sradio []radio.IntRadio, down bool) (float64, float64) {

	radio_3g := sradio[C3G];
	radio_lte := sradio[CLTE];

	energy := Energy{radio_3g.GetActiveTime() +
					radio_3g.GetTailTime(),
					0,
					radio_3g.GetIdleTime(),
					0,
					0};

	e1 := energy.Get3GEnergy() +
			energy.GetIdleEnergy() +
			Energy_3gtolte;
	if down {
		e1 += Energy_lteto3g;
	}
	t1 :=	Switch_2gto3g +
			radio_3g.GetActiveTime() +
			radio_3g.GetTailTime() +
			radio_3g.GetIdleTime() +
			Switch_3gtolte;
	ed2_3g := e1*t1*t1;

    energy = Energy{0,
                        radio_lte.GetActiveTime() +
                        radio_lte.GetTailTime(),
						0,
                        radio_lte.GetIdleTime(),
                        0};

	e2 := energy.GetLTEEnergy() +
			energy.GetIdleEnergy();
	t2 :=		Switch_2gtolte +
			radio_lte.GetActiveTime() +
			radio_lte.GetTailTime() +
			radio_lte.GetIdleTime();

	ed2_lte := e2*t2*t2;

	return ed2_lte, ed2_3g;
}
