# Configure a custom logger

1. Create a `struct` that will wrap around your logger.
    
    ```go
    type CustomAdapter struct {
	    ...
    }
    ```

2. Implement the `logger.Interface()` method.
    ```go
    func (l CustomAdapter) Trace(msg string, args ...any) {
    	...
    }

   ...

    func (l CustomAdapter) With(key string, fields Fields) Interface {
        ...
	    return CustomAdapter{...}
    }

    ```

3. Use the `logger.SetLogger()` method to set the `struct` as a default logger.
    * Make sure to set a logger before retrieving the session. You need to do this because the `logger.SetLogger()` method doesn’t overwrite any logger instances inside already existing sessions. It means that if you retrieve a session that has already existing logger instances but you haven't set a new logger, then this session will use an old or default logger.
    ```go
        logger.SetLogger(CustomAdapter)
    ```

4. Retrieve the logger with the `logger.Default()` method.
    ```go 
        log := logger.Default()
    ```

## Example

```go
package main

import (
	"fmt"
	"maps"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
)

const (
	TraceLevel = iota - 2 // assuming that InfoLevel is default
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type CustomLogger struct {
	loggingLevel int
	key          string
	fields       log.Fields
}

func (l CustomLogger) Fatal(msg string, _ ...any) {
	if l.loggingLevel <= FatalLevel {
		fmt.Println(fmt.Sprintf("[FATAL] %s %s %s", l.key, l.fields.Get(), msg))
	}
}

func (l CustomLogger) Fatalf(msg string, args ...any) {
	if l.loggingLevel <= FatalLevel {
		fmt.Println(fmt.Sprintf("[FATAL] %s %s %s", l.key, l.fields.Get(), fmt.Sprintf(msg, args)))
	}
}

func (l CustomLogger) Error(msg string, _ ...any) {
	if l.loggingLevel <= ErrorLevel {
		fmt.Println(fmt.Sprintf("[ERROR] %s %s %s", l.key, l.fields.Get(), msg))
	}
}

func (l CustomLogger) Errorf(msg string, args ...any) {
	if l.loggingLevel <= ErrorLevel {
		fmt.Println(fmt.Sprintf("[ERROR] %s %s %s", l.key, l.fields.Get(), fmt.Sprintf(msg, args)))
	}
}

func (l CustomLogger) Warn(msg string, _ ...any) {
	if l.loggingLevel <= WarnLevel {
		fmt.Println(fmt.Sprintf("[WARN] %s %s %s", l.key, l.fields.Get(), msg))
	}
}

func (l CustomLogger) Warnf(msg string, args ...any) {
	if l.loggingLevel <= WarnLevel {
		fmt.Println(fmt.Sprintf("[WARN] %s %s %s", l.key, l.fields.Get(), fmt.Sprintf(msg, args)))
	}
}

func (l CustomLogger) Info(msg string, _ ...any) {
	if l.loggingLevel <= InfoLevel {
		fmt.Println(fmt.Sprintf("[INFO] %s %s %s", l.key, l.fields.Get(), msg))
	}
}

func (l CustomLogger) Infof(msg string, args ...any) {
	if l.loggingLevel <= InfoLevel {
		fmt.Println(fmt.Sprintf("[INFO] %s %s %s", l.key, l.fields.Get(), fmt.Sprintf(msg, args)))
	}
}

func (l CustomLogger) Debug(msg string, _ ...any) {
	if l.loggingLevel <= DebugLevel {
		fmt.Println(fmt.Sprintf("[DEBUG] %s %s %s", l.key, l.fields.Get(), msg))
	}
}

func (l CustomLogger) Debugf(msg string, args ...any) {
	if l.loggingLevel <= DebugLevel {
		fmt.Println(fmt.Sprintf("[DEBUG] %s %s %s", l.key, l.fields.Get(), fmt.Sprintf(msg, args)))
	}
}

func (l CustomLogger) Trace(msg string, _ ...any) {
	if l.loggingLevel <= TraceLevel {
		fmt.Println(fmt.Sprintf("[TRACE] %s %s %s", l.key, l.fields.Get(), msg))
	}
}

func (l CustomLogger) Tracef(msg string, args ...any) {
	if l.loggingLevel <= TraceLevel {
		fmt.Println(fmt.Sprintf("[TRACE] %s %s %s", l.key, l.fields.Get(), fmt.Sprintf(msg, args)))
	}
}

func (l CustomLogger) With(key string, fields log.Fields) log.Interface {
	maps.Copy(l.fields, fields)
	return CustomLogger{loggingLevel: l.loggingLevel, key: l.key + key, fields: l.fields}
}

func main() {
	edgerc, _ := edgegrid.New(
		edgegrid.WithFile("~/.edgerc"),
		edgegrid.WithSection("default"),
	)

	l := CustomLogger{loggingLevel: DebugLevel, key: "MyCustomLogger", fields: map[string]interface{}{"k1": "v1"}}

	l1 := l.With("_Modified", map[string]interface{}{"k2": "v2"})

	sess, _ := session.New(
		session.WithSigner(edgerc),
		session.WithHTTPTracing(true),
	)

	var userProfile struct {
		FirstName         string `json:"firstName"`
		LastName          string `json:"lastName"`
		UserName          string `json:"uiUserName"`
		TimeZone          string `json:"timeZone"`
		Country           string `json:"country"`
		PreferredLanguage string `json:"preferredLanguage"`
		SessionTimeOut    *int   `json:"sessionTimeOut"`
	}

	req, _ := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)

	req = req.WithContext(session.ContextWithOptions(req.Context(), session.WithContextLog(l1)))

	result, err := sess.Exec(req, &userProfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result, userProfile)
}
```