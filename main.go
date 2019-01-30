package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"fmt"
	"go/types"
	"strconv"
	"time"
)

func main() { // Main always gets called as the entry point
	log.SetReportCaller(true)

	app := cli.NewApp()
	app.Name = "UptimeParserGo"
	app.Usage = ""
	app.HideVersion = true
	app.HideHelp = false
	app.EnableBashCompletion = true

	// Setup flags here
	var DebugMode bool
	var IpAddress string
	var Snmp string
	flags := []cli.Flag{
		cli.BoolFlag{

			Name:        "debug, d",
			Usage:       "enable debug mode",
			Destination: &DebugMode,
		},
		cli.StringFlag{

			Name:        "ip, i",
			Usage:       "IP address to scan",
			Destination: &IpAddress,
		},
		cli.StringFlag{

			Name:        "snmp, s",
			Usage:       "SNMP community string",
			Destination: &Snmp,
		},
	}

	// Commands to be run go here, after parsing variables
	app.Commands = []cli.Command{
		{
			UseShortOptionHandling: true,
			Name:    "xml",
			Aliases: []string{"x"},
			Usage:   "export as XML",
			Category: "output",
			Action: func(c *cli.Context) error {

				//output.GenerateXML()

				return nil
			},
		},
		{
			UseShortOptionHandling: true,
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "export as JSON",
			Category: "output",
			Action: func(c *cli.Context) error {
				args := c.Args()
				log.Info(args)

				//DO WORK HERE
				return nil
			},
		},
	}

	app.Flags = flags // Assign flags via parse right before we start work
	app.Before = func(c *cli.Context) error {
		// Actions to run before running parsed commands
		if DebugMode {
			log.SetLevel(5)
			log.Info("Debug Mode")
		} else {
			log.SetLevel(3)
			// open a file
			f, err := os.OpenFile("uptime.log", os.O_CREATE | os.O_RDWR, 0666) // Create new log file every run
			//f, err := os.OpenFile("uptime.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
			if err != nil {
				fmt.Printf("error opening file: %v", err)
			}
			log.SetOutput(f)
			log.Warn("Normal Mode")
		}
		return nil
	}
	// Parse Commands and flags here
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("EOP")
}
type Device struct{
	name string
	up_time types.Nil
	snmp_comm string

}
//Constructor
func (d *Device) NewDevice(Name_in string, Snmp_comm_in string) Device {
	obj := new(Device)
	obj.name = Snmp_comm_in
	obj.snmp_comm = Name_in
	return *obj
}
func (d *Device) UpdateUptime(snmp_port_in int){
	var snmpPort int
	if snmp_port_in == 0{snmpPort = 161
	}else {snmpPort = snmp_port_in}
	log.Debug("SNMP Port:" + strconv.Itoa(snmpPort))


}

func (d *Device) IsOverXHours(overHoursIn int)  {
	var overHours int
	if overHoursIn == 0{overHours = 24
	}else {overHours = overHoursIn}
	log.Debug("Over hour amount: " + strconv.Itoa(overHours))

	//TODO Compare device current time with uptime delta
	t := time.Now()
	log.Debug(t)
}
func GenerateXML()  {
	log.Debug("Start XML generation")
	println("SOME XML")
}