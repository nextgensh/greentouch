/* Author : Shravan Aras. */

package lib

const C3G = 0;
const CLTE = 1;
const Bandwidth_3g = (9.38 * 1024) / 8;
const Bandwidth_lte = (23.51 * 1024) / 8;
const Bandwidth_2g = 6.8;
const Switch_3gtolte = 1587.0 / 1000.0;
const Switch_lteto3g = 7487.0 / 1000.0;
const Switch_2gto3g = 7.46;
const Switch_2gtolte = 7.23;
//const Timeout = 50
const Timeout = 121
const Tail_LTE = 12.56108;
const Tail_3G = 12.56108;
const Power_active_lte = 1.607;
const Power_active_3g = 1.401;
const Power_idle_lte = 0.0766478;
const Power_idle_3g = 0.0548549;
const Energy_3gtolte = 2.425;
const Energy_lteto3g = 8.65;
const Rolling_window_size = 5;
const SpikeMin = 20;
const SpikeWait = 3;
//const SpikeMinTime = 23;
//const SpikeMinTime = 32;
const SpikeMinTime = 121;
const SpikeData = 1024;
