package helper

import log "github.com/sirupsen/logrus"

func LoggerInit() {
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyMsg:  "message",
			log.FieldKeyTime: "timestamp",
		},
	})

}
