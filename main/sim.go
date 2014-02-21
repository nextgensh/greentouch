/* Author : Shravan Aras. */
/* Handles some basic simulation functions. */

package main

import ("lib"
		"myio"
		"modules"
		"radio")

/* Read the bandwidth and event files into memory. */
func ReadData(bfile string, bsize int, efile string, esize int) {
	lib.Bandwidths = myio.ReadBandwidthFile(bfile, bsize);
	lib.Events = myio.ReadEventFile(efile, esize);
}

/* Code to start the simulation. */
func StartSimulation(model modules.Module,
						startstate int) (energy lib.Energy,
										radios []radio.IntRadio,
										graphic []lib.Visualize) {

	lib.Bindex = 0;
	lib.Eindex = 0;
	freeze := false;
	predict := 2;
	graphic  = make([]lib.Visualize, len(lib.Bandwidths));
	gindex := 0;
	/* Initialioze the simulation module. */
	model.Init();
	if startstate >= 0 {
		model.SetState(startstate);
	}
	/* Initialize radios used in simulation. */
	radio_lte := radio.Radio{"LTE", lib.Tail_LTE, 0, 0, 0, 0, radio.ACTIVE};
	radio_3g := radio.Radio{"3G", lib.Tail_3G, 0, 0, 0, 0, radio.ACTIVE};
	radios = []radio.IntRadio{&radio_3g, &radio_lte};

	/* Shadow radios for energy-delay2. */
	s_radios := lib.InitRadios();

	/* Replay the traces through the models. */
	for ; lib.Bindex < len(lib.Bandwidths); lib.Bindex ++ {

		if lib.Bandwidths[lib.Bindex].Bandwidth < 0 {
			ed2_lte, ed2_3g := lib.EnergyDelay(s_radios, false);
			if ed2_lte < ed2_3g {
				if printfmt == "firstclick2" {
					startstate = lib.CLTE;
				}
				if printfmt == "firstclick3" {
					lib.BitInc2(predict);
				}
			} else {
				if printfmt == "firstclick2" {
					startstate = lib.C3G;
				}
				if printfmt == "firstclick3" {
					lib.BitDec(predict);
				}
			}
			if printfmt == "firstclick2" {
				model.SetState(startstate);
			}
			if printfmt == "firstclick3" {
				model.SetState(lib.BitBand2(predict));
			}
			model.Reset();
			freeze = false;
			s_radios = lib.InitRadios();
			continue;
		}

		ret := model.ServeData(lib.Bandwidths[lib.Bindex].Bandwidth);
		if !ret {
			continue;
		}
		band, change := model.GetState();
		if change {
			radios[band].ToActive();
			graphic[gindex].AddPoint(lib.Bandwidths[lib.Bindex].Timestamp,
										band);
			gindex ++;
		}
		radios[band].ServeBandwidth(lib.Bandwidths[lib.Bindex].Bandwidth);
		if lib.GetBand(lib.Bandwidths[lib.Bindex].Bandwidth) == lib.CLTE {
			freeze = true;
		}
		if !freeze {
			s_radios[lib.C3G].ServeBandwidth(
							lib.Bandwidths[lib.Bindex].Bandwidth);
			s_radios[lib.CLTE].ServeBandwidth(
							lib.Bandwidths[lib.Bindex].Bandwidth);
		}

		/* Check to see if there is an event coming up next. */
		check, event, eindex := isEvent();
		if check {
			nexttime := int64(0);
			if eindex + 1 >= len(lib.Events) {
				nexttime = lib.Bandwidths[len(lib.Bandwidths)-1].Timestamp;
			} else {
				nexttime = lib.Events[eindex + 1].Timestamp;
			}
			isspike, data, time := lib.IsDataSpike(event.Timestamp,
													nexttime,
													lib.Bindex+1,
													lib.Bandwidths);
			model.HandleEvent(event, data, time, isspike);
		}
	}

	/* Create an energy object and return it. */
	energy = lib.Energy{radio_3g.GetActiveTime() +
						radio_3g.GetTailTime() +
						model.GetDelayTransmission(),
						radio_lte.GetActiveTime() +
						radio_lte.GetTailTime(),
						radio_3g.GetIdleTime(),
						radio_lte.GetIdleTime(),
						model.GetSwitchingEnergy()};

	return energy, radios, graphic;
}

/* Function which tells me if there is an event approaching and returns
 * that event if it is. */
func isEvent() (flag bool, event lib.Event, eindex int) {
	end := lib.Bindex + 1;

	if lib.Eindex > len(lib.Events) {
		return false, lib.Events[0], lib.Eindex;
	}

	if end >= len(lib.Bandwidths) {
		end = len(lib.Bandwidths) - 1;
	}

	count := 0;
	for ; lib.Eindex < len(lib.Events) &&
				lib.Events[lib.Eindex].Timestamp <=
				lib.Bandwidths[end].Timestamp; lib.Eindex++ {
		count ++;
	}

	if count > 0 {
		return true, lib.Events[lib.Eindex-1], lib.Eindex-1;
	}

	return false, lib.Events[0], lib.Eindex;
}
