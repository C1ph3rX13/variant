package alive

import (
	"fmt"
	"time"
	"variant/log"

	"github.com/capnspacehook/taskmaster"
)

func SetWinTask(execPath string, taskPath string) {
	// 创建初始化计划任务
	taskService, _ := taskmaster.Connect()
	defer taskService.Disconnect()

	// 定义新的计划任务
	newTaskDef := taskService.NewTaskDefinition()

	// 添加执行程序的路径
	newTaskDef.AddAction(taskmaster.ExecAction{
		Path: "cmd.exe /c start /b",
		Args: execPath,
	})

	// 定义计划任务开始执行的时间
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Local)

	// 基于时间的触发器
	newTaskDef.AddTrigger(
		taskmaster.DailyTrigger{
			TaskTrigger: taskmaster.TaskTrigger{
				Enabled:       true,
				StartBoundary: startTime,
			},
			DayInterval: 1,
		})

	// 基于事件的触发器
	newTaskDef.AddTrigger(
		taskmaster.LogonTrigger{
			TaskTrigger: taskmaster.TaskTrigger{
				Enabled:       true,
				StartBoundary: startTime,
			},
		},
	)

	// 创建计划任务
	_, _, err := taskService.CreateTask(fmt.Sprintf("\\%s", taskPath), newTaskDef, true)
	if err != nil {
		log.Fatal(err)
	}
}
