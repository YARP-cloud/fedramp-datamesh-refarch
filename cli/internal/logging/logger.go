package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Logger is a simple logger for the CLI
type Logger struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	debugLog *log.Logger
	verbose  bool
}

// NewLogger creates a new logger
func NewLogger() *Logger {
	// Create logs directory in user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Could not find home directory: %v", err)
		return &Logger{
			errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
			debugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
			verbose:  false,
		}
	}
	
	logsDir := filepath.Join(home, ".fedramp-data-mesh", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Printf("Could not create logs directory: %v", err)
		return &Logger{
			errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
			debugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
			verbose:  false,
		}
	}
	
	// Open log files
	errorFile, err := os.OpenFile(
		filepath.Join(logsDir, "error.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Printf("Could not open error log file: %v", err)
		errorFile = os.Stderr
	}
	
	infoFile, err := os.OpenFile(
		filepath.Join(logsDir, "info.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Printf("Could not open info log file: %v", err)
		infoFile = os.Stdout
	}
	
	debugFile, err := os.OpenFile(
		filepath.Join(logsDir, "debug.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Printf("Could not open debug log file: %v", err)
		debugFile = os.Stdout
	}
	
	return &Logger{
		errorLog: log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(infoFile, "INFO: ", log.Ldate|log.Ltime),
		debugLog: log.New(debugFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		verbose:  false,
	}
}

// SetVerbose turns on verbose logging
func (l *Logger) SetVerbose(verbose bool) {
	l.verbose = verbose
}

// Infof logs an info level message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.infoLog.Output(2, fmt.Sprintf(format, args...))
	if l.verbose {
		fmt.Printf("INFO: "+format+"\n", args...)
	}
}

// Errorf logs an error level message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.errorLog.Output(2, fmt.Sprintf(format, args...))
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}

// Debugf logs a debug level message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.debugLog.Output(2, fmt.Sprintf(format, args...))
	if l.verbose {
		fmt.Printf("DEBUG: "+format+"\n", args...)
	}
}
