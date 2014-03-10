package main

import (
		"fmt"
		"modules"
		"os"
		"strconv"
		"lib"
		)

var printfmt string;

func main(){

	args := os.Args;

    if len(args) < 7 {
        fmt.Println("Usage : simulator  <eventfile> <length> <bandwidthfile> <length> <energy / delay> <start>");
        return;
    }

    event_file := args[1];
    event_length, _ := strconv.Atoi(args[2]);
    bandwidth_file := args[3];
    bandwidth_length, _ := strconv.Atoi(args[4]);
	printfmt = args[5];
	start := args[6];

	startstate := lib.C3G;

	if start == "lte" {
		startstate = lib.CLTE;
	} else {
		startstate = lib.C3G;
	}

	ReadData(bandwidth_file, bandwidth_length, event_file, event_length);

	modlte := modules.Module_lte{};
	modreactive := modules.Module_RO{};
	modgto := modules.Module_GTO{};
	modgto1 := modules.Module_GTO1{};
	modgto2 := modules.Module_GTO2{};
	modgto2h := modules.Module_GTO2H{};
	modoracle := modules.Module_O{};
	mod3g := modules.Module_3g{};
	modgto4s := modules.Module_GTO4S{};
	modgto4a := modules.Module_GTO4A{};
	modgto3 := modules.Module_GTO3{};
	modgto4 := modules.Module_GTO4{};
	modgto5 := modules.Module_GTO5{};
	modgto1h := modules.Module_GTO1H{};
	modgto3h := modules.Module_GTO3H{};
	modgto4h := modules.Module_GTO4H{};
	modgto5h := modules.Module_GTO5H{};

	arr_modules := []modules.Module{&modlte, &modreactive,
									&modgto, &modgto1,
									&modgto2, &modgto2h,
									&modoracle, &mod3g,
									&modgto4s, &modgto4a,
									&modgto3, &modgto4,
									&modgto5, &modgto1h,
									&modgto3h, &modgto4h,
									&modgto5h};

	arr_energy := make([]lib.Energy, len(arr_modules));

	arr_names := []string{"LTE", "Reactive", "GTAverage", "GTO1", "GTO2",
							"GTO2H", "Oracle", "3G", "GTS",
							"GTA", "GTO3", "GT", "GTO5",
							"GTO1H", "GTO3H", "GTO4H", "GTO5H"};

	arr_order := []int{0, 1, 5, 6, 7, 8, 9, 3, 4, 10, 11, 12, 13, 14, 15, 16};
	for a:=0; a < len(arr_order); a++ {
		i := arr_order[a];
		if arr_names[i] == "LTE" || arr_names[i] == "3G" {
			arr_energy[i], _, _ = StartSimulation(arr_modules[i], -1);
		} else {
			arr_energy[i], _, _ = StartSimulation(arr_modules[i], startstate);
		}
	}

	if printfmt == "spiketime" {
		fmt.Printf("%.4f\n", (modgto4.GetAvgSpikeTime()));
	}

	arr_order = []int{0, 7, 6, 1, 11};
	if printfmt == "energy" {
		ltetotal := arr_energy[0].TotalEnergy();
		for a:=0; a < len(arr_order); a++ {
			i := arr_order[a];
			norm := arr_energy[i].TotalEnergy() /
					ltetotal;
			fmt.Printf("%s %.4f %.4f %.4f %.4f\n", arr_names[i],
							arr_energy[i].PerGetLTEEnergy() * norm,
							arr_energy[i].PerGet3GEnergy() * norm,
							arr_energy[i].PerGetIdleEnergy() * norm,
							arr_energy[i].PerGetSwitchingEnergy() * norm);
		}
		fmt.Println("");
	}

	arr_order = []int{7, 1, 11};
	if printfmt == "delay" {
		for a:=0; a < len(arr_order); a++ {
			i := arr_order[a];
			if i != 11 {
				fmt.Printf("%s %.4f\n", arr_names[i],
							arr_modules[i].GetAvgDelayTransition()+
							arr_modules[i].GetAvgDelayTransmission());
			} else {
				fmt.Printf("%s %.4f %.4f\n", arr_names[i],
							arr_modules[i].GetAvgDelayTransition()+
							arr_modules[i].GetAvgDelayTransmission(),
							modgto4.GetAvgDelayLearnTransition());
			}
		}
		fmt.Println("");
	}

	if printfmt == "delaywaste" {
		fmt.Printf("%.4f\n", modgto4.GetAvgDelayWasteTransition());
	}

	if printfmt == "compdelay" {
		fmt.Printf("%.4f %.4f\n", modreactive.GetDelayTransition(),
					modgto4.GetDelayTransition() +
					modgto4.GetDelayWasteTransition());
		fmt.Println("");
	}

	arr_order = []int{2,4,5};
	if printfmt == "accuracy1" {
		for a:=0; a < len(arr_order); a++ {
			//wrong := arr_modules[arr_order[a]].GetMissed() +
			//			arr_modules[arr_order[a]].GetUnnecesary();
			//correct := arr_modules[arr_order[a]].GetTotal() - wrong;
			//norm := arr_modules[arr_order[a]].GetTotal();
			fmt.Printf("%d %d %d \n",
				arr_modules[arr_order[a]].GetTotalLTE(),
				arr_modules[arr_order[a]].GetTotalLTE()-
				arr_modules[arr_order[a]].GetMissed(),
				arr_modules[arr_order[a]].GetMissed());
		}
		fmt.Println("");
	}

	if printfmt == "accuracy2" {
		for a:=0; a < len(arr_order); a++ {
			//wrong := arr_modules[arr_order[a]].GetMissed() +
			//			arr_modules[arr_order[a]].GetUnnecesary();
			//correct := arr_modules[arr_order[a]].GetTotal() - wrong;
			//norm := arr_modules[arr_order[a]].GetTotal();
			fmt.Printf("%d %d \n",
						arr_modules[arr_order[a]].GetTotal() -
						arr_modules[arr_order[a]].GetTotalLTE(),
						arr_modules[arr_order[a]].GetUnnecesary());
		}
		fmt.Println("");
	}


	if printfmt == "datastats" {
		modstats := modules.Module_stats{};
		_, _, _ = StartSimulation(&modstats, -1);
		fmt.Printf("%.4f %.4f\n", modstats.GetDataLTE(),
									modstats.GetData3G());
	}

	if printfmt == "tracestats" {
		modstats := modules.Module_stats{};
		_, _, _ = StartSimulation(&modstats, -1);
		fmt.Printf("%.4f %d %d %.4f %d %d\n", modstats.GetTotalTimeHR(),
							modstats.GetTotalLTE(),
							modstats.GetTotal()-
							modstats.GetTotalLTE(),
							(modstats.GetDataLTE()+
							modstats.GetData3G())/1024.0,
							modgto4.UniqueLTE(),
							modgto4.Unique3G());
	}



	if printfmt == "timestats" {
		modstats := modules.Module_stats{};
		_, radios, _ := StartSimulation(&modstats, -1);
		fmt.Printf("%d %d %f\n", modstats.GetTotalTime()-
								modstats.GetTime3G(),
								modstats.GetTime3G() -
								int(radios[lib.CLTE].GetIdleTime()),
								radios[lib.CLTE].GetIdleTime());
	}

	if printfmt == "energydelay" {
		//ed2lte, ed23g := EnergyDelay();
		//fmt.Printf("%.4f %.4f\n", ed2lte, ed23g);
		//fmt.Println("");
	}

	arr_order = []int{5};
	if printfmt == "firstclick1" || printfmt == "firstclick2" ||
			printfmt == "firstclick3" {
		for a:=0; a < len(arr_order); a++ {
			fmt.Printf("%.4f\n",
				arr_modules[arr_order[a]].GetFirstAvgDelayTransition());
		}
	}

	arr_order = []int{0, 8, 9, 5, 7};
	if printfmt == "energy3" {
		ltetotal := arr_energy[0].TotalEnergy();
		for a:=0; a < len(arr_order); a++ {
			i := arr_order[a];
			norm := arr_energy[i].TotalEnergy() /
					ltetotal;
			fmt.Printf("%s %.4f %.4f %.4f %.4f\n", arr_names[i],
							arr_energy[i].PerGetLTEEnergy() * norm,
							arr_energy[i].PerGet3GEnergy() * norm,
							arr_energy[i].PerGetIdleEnergy() * norm,
							arr_energy[i].PerGetSwitchingEnergy() * norm);
		}
		fmt.Println("");
	}

	if printfmt == "visualize" {
		_, _, graphic := StartSimulation(&modgto4s, lib.C3G);
		for a:=0; a < len(graphic); a++ {
			t, r := graphic[a].GetPoint();
			if t > 0 {
				fmt.Println(t, r);
		}
		}
	}

	arr_order = []int{3, 4, 10, 11, 12};
	if printfmt == "acccompare" {
		for a:=0; a < len(arr_order); a++ {
			norm := arr_modules[arr_order[a]].GetTotalLTE();
			if norm == 0 {
				norm = 1;
			}
			fmt.Printf("%.4f \n",
				(1-(float64(arr_modules[arr_order[a]].GetMissed())/
				float64(norm)))*100);
		}
		fmt.Println("");
	}
	arr_order = []int{3, 4, 10, 11, 12};
	if printfmt == "acccompare2" {
		for a:=0; a < len(arr_order); a++ {
			norm := arr_modules[arr_order[a]].GetTotal()-
					arr_modules[arr_order[a]].GetTotalLTE();
			if norm == 0 {
				norm = 1;
			}
			fmt.Printf("%.4f \n",
				((float64(arr_modules[arr_order[a]].GetUnnecesary())/
				float64(norm)))*100);
		}
		fmt.Println("");
	}

	arr_order = []int{3, 4, 10, 11, 12};
	if printfmt == "delaycompare" {
		for a:=0; a < len(arr_order); a++ {
			i := arr_order[a];
			fmt.Printf("%.4f\n",
						arr_modules[i].GetAvgDelayTransition()+
						arr_modules[i].GetAvgDelayTransmission());
		}
		fmt.Println("");
	}

	arr_order = []int{0, 8, 9, 7};
	if printfmt == "energy3compare" {
		for a:=0; a < len(arr_order); a++ {
			i := arr_order[a];
			norm := arr_energy[0].TotalEnergy();
			gto := (arr_energy[11].TotalEnergy() / norm) * 100;
			fmt.Printf("%.4f ",
						((arr_energy[i].TotalEnergy() / norm) * 100) -
							gto);
		}
		fmt.Println("");
		fmt.Println("");
	}

	arr_order = []int{5, 10};
	if printfmt == "acccompare3" {
		for a:=0; a < len(arr_order); a++ {
			norm := arr_modules[arr_order[a]].GetTotalLTE();
			if norm == 0 {
				norm = 1;
			}
			fmt.Printf("%.4f \n",
				(1-(float64(arr_modules[arr_order[a]].GetMissed())/
				float64(norm)))*100);
		}
		fmt.Println("");
	}

	arr_order = []int{13, 5, 14, 15, 16};
	if printfmt == "acccompare4" {
		for a:=0; a < len(arr_order); a++ {
			norm := arr_modules[arr_order[a]].GetTotal()-
					arr_modules[arr_order[a]].GetTotalLTE();
			fmt.Printf("%.4f \n",
				((float64(arr_modules[arr_order[a]].GetUnnecesary())/
				float64(norm)))*100);
		}
		fmt.Println("");
	}

	arr_order = []int{13, 5, 14, 15, 16};
	if printfmt == "acccompare5" {
		for a:=0; a < len(arr_order); a++ {
			norm := arr_modules[arr_order[a]].GetTotalLTE();
			if norm == 0 {
				norm = 1;
			}
			fmt.Printf("%.4f \n",
				(1-(float64(arr_modules[arr_order[a]].GetMissed())/
				float64(norm)))*100);
		}
		fmt.Println("");
	}

	if printfmt == "mastercompare" {
		for a:=1; a <= 5; a++ {
			for b:=1; b <=5; b++ {
				mod := modules.Module_GTO2HN{};
				mod.SetHistoryBit(a);
				mod.SetPredictBit(b);
				_, _, _ = StartSimulation(&mod, lib.C3G);
				norm := mod.GetTotalLTE();
				fmt.Printf("%.4f ",
					(1-(float64(mod.GetMissed())/
					float64(norm)))*100);
			}
			fmt.Println("");
		}
		fmt.Println("");
	}

	if printfmt == "mastercompare1" {
		for a:=1; a <= 5; a++ {
			for b:=1; b <=5; b++ {
				mod := modules.Module_GTO2HN{};
				mod.SetHistoryBit(a);
				mod.SetPredictBit(b);
				_, _, _ = StartSimulation(&mod, lib.C3G);
				norm := mod.GetTotal()-
						mod.GetTotalLTE();
				fmt.Printf("%.4f ",
					((float64(mod.GetUnnecesary())/
					float64(norm)))*100);

			}
			fmt.Println("");
		}
		fmt.Println("");
	}

}
