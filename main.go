package main

import (
	"fmt"
	"github.com/beevik/etree"
	log "github.com/sirupsen/logrus"
	snmp "github.com/soniah/gosnmp"
	"github.com/urfave/cli"
	"net"
	"os"
	"strconv"
	"time"
)

func main() { // Main always gets called as the entry point
	log.SetReportCaller(true) // Enable long format logging

	// Setup CLI Info
	app := cli.NewApp()
	app.Name = "UptimeParserGo"
	app.Usage = ""
	app.HideVersion = true
	app.HideHelp = false
	app.EnableBashCompletion = true

	// Setup flags here
	var DebugMode bool
	var CidrIpAddress_in string
	var Snmp_in string
	var overTime_in int

	flags := []cli.Flag{
		cli.BoolFlag{

			Name:        "debug, d",
			Usage:       "Enable debug mode",
			Destination: &DebugMode,
		},
		cli.StringFlag{

			Name:        "ip, i",
			Usage:       "IP address to scan",
			Destination: &CidrIpAddress_in,
		},
		cli.StringFlag{

			Name:        "snmp, s",
			Usage:       "SNMP community string ('public' default)",
			Destination: &Snmp_in,
		},
		cli.IntFlag{

			Name:        "overTime, o",
			Usage:       "Time over threshold (24hr default)",
			Destination: &overTime_in,
		},
	}

	// Commands to be run go here, after parsing variables
	app.Commands = []cli.Command{
		{
			UseShortOptionHandling: true,
			Name:     "xml",
			Aliases:  []string{"x"},
			Usage:    "Export as XML",
			Category: "output",
			Action: func(c *cli.Context) error {

				// Set defaults
				var CidrIpAddress_out string
				if CidrIpAddress_in == "" {
					println("IP Address missing from args")
					os.Exit(1)
				} else {
					CidrIpAddress_out = CidrIpAddress_in
				}

				var overTime_out int
				if overTime_in == 0 {
					overTime_out = 24
				} else {
					overTime_out = overTime_in
				}

				var Snmp_out string
				if Snmp_in == "" {
					Snmp_out = "public"
				} else {
					Snmp_out = Snmp_in
				}

				// Pass CLI args and run main logic method
				output := MainLogic(CidrIpAddress_out, Snmp_out, overTime_out)
				print(output)
				return nil
			},
		},
		{
			UseShortOptionHandling: true,
			Name:     "json",
			Aliases:  []string{"j"},
			Usage:    "export as JSON",
			Category: "output",
			Action: func(c *cli.Context) error {
				args := c.Args()
				log.Debug(args)
				log.Fatal("NOT Implemented yet :(")

				//DO WORK HERE
				return nil
			},
		},
	}

	app.Flags = flags // Assign flags via parse right before we start work
	app.Before = func(c *cli.Context) error {
		// Actions to run before running parsed commands
		if DebugMode {
			log.SetLevel(5) // 5=Debug 4=Info
			log.Info("Debug Mode")
		} else {
			log.SetLevel(3)
			// open a file
			f, err := os.OpenFile("uptime.log", os.O_CREATE|os.O_RDWR, 0666) // Create new log file every run
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
	log.Debug("EOP")
}

func MainLogic(ip_CIDR_in string, snmp_in string, overTime int) string {
	// Setup output variables ahead of time
	var outputToConsole string

	// Generate list of IP addresses from console input
	var IPlist, err = Hosts(ip_CIDR_in)
	if err != nil {
		log.Fatal("Invalid IP")
	}

	// Generate list of devices
	log.Debug(IPlist)
	var device_list []Device
	for _, i := range IPlist {
		x := Device{}
		x.name = i
		x.snmp_comm = snmp_in
		device_list = append(device_list, x)
	}

	// Update list of devices
	device_list = UpdateDeviceObjUptimeList(device_list)

	// Generate device over limit count
	var device_over_limit int
	var sec_in_day = uint(60 * 60 * overTime)
	for _, i := range device_list {
		if i.up_time_sec > sec_in_day {
			log.Debug("Overtime - ", i.name)
			device_over_limit = device_over_limit + 1
		}

	}

	var XML_output = map[string]int{} //Key value pairs (like a dictionary)
	XML_output["Up Device count"] = len(device_list)
	XML_output["Device over time limit"] = device_over_limit
	outputToConsole = GenerateXML(XML_output, "Ok")
	return outputToConsole
}

func UpdateDeviceObjUptimeList(device_list_in []Device) []Device {
	// Make list of channels
	log.Info("Creating channels")
	var chan_list []chan Device
	for range device_list_in {
		dev_chan := make(chan Device)
		chan_list = append(chan_list, dev_chan)
	}

	// Dispatch work to threads
	log.Info("Starting work...")
	for i, item := range device_list_in {
		go Update_single_device_uptime(item, chan_list[i])
	}

	// Generate list of devices for output
	log.Info("Merging channels...")
	var device_list_proc []Device
	for i, item := range chan_list {
		log.Debug("Received: ", i)
		device_list_proc = append(device_list_proc, <-item)
		close(item)
	}

	// Print Debug output
	log.Info("Cleaning output...")
	var device_list_out []Device
	for _, item := range device_list_proc {
		if item.up_time_sec != 0 {
			device_list_out = append(device_list_out, item)
			log.Info("Name: " + item.name + " Uptime:" + fmt.Sprint(item.up_time_sec))
		}
	}
	return device_list_out
}

// For Multithread processing
func Update_single_device_uptime(device_in Device, out_dev chan Device) {
	device_in.UpdateUptime()
	out_dev <- device_in
}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var mask = cidr[len(cidr)-2:]
	if mask == "32" {
		var ips []string
		var singleIp = cidr[:len(cidr)-3]
		ips = append(ips, singleIp)
		return ips, nil
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); host_inc(ip) {
		ips = append(ips, ip.String())
	}

	return ips[1 : len(ips)-1], nil // pop network address and broadcast address
}

func host_inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

//GenerateXML Takes key value pairs and output XML that PRTG can ingest
func GenerateXML(data_in map[string]int, msg_in string) string {
	doc := etree.NewDocument()
	prtg := doc.CreateElement("prtg")

	for k, v := range data_in {
		result := prtg.CreateElement("result")
		chan_ele := result.CreateElement("channel")
		chan_ele.CreateText(k)
		val_ele := result.CreateElement("value")
		val_ele.CreateText(strconv.Itoa(v))
	}

	text := prtg.CreateElement("text")
	text.CreateText(msg_in)

	doc.Indent(0)
	XmlOutput, _ := doc.WriteToString()
	return XmlOutput
}

func GenerateJSON() {

}

type Device struct {
	//TODO maybe change uptime to a signed int and use -1 to indicate error?
	name        string
	up_time_sec uint // Value of 0 == bad lookup
	snmp_comm   string
}

//Constructor
func (d *Device) NewDevice(Name_in string, Snmp_comm_in string) Device {
	obj := new(Device)
	obj.name = Snmp_comm_in
	obj.snmp_comm = Name_in
	return *obj
}
func (d *Device) UpdateUptime() {
	d.up_time_sec = d.GetSNMP().(uint) / 100
}

//GetSNMP Defaults to getting up time if no oid is specified
func (d *Device) GetSNMP(oid_in ...string) interface{} {
	var oids []string
	if oid_in == nil {
		oids = []string{"1.3.6.1.2.1.1.3.0"}
	} else {
		oids = oid_in
	}

	// build our own GoSNMP struct, rather than using snmp.Default
	params := &snmp.GoSNMP{
		Target: d.name,
		Port:   161, // When trying to pass a uint16 or convert from int to uint16,
		// the call freezes, just going to hard code it
		Version:   snmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Community: d.snmp_comm,
	}
	log.Debug("Trying: " + d.name)
	err := params.Connect()
	if err != nil {
		log.Debug("Connect() err: %v", err) // Normally this would log as a FATAL
	}
	defer params.Conn.Close()

	result, err2 := params.Get(oids) // Get() accepts up to snmp.MAX_OIDS

	if err2 != nil {
		log.Debug("Get() err: ", d.name, " - ", err2)
	}

	if err != nil || err2 != nil {
		return uint(0) // Normally we would return a better value, but we will deal with it up stream
	} else {
		log.Debug("Result: ", d.name, " ", result.Variables[0].Value)
		return result.Variables[0].Value
	}

}
