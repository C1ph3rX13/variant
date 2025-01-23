package persistence

import "encoding/xml"

type Task struct {
	XMLName          xml.Name         `xml:"Task"`
	Version          string           `xml:"version,attr"`
	Xmlns            string           `xml:"xmlns,attr"`
	RegistrationInfo RegistrationInfo `xml:"RegistrationInfo"`
	Triggers         Triggers         `xml:"Triggers"`
	Principals       Principals       `xml:"Principals"`
	Settings         Settings         `xml:"Settings"`
	Actions          Actions          `xml:"Actions"`
}

type RegistrationInfo struct {
	Date        string `xml:"Date"`
	Author      string `xml:"Author"`
	Version     string `xml:"Version"`
	Description string `xml:"Description"`
}

type Triggers struct {
	CalendarTrigger *CalendarTrigger `xml:"CalendarTrigger"`
	TimeTrigger     *TimeTrigger     `xml:"TimeTrigger"`
	BootTrigger     *BootTrigger     `xml:"BootTrigger"`
	LogonTrigger    *LogonTrigger    `xml:"LogonTrigger"`
}

type CalendarTrigger struct {
	StartBoundary string        `xml:"StartBoundary"`
	EndBoundary   string        `xml:"EndBoundary"`
	Repetition    Repetition    `xml:"Repetition"`
	ScheduleByDay ScheduleByDay `xml:"ScheduleByDay"`
}

/*
TimeTrigger
时间触发器
<TimeTrigger>
<StartBoundary>2005-10-11T13:21:17-08:00</StartBoundary>
<EndBoundary>2006-01-01T00:00:00-08:00</EndBoundary>
<Enabled>true</Enabled>
<ExecutionTimeLimit>PT5M</ExecutionTimeLimit>
</TimeTrigger>
*/
type TimeTrigger struct {
	StartBoundary      string `xml:"StartBoundary"`
	EndBoundary        string `xml:"EndBoundary"`
	Enabled            bool   `xml:"Enabled"`
	ExecutionTimeLimit string `xml:"ExecutionTimeLimit"`
}

/*
BootTrigger
启动触发器
<BootTrigger>
<StartBoundary>2005-10-11T13:21:17-08:00</StartBoundary>
<EndBoundary>2006-01-01T00:00:00-08:00</EndBoundary>
<Enabled>true</Enabled>
<ExecutionTimeLimit>PT5M</ExecutionTimeLimit>
</BootTrigger>
*/
type BootTrigger struct {
	StartBoundary      string `xml:"StartBoundary"`
	EndBoundary        string `xml:"EndBoundary"`
	Enabled            bool   `xml:"Enabled"`
	ExecutionTimeLimit string `xml:"ExecutionTimeLimit"`
}

/*
LogonTrigger
登录触发器
<LogonTrigger>
<StartBoundary>2005-10-11T13:21:17-08:00</StartBoundary>
<EndBoundary>2006-01-01T00:00:00-08:00</EndBoundary>
<Enabled>true</Enabled>
<UserId>DOMAIN_NAME\UserName</UserId>
</LogonTrigger>
*/
type LogonTrigger struct {
	StartBoundary string `xml:"StartBoundary"`
	EndBoundary   string `xml:"EndBoundary"`
	Enabled       bool   `xml:"Enabled"`
	UserId        string `xml:"UserId"`
}

type Repetition struct {
	Interval string `xml:"Interval"`
	Duration string `xml:"Duration"`
}

type ScheduleByDay struct {
	DaysInterval int `xml:"DaysInterval"`
}

type Principals struct {
	Principal Principal `xml:"Principal"`
}

type Principal struct {
	Id        string `xml:"id,attr"`
	UserId    string `xml:"UserId"`
	RunLevel  string `xml:"RunLevel"`
	LogonType string `xml:"LogonType"`
}

type Settings struct {
	Enabled            bool `xml:"Enabled"`
	AllowStartOnDemand bool `xml:"AllowStartOnDemand"`
	AllowHardTerminate bool `xml:"AllowHardTerminate"`
	Hidden             bool `xml:"Hidden"`
}

type Actions struct {
	Context string `xml:"Context,attr"`
	Exec    Exec   `xml:"Exec"`
}

type Exec struct {
	Command string `xml:"Command"`
}
