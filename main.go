package main

import (
	"fmt"
)

var GoWhatsClient IGoWhatsBot

func main() {

	loadConfig()

	GoWhatsClient = NewGoWhatsBot(WhatsConfig)
	GoWhatsClient.AddEventHandler(pluginMonitor)
	// GoWhatsClient.AddEventHandler(plug)

	var driver = GoWhatsClient.GetConfig("driver")
	var address = GoWhatsClient.GetConfig(driver)

	fmt.Println("Driver :", driver)
	fmt.Println("Address :", address)
	fmt.Println("Client Name :", *GoWhatsClient.GetDeviceProps().Os)
	fmt.Println("Client PlatformType :", GoWhatsClient.GetDeviceProps().PlatformType.String())
	fmt.Println("Client Version :", GoWhatsClient.GetDeviceProps().Version)

	loadPlugins()

	WhatsConfig.SaveToFile(ConfigPath)

	GoWhatsClient.Run()
	GoWhatsClient.Stop()
}
