package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	ConfigSettings struct {
		AdbPathFolder string `json:"adbpathfolder"`
	} `json:"settings"`

	ConfigFastbootflash struct {
		Boot     string `json:"boot"`
		System   string `json:"system"`
		Vendor   string `json:"vendor"`
		Recovery string `json:"recovery"`
	} `json:"config_fastbootflash"`

	ConfigTwrpflash struct {
		TwrpFilePath  string `json:"twrpfilepath"`
		SystemZipFile string `json:"systemzipfile"`
	} `json:"config_twrpflash"`
}

// initializing JSON struct getting adb path
var config Config

func main() {
	fmt.Println(`

███████████████████████████████████████████████████████████████████████████████████████████
█░░░░░░░░░░░░░░█░░░░░░░░░░░░░░█░░░░░░█████████░░░░░░░░░░░░░░█░░░░░░░░░░░░░░█░░░░░░██░░░░░░█
█░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀░░█████████░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀░░██░░▄▀░░█
█░░░░░░░░░░▄▀░░█░░▄▀░░░░░░░░░░█░░▄▀░░█████████░░▄▀░░░░░░▄▀░░█░░▄▀░░░░░░░░░░█░░▄▀░░██░░▄▀░░█
█████████░░▄▀░░█░░▄▀░░█████████░░▄▀░░█████████░░▄▀░░██░░▄▀░░█░░▄▀░░█████████░░▄▀░░██░░▄▀░░█
█░░░░░░░░░░▄▀░░█░░▄▀░░░░░░░░░░█░░▄▀░░█████████░░▄▀░░░░░░▄▀░░█░░▄▀░░░░░░░░░░█░░▄▀░░░░░░▄▀░░█
█░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀░░█████████░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀▄▀▄▀▄▀▄▀░░█
█░░▄▀░░░░░░░░░░█░░▄▀░░░░░░░░░░█░░▄▀░░█████████░░▄▀░░░░░░▄▀░░█░░░░░░░░░░▄▀░░█░░▄▀░░░░░░▄▀░░█
█░░▄▀░░█████████░░▄▀░░█████████░░▄▀░░█████████░░▄▀░░██░░▄▀░░█████████░░▄▀░░█░░▄▀░░██░░▄▀░░█
█░░▄▀░░░░░░░░░░█░░▄▀░░█████████░░▄▀░░░░░░░░░░█░░▄▀░░██░░▄▀░░█░░░░░░░░░░▄▀░░█░░▄▀░░██░░▄▀░░█
█░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀░░█████████░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀░░██░░▄▀░░█░░▄▀▄▀▄▀▄▀▄▀░░█░░▄▀░░██░░▄▀░░█
█░░░░░░░░░░░░░░█░░░░░░█████████░░░░░░░░░░░░░░█░░░░░░██░░░░░░█░░░░░░░░░░░░░░█░░░░░░██░░░░░░█
███████████████████████████████████████████████████████████████████████████████████████████

  ░██         ░████         ░████   
░████        ░██ ░██       ░██ ░██  
  ░██       ░██ ░████     ░██ ░████ 
  ░██       ░██░██░██     ░██░██░██ 
  ░██       ░████ ░██     ░████ ░██ 
  ░██        ░██ ░██       ░██ ░██  
░██████ ░██   ░████   ░██   ░████   
                                                 


╔══════════════════════════╗
║  Welcome to 2Flash Tool  ║
╚══════════════════════════╝

You can change all configuration settings in .json file.
If you want to continue with a valid TWRP recovery, please type 1. Otherwise, type 2 to run the process via fastboot.
`)

	// reading JSON configs
	file, err := os.ReadFile("settings.json")
	if err != nil {
		fmt.Println("Error reading settings.json:", err)
		return
	}

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	/* USER CHOOSE STACK */
	reader := bufio.NewReader(os.Stdin)

	for {
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Println("Starting TWRP flash...")
			fmt.Println("TWRP file:", config.ConfigTwrpflash.TwrpFilePath)
			fmt.Println("System zip:", config.ConfigTwrpflash.SystemZipFile)
			twrpFlash()
		case "2":
			fmt.Println("Starting Fastboot flash...")
			fmt.Println("Boot image:", config.ConfigFastbootflash.Boot)
			fmt.Println("System image:", config.ConfigFastbootflash.System)
			fmt.Println("Vendor image:", config.ConfigFastbootflash.Vendor)
			fmt.Println("Recovery image:", config.ConfigFastbootflash.Recovery)
			fastbootFlash()
		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
	/* end USER CHOOSE STACK */

	fmt.Println("Press send key to close...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// case 1: flash twrp + adb sideload custom ROM
func twrpFlash() {
	fmt.Println("TWRP executing...")

	if err := runCommandInFolder("adb", "devices"); err != nil {
		fmt.Println(err)
		return
	}

	if err := runCommandInFolder("adb", "reboot", "bootloader"); err != nil {
		fmt.Println(err)
		return
	}

	countdown(20)

	if err := runCommandInFolder("fastboot", "flash", "recovery", config.ConfigTwrpflash.TwrpFilePath); err != nil {
		fmt.Println(err)
		return
	}

	countdown(20)

	if err := runCommandInFolder("fastboot", "boot", config.ConfigTwrpflash.TwrpFilePath); err != nil {
		fmt.Println(err)
		return
	}

	countdown(20)

	if err := runCommandInFolder("adb", "sideload", config.ConfigTwrpflash.SystemZipFile); err != nil {
		fmt.Println(err)
		return
	}

	// wait more time for sideload operations: 8 min.
	countdown(480)

	if err := runCommandInFolder("adb", "reboot"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Process successfully complete")
}

// case 2: direct flashing of syst images
func fastbootFlash() {
	fmt.Println("Fastboot executing...")

	// Step 1: check device connection
	if err := runCommandInFolder("adb", "devices"); err != nil {
		fmt.Println(err)
		return
	}

	// Step 2: reboot to bootloader
	if err := runCommandInFolder("adb", "reboot", "bootloader"); err != nil {
		fmt.Println(err)
		return
	}

	// Wait for bootloader to be ready
	countdown(20)

	// Step 3: flash partitions if paths are set in config
	if config.ConfigFastbootflash.Boot != "" {
		fmt.Println("Flashing boot:", config.ConfigFastbootflash.Boot)
		if err := runCommandInFolder("fastboot", "flash", "boot", config.ConfigFastbootflash.Boot); err != nil {
			fmt.Println(err)
			return
		}
	}

	if config.ConfigFastbootflash.System != "" {
		fmt.Println("Flashing system:", config.ConfigFastbootflash.System)
		if err := runCommandInFolder("fastboot", "flash", "system", config.ConfigFastbootflash.System); err != nil {
			fmt.Println(err)
			return
		}
	}

	if config.ConfigFastbootflash.Vendor != "" {
		fmt.Println("Flashing vendor:", config.ConfigFastbootflash.Vendor)
		if err := runCommandInFolder("fastboot", "flash", "vendor", config.ConfigFastbootflash.Vendor); err != nil {
			fmt.Println(err)
			return
		}
	}

	if config.ConfigFastbootflash.Recovery != "" {
		fmt.Println("Flashing recovery:", config.ConfigFastbootflash.Recovery)
		if err := runCommandInFolder("fastboot", "flash", "recovery", config.ConfigFastbootflash.Recovery); err != nil {
			fmt.Println(err)
			return
		}
	}

	// Step 4: reboot
	if err := runCommandInFolder("fastboot", "reboot"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Fastboot process completed successfully ✅")
}

// run adb/fastboot commands
func runCommandInFolder(name string, args ...string) error {
	exeName := name

	if os.PathSeparator == '\\' && !strings.HasSuffix(name, ".exe") {
		exeName += ".exe"
	}

	// Build the full path
	fullPath := filepath.Join(config.ConfigSettings.AdbPathFolder, exeName)

	cmd := exec.Command(fullPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing %s: %v [in path: %s]", fullPath, err, fullPath)
	}
	return nil
}

func countdown(seconds int) {
	fmt.Println()

	barLength := 30 // lunghezza della barra di progresso
	for i := seconds; i > 0; i-- {
		progress := (seconds - i) * barLength / seconds
		bar := strings.Repeat("/", progress) + strings.Repeat(" ", barLength-progress)

		fmt.Printf("\r[%s] %d sec remaining...", bar, i)
		time.Sleep(1 * time.Second)
	}
	// barra piena al termine
	fmt.Printf("\r[%s] 0 sec remaining... Done!\n", strings.Repeat("/", barLength))
	fmt.Println("Starting next operation...\n")
}
