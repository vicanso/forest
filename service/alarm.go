// Copyright 2020 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"crypto/tls"
	"sync"

	"github.com/vicanso/forest/config"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var (
	mailDialer = newMailDialer()
	mailSender string

	sendingMailMutex = new(sync.Mutex)

	basicInfo   config.BasicConfig
	alarmConfig config.AlarmConfig
)

func newMailDialer() *gomail.Dialer {
	basicInfo = config.GetBasicConfig()
	alarmConfig = config.GetAlarmConfig()
	mailConfig := config.GetMailConfig()
	if mailConfig.Host == "" {
		return nil
	}
	mailSender = mailConfig.User
	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.User, mailConfig.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d
}

// AlarmError 发送出错警告
func AlarmError(message string) {
	logger.Error(message,
		zap.String("app", basicInfo.Name),
		zap.String("category", "alarm-error"),
	)
	if mailDialer != nil {
		m := gomail.NewMessage()
		receivers := alarmConfig.Receivers
		m.SetHeader("From", mailSender)
		m.SetHeader("To", receivers...)
		m.SetHeader("Subject", "Alarm-"+basicInfo.Name)
		m.SetBody("text/plain", message)
		// 避免发送邮件时太慢影响现有流程
		go func() {
			// 一次只允许一个email发送（由于使用的邮件服务有限制）
			sendingMailMutex.Lock()
			defer sendingMailMutex.Unlock()
			err := mailDialer.DialAndSend(m)
			if err != nil {
				logger.Error("send mail fail",
					zap.Error(err),
				)
			}
		}()
	}
}
