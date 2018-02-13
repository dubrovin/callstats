# callstats


### how to run

#### 1) build and run docker container
```
docker build -t callstats:latest .
docker run -d -p 8080:8080 callstats:latest 
```

#### 2) go run main.go or go build and run binary

```
    go run main.go
```

##### application provides next parameters:
```
addr       = flag.String("addr", ":8080", "http service address")
delay      = flag.String("delay", "1µs", "interval delay")
windowSize = flag.Int("window_size", 100, "size for sliding window")
filePath   = flag.String("file_path", "tests/test2.csv", "path to test file")
``` 


##### running application provides 2 methods: addDelay and getMedian:

##### addDelay:
```
curl -X POST \
  http://127.0.0.1:8080/add_delay \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 37acf7ae-a26a-432f-8612-35209289d6a4' \
  -d '{
	"delay": "1s"
}'
```

##### getMedian:
```
curl -X GET \
  http://127.0.0.1:8080/get_median \
  -H 'Cache-Control: no-cache' \
  -H 'Postman-Token: 4fad6579-9c88-d08e-a8c4-89179197abee'
```

##### while application works it persist calculated medians to out folder
