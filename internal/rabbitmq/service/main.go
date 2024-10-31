package service

import (
	"encoding/json"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	customLog "golang-api-starter/internal/modules/log"
	"golang-api-starter/internal/notification/email"
)

type RbmqWorker struct{}

func (rw *RbmqWorker) HandleLogFromQueue(msgBody []byte) {
	var log = new(customLog.Log)
	if err := json.Unmarshal(msgBody, log); err != nil {
		logger.Errorf("failed to Unmarshal log, err: %+v", err)
	}

	logger.Infof("Received log: %v bytes\n", *log)

	// create log to database,
	customLog.Srvc.Create([]*customLog.Log{log})
}

func (rw *RbmqWorker) HandleEmailFromQueue(msgBody []byte) {
	var emailInfo = new(email.EmailInfo)
	if err := json.Unmarshal(msgBody, emailInfo); err != nil {
		logger.Errorf("failed to Unmarshal email, err: %+v", err)
	}

	logger.Infof("Received email: %v bytes\n", *emailInfo)

	// send email
	if err := email.TemplateEmail(*emailInfo); err != nil {
		logger.Errorf("failed to send email in rbmq worker, err: %+v", err)
	}
}
