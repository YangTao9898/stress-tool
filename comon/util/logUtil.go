package util

import "github.com/phachon/go-logger"

var logger *go_logger.Logger

func init()  {
	logger = initLogger()
}

func initLogger() *go_logger.Logger{

	logger := go_logger.NewLogger()

	logger.Detach("console")

	// console adapter config
	consoleConfig := &go_logger.ConsoleConfig{
		Color: false,      // Does the text display the color
		JsonFormat: false, // Whether or not formatted into a JSON string
		Format: "",        // JsonFormat is false, logger message output to console format string
	}

	// add output to the console
	//      LOGGER_LEVEL_EMERGENCY  = 0
	//  	LOGGER_LEVEL_ALERT      = 1
	//  	LOGGER_LEVEL_CRITICAL   = 2
	//  	LOGGER_LEVEL_ERROR      = 3
	//  	LOGGER_LEVEL_WARNING    = 4
	//  	LOGGER_LEVEL_NOTICE     = 5
	//  	LOGGER_LEVEL_INFO       = 6
	//  	LOGGER_LEVEL_DEBUG      = 7
	logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)
	return  logger
}

func GetLogger() *go_logger.Logger{
	return logger
}
