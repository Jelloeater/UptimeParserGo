package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"fmt"
	"go/types"
	"strconv"
	"time"
	snmp "github.com/soniah/gosnmp"
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

	x := Device{}
	x.snmp_comm = "public"
	x.name = "192.168.1.1"
	x.GetSNMP()
	log.Info("EOP")
}




type Device struct{
	name string
	up_time types.Nil
	snmp_comm string
	snmp_port int

}
//Constructor
func (d *Device) NewDevice(Name_in string, Snmp_comm_in string) Device {
	obj := new(Device)
	obj.name = Snmp_comm_in
	obj.snmp_comm = Name_in
	obj.snmp_port = 161
	return *obj
}
func (d *Device) UpdateUptime(){
	log.Debug("SNMP Port")



}

//GetSNMP Defaults to getting up time if no oid is specified
func (d *Device) GetSNMP(oid_in ...string)interface{}{
	var oids []string
	if oid_in ==nil {
		oids = []string{"1.3.6.1.2.1.1.3.0"}
	}else {oids = oid_in}

	// build our own GoSNMP struct, rather than using snmp.Default
	params := &snmp.GoSNMP{
		Target:    d.name,
		Port:      161, // When trying to pass a uint16 or convert from int to uint16,
						// the call freezes, just going to hard code it
		Version:   snmp.Version2c,
		Timeout:   time.Duration(30) * time.Second,
		Community: d.snmp_comm,
	}
	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	result, err2 := params.Get(oids) // Get() accepts up to snmp.MAX_OIDS

	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	log.Debug(result.Variables[0].Value)
	return result.Variables[0].Value


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

type MainLogic struct {

}

func (MainLogic) MainLogic()  {

}

func (MainLogic) UpdateDeviceObjUptime(){

}

func (MainLogic) GenerateSensorData (){

}

func (MainLogic) GenerateXML (){

}

func (MainLogic) GenerateJSON()  {

}