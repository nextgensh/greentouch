GreenTocuh source code -

Link to the paper - https://dl.acm.org/citation.cfm?id=2971660

Paper Abstract - Smartphones come equipped with multiple radios for cellular data communication such as 4G LTE, 3G, and 2G, that offer different bandwidths and power profiles. 4G LTE offers the highest bandwidth and is desired by users as it offers quick response while browsing the Internet, streaming media, or utilizing numerous network aware applications available to users. However, majority of the time this high bandwidth level is unnecessary, and the bandwidth demand can be easily met by 3G radios at a reduced power level. While 2G radios demand even lower power, they do not offer adequate bandwidth to meet the demand of interactive applications; however, the 2G radio may be utilized to provide connectivity when the phone is in the standby mode. To address different demands for bandwidth, we propose GreenTouch, a system that dynamically adapts to the bandwidth demand and system state by switching between 4G LTE, 3G, and 2G with the goal of minimizing delays and maximizing energy efficiency. GreenTouch associates users' behavior to network activity through capturing and correlating user interactions with the touch display. We have used top applications on the Google play store to show the potential of GreenTouch to reduce energy consumption of the radios by 10%, on average, compared to running the applications in the standard Android. This translates to an overall energy savings of 7.5% for the entire smartphone.

radio/	- This package is used to simulate the state diagram of a cellular radio.
		The FSM related to this radio has 3 states - Active, Idle and Tail.
		We do not simulate the microstate transitions in Active and Idle, but treat them
		as a whole.
		This interface is used to create a specific radios such as LTE or 3G by populating 
		it with the correct state transition values.

lib/	- This package primarily contains the energydelay^2 calculation code, a metric 
		we use in the paper to determine the most optimal radio to pick.
		The package also contains functions to calculate energy, organize date for plotting,
		as well as calculate the accuracy of the simulator for various models.

modules/ - This package contains various prediction strategies used in the paper to
		predict the selection of the correct cellular radio. Not all strategies made it to the
		final paper. This is a good place to look at if you want to see all the various 
		strategies we explored before picking the best ones.

myio/	- Helper functions to read the trace file. The bandwidth file contained 
		bandwidth related information from the application traces while the Events file
		contained user interaction events. (These events where obtained by instrumenting
		the android framework, so log events as user clicked on various elements on the screen).
