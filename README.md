
**Introduction**
The provided Go package `main`  is a comprehensive logging library that facilitates logging messages to both files and Kafka topics. It offers various customization options, enabling fine-grained control over the logging behavior.

**Initialization**
1. Import the package:
   
```go

import "github.com/vizvasrj/l"

```

2. Define the desired storage types:
   - **File Storage**:
     
```go

fileStorageType := l.FileStorageType{
         Filename:   "test.log", // File name to write logs to
         MaxSize:    1,         // Maximum size of the log file in megabytes. When the file reaches this size, it will be rotated.
         MaxBackups: 3,         // Maximum number of backup log files to keep. Older logs will be deleted as new logs are created.
         MaxAge:     7,         // Maximum age of the log files in days. Logs older than this age will be deleted.
     }

```
   - **Kafka Storage**:
     
```go

kafkaStorageType := l.KafkaStorageType{
         Brokers: []string{"kafka-257bfc54-lakjos-f2b6.a.aivencloud.com:19190"}, // List of Kafka brokers to connect to
         Topic:   "logs",                                      // Kafka topic to write logs to
     }

```

3. Create a logger instance:
   
```go

log := l.Logger{
       FileStorageType:  fileStorageType,
       KafkaStorageType: kafkaStorageType,
   }

```

**Logging Levels**
This logger supports the following logging levels:
- `Debug`
- `Info`
- `Warning`
- `Error`
- `Fatal`

**Logging Messages**
To log a message, use the appropriate logging function, passing in the log level and the message to log. For example, to log an informational message:
```go

log.InfoF("This is an informational message")

```
You can also use formatted strings with `%v`, `%s`, `%d`, etc., to include dynamic values in your log messages.

**Conclusion**
This package offers a flexible logging solution. By combining file and Kafka storage options, it provides comprehensive logging capabilities, making it adaptable to diverse logging requirements.