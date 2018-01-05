
/* #The webserver listeing on port 4040. The server uses map which contains key and values as strings.
#It uses update function for updating the rectanlge keys and values.
#show function used to show the map key values.
#IsIntersectRectangles function return true if rectangles r1 & r2 are intersecting otherwise false.
#IntersectingRectangle function returns intersecting rectanlge.

#User can send from curl client in the following format where r1 rectanlge one and it's coordinates are 0,10,10,0 then it displays on client terminal #"intersecting rectnalge" if they are intersecting otherwise "No Intersection".

#curl -X PUT localhost:4040/entry/r1/0,10,10,0 
#curl -X PUT localhost:4040/entry/r2/1,11,11,1
# Calculating Intersecting Rectangle
# Intersecting Rectangle: 1,10,10,1 */



// ## Imports and globals
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"strings"
	"strconv"

	// This is `httprouter`. Ensure to install it first via `go get`.
	"github.com/julienschmidt/httprouter"
)

// We need a data store. For our purposes, a simple map
// from string to string is completely sufficient.
type store struct {
	data map[string]string

	// Handlers run concurrently, and maps are not thread-safe.
	// This mutex is used to ensure that only one goroutine can update `data`.
	m sync.RWMutex
}

type Rectangle struct {
  	left, top, right, bottom int //left(x1) top(y1) right(x2) bottom(y2)
  }

var (
	// We need a flag for setting the listening address.
	// We set the default to port 4040, which is a HTTP port
	// for servers with local-only access.
	addr = flag.String("addr", ":4040", "http service address")

	// Now we create the data store.
	s = store{
		data: map[string]string{},
		m:    sync.RWMutex{},
	}
	// Now we create the rectangle
	r1 = Rectangle{0,0,0,0}
	r2 = Rectangle{0,0,0,0}
	r3 = Rectangle{0,0,0,0}
)

// ## main
func main() {
	// The main function starts by parsing the commandline.
	flag.Parse()

	// Now we can create a new `httprouter` instance...
	r := httprouter.New()

	// ...and add some routes.
	// `httprouter` provides functions named after HTTP verbs.
	// So to create a route for HTTP GET, we simply need to call the `GET` function
	// and pass a route and a handler function.
	// The first route is `/entry` followed by a key variable denoted by a leading colon.
	// The handler function is set to `show`.
	r.GET("/entry/:key", show)

	// We do the same for `/list`. Note that we use the same handler function here;
	// we'll switch functionality within the `show` function based on the existence
	// of a key variable.
	r.GET("/list", show)

	// For updating, we need a PUT operation. We want to pass a key and a value to the URL,
	// so we add two variables to the path. The handler function for this PUT operation
	// is `update`.
	r.PUT("/entry/:key/:value", update)

	// Finally, we just have to start the http Server. We pass the listening address
	// as well as our router instance.
	err := http.ListenAndServe(*addr, r)

	// For this demo, let's keep error handling simple.
	// `log.Fatal` prints out an error message and exits the process.
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// ## The handler functions

// Let's implement the show function now. Typically, handler functions receive two parameters:
//
// * A Response Writer, and
// * a Request object.
//
// `httprouter` handlers receive a third parameter of type `Params`.
// This way, the handler function can access the key and value variables
// that have been extracted from the incoming URL.
func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// To access these parameters, we call the `ByName` method, passing the variable name that we chose when defining the route in `main`.
	k := p.ByName("key")

	// The show function serves two purposes.
	// If there is no key in the URL, it lists all entries of the data map.
	if k == "" {
		// Lock the store for reading.
		s.m.RLock()
		fmt.Fprintf(w, "Read list: %v", s.data)
		s.m.RUnlock()
		return
	}

	// If a key is given, the show function returns the corresponding value.
	// It does so by simply printing to the ResponseWriter parameter, which
	// is sufficient for our purposes.
	s.m.RLock()
	fmt.Fprintf(w, "Read entry: s.data[%s] = %s", k, s.data[k])
	s.m.RUnlock()
}

// The update function has the same signature as the show function.
func update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Fetch key and value from the URL parameters.
	k := p.ByName("key")
	v := p.ByName("value")

	// We just need to either add or update the entry in the data map.
	s.m.Lock()
	s.data[k] = v
	s.m.Unlock()
	
	// Fill Rectangle
	fillRectangle(k)
	
	// we print the result to the ResponseWriter.
	fmt.Fprintf(w, "Updated: s.data[%s] = %s", k, v)
	
	if(strings.Compare(k, "r2") == 0){
		
		fmt.Fprintf(w, "\n Calculating Intersecting Rectangle")
		
		if r1.IsIntersectRectangles(r2){
			r3 = r1.IntersectingRectangle(r2);
			result := []string{}
			result = append(result, strconv.Itoa(r3.left), strconv.Itoa(r3.top), strconv.Itoa(r3.right),  strconv.Itoa(r3.bottom))
			fmt.Fprintf(w, "\n Intersecting Rectangle: %s", strings.Join(result, ","))
		}else{
			fmt.Fprintf(w, "\n No Intersection")
		}
	}
}

// This function fills the rectangle objects with map values based on key
func fillRectangle(key string){

	// Split on comma.
    result := strings.Split(s.data[key], ",")
	
	r:= Rectangle{0,0,0,0};
	
	//fill the rectangle
	r.left,_ = strconv.Atoi(result[0])
	r.top,_ = strconv.Atoi(result[1])
	r.right,_ = strconv.Atoi(result[2])
	r.bottom,_ = strconv.Atoi(result[3])
	
	if (strings.Compare(key, "r1") == 0){
		r1 = r;
	} else if(strings.Compare(key, "r2") == 0){
		r2 = r;
	} else {
		fmt.Println("\n Invalid key for rectangle, allowed keys are r1/r2")
	}
	
}

func Min(x, y int) int {
    if x < y {
        return x
    }
    return y
}

func Max(x, y int) int {
    if x > y {
        return x
    }
    return y
}

//This method returns true if rectangles r1 & r2 are intersecting otherwise false
func (r1 Rectangle) IsIntersectRectangles (r2 Rectangle) bool{
	
	if r2.left < r2.right && r2.right > r1.left && r2.bottom < r1.top && r2.top > r1.bottom	{
			return true;
		}

	return false;
}

// This method takes two rectangles as parameters
// The rectangles have left(x1), top(y1), right(x2) and bottom(y2) as accessible properties.
// It return a rectangle representing the intersecting between two rectangular regions
// If there is no intersecting, print not intersecting. Assume positive left(x1), top(y1), right(x2) and bottom(y2)
func (r1 Rectangle) IntersectingRectangle (r2 Rectangle) Rectangle{
    r3:= Rectangle{}
	r3.left = Max(r1.left, r2.left);
	r3.top = Min(r1.top, r2.top);
	r3.right = Min(r1.right, r2.right);
	r3.bottom = Max(r1.bottom, r2.bottom);
	
	return r3;
}
