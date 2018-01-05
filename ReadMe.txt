Install the go and set the path vataibles, Please check @ https://golang.org/doc/install


Router
Since we're trying to build out a webserver, let's find a router package. When searching for a go package to use, GoDoc is always a good place to start. There are many excellect candidates, but let's use the httprouter package.

$ go get github.com/julienschmidt/httprouter

install 
go get github.com/ashokpottabathini/RestIntersectRect


we can run the code locally by calling
```
go run RestIntersectRect.go
```
Now we can call our server. For this, let's use curl. Curl is an HTTP client for the command line. By default, it sends GET requests, but the X parameter lets us create a PUT request instead.
First, let's fill the map with some entries. We do that by sending PUT requests with a key and a value.
Then we can request a list of all entries, as well as rectangle entries by name key r1 and cordinate values as x1,y1,x2,y2
```
curl -X PUT localhost:4040/entry/r1/0,10,10,0 
curl -X PUT localhost:4040/entry/r2/1,11,11,1

curl localhost:4040/list
curl localhost:4040/entry/r1
curl localhost:4040/entry/r2

```