## Logger package

______

### How to use logger package

#### 1. Download the package :

```go
    go get github.com/cploutarchou/go-elastic-logger/v7
```

#### 2. Import the package :

```go
    import github.com/cploutarchou/go-elastic-logger/v7
```
#### 3. Set the Elasticsearch client:
```go
   cfg := elasticsearch.Config{
    Addresses: []string{"http://localhost:9200"},
    }
    client, err := elasticsearch.NewClient(cfg)
    if err != nil {
        panic(err)
    }
    logger.SetElasticsearchClient(client)

```
          
#### 4. Set the index name:
                            
```go
    logger.SetIndex("myindex")
```
         
#### 5. Set the log level
```go
    logger.SetLogLevel(logger.DEBUG)
```

#### 6. Use the log level function
```go
    logger.Debug("Debug message")
    logger.Debugf("Debug message with format: %s", "hello")
    logger.Info("Info message")
    logger.Infof("Info message with format: %d", 123)
    logger.Warning("Warning message")
    logger.Warningf("Warning message with format: %f", 3.14)
    logger.Error("Error message")
    logger.Errorf("Error message with format: %t", true)
    logger.Fatal("Fatal message")
    logger.Fatalf("Fatal message with format: %v", "goodbye")
```


### Notes:
* The package will output logs to stdout and index them in Elasticsearch.
* The package is concurrent safe and the elasticsearch client will be initiated only once.
* If you want to change the timestamp format you can change the format by using the `SetTimeFormat()`function. Default layout is 2023-01-02 15:04:05
* You can use the SetLogLevel function to filter the logs you want to see based on the log level.
                      

### Example: 
```go  
package main

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// setElasticClient sets the elasticsearch client to be used for logging.
func setElasticClient() {
	cnf := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "your_username",
		Password:  "your_password",
	}
	client, err := elasticsearch.NewClient(cnf)
	if err != nil {
		panic(err)
	}
	//sets the client as the Elasticsearch client to be used.
	SetElasticsearchClient(client)

}

func testFnc() {
	// sets the index name to be used for logging.
	SetIndex("logs")
	//sets the log level to debug.
	SetLogLevel(DEBUG)
	//sets the format of the timestamp in the logs.
	SetTimeFormat("2006-01-02 15:04:05")
	//logs a debug message.
	Debug("This is a debug message")
	return
}
func main() {
	//sets the elasticsearch client to be used for logging.
	setElasticClient()
	testFnc()
}
```
       
### Elastic Doc Example:
```json
{
  "_index": "logs",
  "_id": "LmUPsIUBS2V1Rn1ibiUZ",
  "_score": 1.0,
  "_source": {
    "file": "/Users/christos/workspace/go-elastic-logger/example.go:33",
    "function": "testFnc()",
    "level": "DEBUG",
    "message": "This is a debug message",
    "timestamp": "2023-01-14 13:34:38"
  }
}
```
