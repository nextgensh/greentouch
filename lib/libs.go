/* Author : Shravan Aras. */
/* A bunch of common library functions. */

package lib

import ("math"
		"radio")

func GetBand(bandwidth float64) int {

    if bandwidth >= Bandwidth_3g {
        return CLTE;
    } else if bandwidth >= Bandwidth_2g {
		return C3G;
    }

	return -1;
}

/* Method used to peak in the future. */
func LookAhead(time int) (size int, bandwidth []Bandwidth){
	temp := time + Bindex;

	if temp > len(Bandwidths) {
		temp = len(Bandwidths);
	}

	return temp - time, Bandwidths[Bindex:temp];
}

func IsWasteful() bool {
	for a:=Bindex; a < Bindex + SpikeMinTime &&
						a < len(Bandwidths); a++ {
		if GetBand(Bandwidths[a].Bandwidth) == CLTE {
			return false;
		}
	}

	return true;
}

/* Method which calculates the energy-delay2 for falling down. */
func EnergyDelayFall(islte bool, data float64) (fall bool) {
	radios := InitRadios();

	a:=Bindex;

	/* If is LTE is there then skip the immediate LTE spike. */
	if islte {
		for ; a < len(Bandwidths) &&
					GetBand(Bandwidths[a].Bandwidth) != CLTE; a++ {

		}
		for ; a < len(Bandwidths) &&
					GetBand(Bandwidths[a].Bandwidth) != C3G; a++ {

		}
	}

	for ; a < len(Bandwidths) &&
		GetBand(Bandwidths[a].Bandwidth) != CLTE; a++ {
		radios[CLTE].ServeBandwidth(Bandwidths[a].Bandwidth);
		radios[C3G].ServeBandwidth(Bandwidths[a].Bandwidth);
	}

	ed2lte, ed23g := EnergyDelay(radios, true);

	if ed23g <= ed2lte {
		return true;
	}

	return false;
}

func ShouldISwitch(data float64, spike_time float64) (jump bool) {

    temp := (Switch_3gtolte - spike_time);

    if temp < 0 {
        temp = 0;
    }

    /* Calculate the time it will take in both configurations. */
    timeinlte := temp + (data / Bandwidth_lte);
    timein3g := data / Bandwidth_3g;

    /* Calculate the energy in both. */
    energyinlte := (data / Bandwidth_lte) * Power_active_lte;
    energyin3g := timein3g * Power_active_3g;

    costlte := (Energy_3gtolte + energyinlte) * (timeinlte * timeinlte);
    cost3g := (energyin3g) * (timein3g * timein3g);

    if(costlte  < cost3g){
        return true;
    }
    return false;
}

func IsDataSpike(time1, time2 int64, b_start int, bandwidth []Bandwidth) (
					status bool, data float64, spike_data_time float64){ 

    spike_data_time  = 0;
    data = 0;
    c_bandwidth := float64(0.0);
    hint_lte := false;  /* Tell me if I have been to the LTE range or not. */
    hint_wcdma := false;
	flag := false;
	timespent := 0;

    for ; b_start < len(bandwidth) && 
					bandwidth[b_start].Timestamp <= time2; b_start++ { 
        c_bandwidth = bandwidth[b_start].Bandwidth;

		timespent ++;

		if !hint_lte && timespent > SpikeMinTime {
			return false, data, spike_data_time;
		}

        data = data + c_bandwidth;
		current_state := GetBand(c_bandwidth);

		if bandwidth[b_start].Bandwidth != 0.0 &&
								!flag {
			flag = true;
			spike_data_time = float64(bandwidth[b_start].Timestamp - 
										time1) / 1000;
		}

        if(current_state == C3G && !hint_wcdma) {
            hint_wcdma = true;
        }

        if(current_state == CLTE && !hint_lte){
            hint_lte = true;
        }

        if ((!hint_lte && hint_wcdma && bandwidth[b_start].Bandwidth == 0.0)) {
			if  timespent <= SpikeMinTime {
				data = 0;
				hint_wcdma = false;
			} else {
				return false, data, spike_data_time;
			}
        }

        if((hint_lte && bandwidth[b_start].Bandwidth == 0.0)){
            /* I am in the accepted state, which means I have hopefuly found a 
             * spike and I can return. */
             return true, data, spike_data_time;
        }
    }

    return false, data, spike_data_time;

}

func GetAverage(data_arr_in []float64, data float64, count int) (
							data_avg float64, data_arr []float64) {
    data_arr_in[(count-1) % Rolling_window_size] = data;
    sum := 0.0;

    /* We cannot calculate the rolling average until we have atleast _rolling_window_  
     * number of elements in it. So until I get that this function just returns the 
     * normal average of the elements. */

    for a:=0; a <= (count-1) % Rolling_window_size ; a++ {
        sum += data_arr_in[a];
    }

    return (sum / math.Min(float64(Rolling_window_size), float64(count))), data_arr_in;
}

/* 00 - 3g
 * 01 - 3g
 * 10 - LTE
 * 11 - LTE
 */
func BitBand(bit int) int {
	if bit >= 4 {
		return CLTE;
	}

	return C3G;
}

func BitBand2(bit int) int {
	if bit >= 2 {
		return CLTE;
	}

	return C3G;
}

func BitBand1(bit int) int {
	if bit >= 1 {
		return CLTE;
	}

	return C3G;
}

func BitBandN(bit int, mark int) int {
	if bit >= mark {
		return CLTE;
	}

	return C3G;
}

func BitInc(bit int) int {
	bit ++;
	if bit > 7 {
		bit = 7;
	}

	return bit;
}

func BitInc2(bit int) int {
	bit ++;
	if bit > 3 {
		bit = 3;
	}

	return bit;
}

func BitInc1(bit int) int {
	bit ++;

	if bit > 1 {
		bit = 1;
	}

	return bit;
}

func BitIncN(bit int, size int) int {
	bit ++;

	if bit > size {
		bit = size;
	}

	return bit;
}

func BitDec(bit int) int {
	bit --;
	if bit < 0 {
		bit = 0;
	}

	return bit;
}

/* Anything above a 2 is jump to lte. */
func BitJump(bit int) bool {
	if BitBand(bit) == CLTE {
		return true;
	}

	return false;
}

func BitJumpN(bit int, mark int) bool {
	if BitBandN(bit, mark) == CLTE {
		return true;
	}

	return false;
}

/* Does the exact opposite of BitJump does. */
func BitFall(bit int) bool {
	if BitBand1(bit) == C3G {
		return true;
	}

	return false;
}

func BitAdd(bit int, state int) int {
	return (((bit << 1) | state) & 0x3);
}

func BitAddN(bit int, state int, size int) int {
	return (((bit << 1) | state) & size);
}

func InitRadios() []radio.IntRadio {
    s_radio_lte := radio.Radio{"LTE", Tail_LTE, 0, 0, 0, 0, radio.ACTIVE};
    s_radio_3g := radio.Radio{"3G", Tail_3G, 0, 0, 0, 0, radio.ACTIVE};
    s_radios := []radio.IntRadio{&s_radio_3g, &s_radio_lte};

    return s_radios;
}

func BitValue(bits int) int {
	return int(math.Pow(2.0, float64(bits)));
}
