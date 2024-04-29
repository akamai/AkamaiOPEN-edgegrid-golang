### How to use different logging backend

1. Create struct that will wrap around your logger
    
    ```go
    type CustomAdapter struct {
	    ...
    }
    ```

2. Implement `logger.Interface` methods
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

3. Use `logger.SetLogger` to set the struct as default logger
    ```go
        logger.SetLogger(CustomAdapter)
    ```
    * make sure to set logger before retrieving it, as it will not overwrite already existing instances

4. Retrieve the logger with `logger.Default()`
    ```go 
        log := logger.Default()
    ```
