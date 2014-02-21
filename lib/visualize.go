/* Author : Shravan Aras. */
/* Visualization library to use with gnuplot. */

package lib

type Visualize struct {
	time int64;
	radio int;
}

func (v *Visualize) AddPoint(time int64, radio int) {
	v.time = time;
	v.radio = radio;
}

func (v *Visualize) GetPoint() (int64, int) {
	return v.time, (v.radio+1)*1000;
}
