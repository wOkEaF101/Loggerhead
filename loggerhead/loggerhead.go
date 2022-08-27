package main

import (
	"flag"
	loggerhead "loggerhead/internal"
	"sync"
)

//type EventTrace struct {
//	GUID  string
//	Event []uint16
//	PID   []string
//}

/*<opcode name="BmDetection" message="$(string.opcode_BehaviorMonitoringBmDetection)" value="10"/>
<opcode name="BmProcessStart" message="$(string.opcode_BehaviorMonitoringBmProcessStart)" value="11"/>
<opcode name="BmDriverLoad" message="$(string.opcode_BehaviorMonitoringBmDriverLoad)" value="12"/>
<opcode name="BmModuleLoad" message="$(string.opcode_BehaviorMonitoringBmModuleLoad)" value="13"/>
<opcode name="BmDocumentOpen" message="$(string.opcode_BehaviorMonitoringBmDocumentOpen)" value="14"/>
<opcode name="BmFileCreate" message="$(string.opcode_BehaviorMonitoringBmFileCreate)" value="15"/>
<opcode name="BmFileChange" message="$(string.opcode_BehaviorMonitoringBmFileChange)" value="16"/>
<opcode name="BmFileDelete" message="$(string.opcode_BehaviorMonitoringBmFileDelete)" value="17"/>
<opcode name="BmFileRename" message="$(string.opcode_BehaviorMonitoringBmFileRename)" value="18"/>
<opcode name="BmRegistryKeyCreate" message="$(string.opcode_BehaviorMonitoringBmRegistryKeyCreate)" value="19"/>
<opcode name="BmRegistryKeyRename" message="$(string.opcode_BehaviorMonitoringBmRegistryKeyRename)" value="20"/>
<opcode name="BmRegistryKeyDelete" message="$(string.opcode_BehaviorMonitoringBmRegistryKeyDelete)" value="21"/>
<opcode name="BmRegistryValueSet" message="$(string.opcode_BehaviorMonitoringBmRegistryValueSet)" value="22"/>
<opcode name="BmRegistryValueDelete" message="$(string.opcode_BehaviorMonitoringBmRegistryValueDelete)" value="23"/>
<opcode name="BmNetworkConnect" message="$(string.opcode_BehaviorMonitoringBmNetworkConnect)" value="24"/>
<opcode name="BmNetworkData" message="$(string.opcode_BehaviorMonitoringBmNetworkData)" value="25"/>
<opcode name="BmNetworkListen" message="$(string.opcode_BehaviorMonitoringBmNetworkListen)" value="26"/>
<opcode name="BmNetworkAccept" message="$(string.opcode_BehaviorMonitoringBmNetworkAccept)" value="27"/>
<opcode name="BmProcessTerminate" message="$(string.opcode_BehaviorMonitoringBmProcessTerminate)" value="28"/>
<opcode name="BmNetworkDetection" message="$(string.opcode_BehaviorMonitoringBmNetworkDetection)" value="29"/>
<opcode name="BmBootRecordChange" message="$(string.opcode_BehaviorMonitoringBmBootRecordChange)" value="30"/>
<opcode name="BmRemoteThreadCreate" message="$(string.opcode_BehaviorMonitoringBmRemoteThreadCreate)" value="31"/>
<opcode name="BmRegistryBlockSet" message="$(string.opcode_BehaviorMonitoringBmRegistryBlockSet)" value="46"/>
<opcode name="BmRegistryBlockDelete" message="$(string.opcode_BehaviorMonitoringBmRegistryBlockDelete)" value="47"/>
<opcode name="BmRegistryBlockRename" message="$(string.opcode_BehaviorMonitoringBmRegistryBlockRename)" value="48"/>
<opcode name="BmRegistryReplace" message="$(string.opcode_BehaviorMonitoringBmRegistryReplace)" value="49"/>
<opcode name="BmRegistryRestore" message="$(string.opcode_BehaviorMonitoringBmRegistryRestore)" value="50"/>
<opcode name="BmRegistryBlockReplace" message="$(string.opcode_BehaviorMonitoringBmRegistryBlockReplace)" value="51"/>
<opcode name="BmRegistryBlockRestore" message="$(string.opcode_BehaviorMonitoringBmRegistryBlockRestore)" value="52"/>
<opcode name="BmOpenProcess" message="$(string.opcode_BehaviorMonitoringBmOpenProcess)" value="53"/>
<opcode name="BmRegistryBlockCreate" message="$(string.opcode_BehaviorMonitoringBmRegistryBlockCreate)" value="55"/>
<opcode name="BmEtw" message="$(string.opcode_BehaviorMonitoringBmEtw)" value="60"/>
<opcode name="BmFolderCreate" message="$(string.opcode_BehaviorMonitoringBmFolderCreate)" value="61"/>
<opcode name="BmScavengerTask" message="$(string.opcode_BehaviorMonitoringBmScavengerTask)" value="62"/>
<opcode name="BmProcessTainting" message="$(string.opcode_BehaviorMonitoringBmProcessTainting)" value="63"/>
<opcode name="BmFolderRename" message="$(string.opcode_BehaviorMonitoringBmFolderRename)" value="64"/>
<opcode name="BmFolderEnum" message="$(string.opcode_BehaviorMonitoringBmFolderEnum)" value="65"/>
<opcode name="BmFileHardLink" message="$(string.opcode_BehaviorMonitoringBmFileHardLink)" value="66"/>
*/

func init() {
	flag.Bool("p", false, "A boolean value to trigger kernel process event collection.")
	flag.Bool("c", false, "A boolean value to trigger network connection events.")
	flag.Bool("m", false, "A boolean value to trigger Windows Defender ETW events.")
	flag.Bool("d", false, "A boolean value to trigger Microsoft .NET ETW events.")
}

func FlagExists(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func CollectLogs(guid string, event uint16, wg *sync.WaitGroup) {
	defer wg.Done()
	loggerhead.SetLogging("./loggerhead.log")
	coll := loggerhead.EventTrace{GUID: guid, Event: event}
	loggerhead.ETWSession(coll)
}

func main() {
	var wg sync.WaitGroup

	flag.Parse()
	count := flag.NFlag()

	wg.Add(count)

	//Check if command line flags are set to initialize the ETW provider - super gross but ¯\_(ツ)_/¯
	if FlagExists("p") {
		go CollectLogs("{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}", 1, &wg) //Process
	}
	if FlagExists("c") {
		go CollectLogs("{7DD42A49-5329-4832-8DFD-43D979153A88}", 11, &wg) //Network
	}
	if FlagExists("m") {
		go CollectLogs("{0A002690-3839-4E3A-B3B6-96D8DF868D99}", 16, &wg) //Malz
	}
	if FlagExists("d") {
		go CollectLogs("{E13C0D23-CCBC-4E12-931B-D9CC2EEE27E4}", 252, &wg) //.NET
	}

	wg.Wait()
}
