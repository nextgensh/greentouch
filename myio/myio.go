/* Author : Shravan Aras */

package myio

import(	"lib"
		"io"
		"bufio"
		"strconv"
		"strings"
		"os"
	  )


/* Code used to read the bandwidth file. */
func ReadBandwidthFile(filename string, length int) (bandwidth []lib.Bandwidth){
    bandwidth = make([]lib.Bandwidth, length);

    /* Read data from the disk files first into arrays. */
    fbandwidth, _ := os.Open(filename);
    rbandwidth := bufio.NewReader(fbandwidth);

    /* Read contents from the bandwidth file. */
    line, iseof := rbandwidth.ReadString(0xA);
    for count:=0; iseof != io.EOF; count++ {
        buf := strings.Split(line, " ");
        if len(buf) > 1 {
            bandwidth[count].Timestamp, _ = strconv.ParseInt(buf[0], 10, 64); 
			if len(buf) > 2 {
				bandwidth[count].Bandwidth, _ = strconv.ParseFloat(buf[1], 64); 
				buf[2] = strings.Split(buf[2],"\n" )[0];
				temp, _ := strconv.ParseFloat(buf[2], 64); 
				bandwidth[count].Bandwidth += temp;
			} else {
				buf[1] = strings.Split(buf[1],"\n" )[0];
				bandwidth[count].Bandwidth, _ = strconv.ParseFloat(buf[1], 64); 
			}
        }    
        line, iseof = rbandwidth.ReadString(0xA);
    }    

    return bandwidth;
}

/* Code used to read the event file. */ 
func ReadEventFile(filename string, length int) (event []lib.Event){
    event = make([]lib.Event, length);

    /* Read data from the disk files first into arrays. */
    fevent, _ := os.Open(filename);
    revent := bufio.NewReader(fevent);

    /* Read contents from the events file. */
    line, iseof := revent.ReadString(0xA);
    for count:=0; iseof != io.EOF; count++ {
        buf := strings.Split(line, " ");
        event[count].Timestamp, _ = strconv.ParseInt(buf[0], 10, 64);
        buf[1] = strings.Split(buf[1],"\n" )[0];
        event[count].Name = buf[1];
        line, iseof = revent.ReadString(0xA);
    }

    return event;
}
