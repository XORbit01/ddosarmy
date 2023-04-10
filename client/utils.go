package client

import "github.com/fatih/color"

func DisplayCampInfo(cp CampAPI) {
	color.Yellow("VICTIM : " + cp.Settings.VictimServer)
	color.Green("LEADER : " + cp.Leader.Name)
	color.Blue("SOLDIERS : \n")
	for _, s := range cp.Soldiers {
		color.Cyan("   " + s.Name + "  : " + s.Ip)
	}
	color.Red("DDOS TYPE : " + cp.Settings.DDOSType)
	color.Red("STATUS : " + cp.Settings.Status)
}
