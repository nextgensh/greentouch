package lib

type Event struct {
    Timestamp int64;
    Name string;
}

type Bandwidth struct {
    Timestamp int64;
    Bandwidth float64;
}

var Bandwidths []Bandwidth;
var Events []Event;

/* Points to the point in the trace were we are currently looking at. */
var Bindex int;
var Eindex int;
