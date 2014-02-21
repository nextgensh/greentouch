
package radio

/* 3 states a radio interface can be in. */
const ACTIVE = 0;
const TAIL = 1;
const IDLE = 2;

/* Interfaces which maintain the state of a radio. */
type IntRadio interface {
    ServeBandwidth(bandwidth float64);
    GetActiveTime() float64;
    GetTailTime() float64;
    GetIdleTime() float64;
    ToActive();
}

type Radio  struct {
    Name string;
    Tailtime float64;
    /* Time spent by this interface in active, tail and idle time. */
    Active_time float64;
    Tail_time float64;
    Tail_seen float64;
    Idle_time float64;
    /* The current state this radio is in ACTIVE, TAIL, IDLE. */
    State int;
}

func (r *Radio) ServeBandwidth(bandwidth float64) {

    if bandwidth > 0.0 {
        r.State = ACTIVE;
        r.Active_time ++;
        r.Tail_seen = 0.0;
        return;
    }

    switch r.State {
        case ACTIVE:
            r.State = TAIL;
            r.Tail_seen ++;
        case TAIL:
            r.Tail_seen ++;
            r.Tail_time ++;
            if r.Tail_seen >= r.Tailtime {
                /* Reset the tail seen internal counter. */
                r.Tail_seen = 0.0;
                r.State = IDLE;
                return;
            }
        case IDLE:
            r.Idle_time ++;
    }
}

func (r *Radio) ToActive() {
    r.State = ACTIVE;
    r.Tail_seen = 0.0;
}

func (r *Radio) GetActiveTime() float64 {
    return r.Active_time;
}

func (r *Radio) GetTailTime() float64 {
    return r.Tail_time;
}

func (r *Radio) GetIdleTime() float64 {
    return r.Idle_time;
}
