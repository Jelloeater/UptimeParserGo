package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"fmt"
	"strconv"
	"time"
	snmp "github.com/soniah/gosnmp"
	"github.com/beevik/etree"
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
	var CidrIpAddress string
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
			Destination: &CidrIpAddress,
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

				output := MainLogic(CidrIpAddress,Snmp)
				print(output)
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
	name        string
	up_time_sec uint
	snmp_comm   string

}
//Constructor
func (d *Device) NewDevice(Name_in string, Snmp_comm_in string) Device {
	obj := new(Device)
	obj.name = Snmp_comm_in
	obj.snmp_comm = Name_in
	return *obj
}
func (d *Device) UpdateUptime(){
	d.up_time_sec = d.GetSNMP().(uint)/100



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


func MainLogic(ip_CIDR_in string, snmp_in string)  string{
	var outputToConsole string
	var XML_output = map[string]int{} //Key value pairs (like a dictionary)


	// TODO add list of devices
	//var device_list []Device


	x := Device{}
	x.name = "192.168.1.1"
	x.snmp_comm = snmp_in
	x.UpdateUptime()

	// TODO Take list of devices and update uptime on all of them at once 'multi-threaded'

	// TODO Compare all the sensors and then output the results as XML


	XML_output["Up Device count"]  = 42
	XML_output["Device over time limit"]  = 24
	outputToConsole = GenerateXML(XML_output, "Ok")
	log.Debug("")
	return outputToConsole
}

func UpdateDeviceObjUptime(device_obj_in Device) Device{

	return device_obj_in
}

//GenerateSensorData Takes in a list of devices (slice) and output a dictionary (map key-value pair)
func GenerateSensorData (device_list_in []Device)map[string]int{

	// TODO write go routine, see python script (UptimeParser) for reference in other repo

	return map[string]int{"foo": 1, "bar": 2}
}

//GenerateXML Takes key value pairs and output XML that PRTG can ingest
func GenerateXML (data_in map[string]int,msg_in string)string{
	doc := etree.NewDocument()
	prtg := doc.CreateElement("prtg")
	result := prtg.CreateElement("result")

	for k, v := range data_in{
		chan_ele := result.CreateElement("channel")
		chan_ele.CreateText(k)
		val_ele := result.CreateElement("value")
		val_ele.CreateText(strconv.Itoa(v))
	}

	text := prtg.CreateElement("text")
	text.CreateText(msg_in)

	doc.Indent(0)
	XmlOutput,_ := doc.WriteToString()
	return XmlOutput
}

func GenerateJSON()  {

}