package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Zekeriyyah/hands-on-RestApis/chapter4/railapi/dbutils"
	"github.com/Zekeriyyah/hands-on-RestApis/chapter4/railapi/utils"
	restful "github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const layout = "2006-01-02 15:04:05-07:00"

type Train struct {
	ID              int    `json:"id"`
	DriverName      string `json:"driver_name"`
	OperatingStatus bool   `json:"operating_status"`
}

type Station struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	OpeningTime time.Time `json:"opening_time"`
	ClosingTime time.Time `json:"closing_time"`
}

type Schedule struct {
	ID          int       `json:"id"`
	TrainID     int       `json:"train_id"`
	StationID   int       `json:"station_id"`
	ArrivalTime time.Time `json:"arrival_time"`
}

func (t *Train) Register(container *restful.Container) {
	//Create new webservice
	ws := new(restful.WebService)

	//Define path with expected request and response format
	ws.Path("/v1/trains").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	t = &Train{}
	//Define the various routes for the webservice
	ws.Route(ws.GET("/{train_id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.CreateTrain))
	ws.Route(ws.DELETE("/{train_id}").To(t.removeTrain))

	//Add the webservice to the container
	container.Add(ws)

	/**
		Note that path are the URL endpoints while the routes are the path
		parameters or the query parameters
	**/
}

// Defining the handlers for the train resource.

func (t Train) getTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train_id")

	//Validating the Id in url to avoid sql injection
	utils.ParseID(response, id)

	//Query database for train entity with ID = id
	err := DB.QueryRow("select ID, DRIVER_NAME, OPERATING_STATUS FROM train where id=?", id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)

	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Train could not be found.")
	} else {
		response.WriteEntity(t)
	}

}

func (t Train) CreateTrain(req *restful.Request, resp *restful.Response) {
	log.Println(req.Request.Body)

	//Decoding the request body
	decoder := json.NewDecoder(req.Request.Body)

	newTrain := &Train{}

	err := decoder.Decode(newTrain)

	if err != nil {
		log.Println("Decoding Failed! :", err)
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, "Failed to Decode Request Body!")
		return
	}

	statement, _ := DB.Prepare("insert into train (DRIVER_NAME, OPERATING_STATUS) values (?,?)")
	result, err := statement.Exec(newTrain.DriverName, newTrain.OperatingStatus)
	if err == nil {
		newId, _ := result.LastInsertId()
		newTrain.ID = int(newId)

		resp.WriteHeaderAndEntity(http.StatusOK, newTrain)
	} else {
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (t Train) removeTrain(req *restful.Request, resp *restful.Response) {
	Id := req.PathParameter("train_id")

	//Validating the Id in url to avoid sql injection
	_, parseErr := strconv.ParseInt(Id, 0, 0)
	if parseErr != nil {
		log.Println("Invalid ID!: ", parseErr)
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusBadRequest, "Invalid ID as the request parameter!")
		return
	}

	statement, err := DB.Prepare("delete from train where id=?")
	if err != nil {
		log.Println("Invalid delete sql command", err)
		return
	}
	_, err = statement.Exec(Id)
	if err == nil {
		log.Println("Train Resource successfully deleted.")

		resp.WriteHeader(http.StatusOK)
		resp.AddHeader("Content-Type", "text/plain")
		resp.Write([]byte("Successfully Deleted!"))
	} else {

		log.Println("No deletion occured: ", err)

		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

/**
	Implementing webservices on the Station struct resources following the same way
	the that of train resouces have been implemented
**/

// Implement a Station resource method "Register" to register it to add its webservices to a restful container
func (s *Station) Register(container *restful.Container) {
	//Initialize a new webservice
	ws := new(restful.WebService)

	//Define various Routes and merge them with their respective handler
	ws.Path("/v1/stations").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{station-id}").To(s.getStation))
	ws.Route(ws.POST("").To(s.createStation))
	ws.Route(ws.PUT("/{station-id}").To(s.updateStation))
	ws.Route(ws.DELETE("/{station-id}").To(s.deleteStation))

	//Add the webservice to container
	container.Add(ws)
}

func (s Station) createStation(req *restful.Request, resp *restful.Response) {
	//Create a new instance of station resource
	newStation := &Station{}

	//Decode the request body into it. If returned err is not nil, display error msg otherwise write a query to store the new station
	//Give a success response
	decoder := json.NewDecoder(req.Request.Body)

	err := decoder.Decode(newStation)
	if err != nil {
		log.Println(err)
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, "Could not decode request body...")
		return
	}

	statement, creatErr := DB.Prepare("insert into Station (NAME, OPENING_TIME, CLOSING_TIME) values (?,?,?)")
	if creatErr != nil {
		panic(creatErr)
	}

	result, err := statement.Exec(newStation.Name, newStation.OpeningTime, newStation.ClosingTime)
	if err == nil {
		newId, _ := result.LastInsertId()
		newStation.ID = int(newId)
		resp.WriteAsJson(newStation)
	} else {
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, err.Error())
	}
}

func (s Station) getStation(req *restful.Request, resp *restful.Response) {
	s = Station{}
	id := req.PathParameter("station-id")

	//validate the id
	if utils.ParseID(resp, id) != nil {
		return
	}

	//Query the database
	var openingTimeStr string
	var closingTimeStr string

	err := DB.QueryRow("select ID, NAME, OPENING_TIME, CLOSING_TIME FROM station where id=?", id).Scan(&s.ID, &s.Name, &openingTimeStr, &closingTimeStr)
	if err != nil {
		log.Print("Scanning station resource failed: ", err)
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusNotFound, "Station With ID NOT FOUND!!")
		return
	}

	//preparing the response
	s.OpeningTime, err = time.Parse(layout, openingTimeStr)
	if err != nil {
		log.Println("Could not parse openingTimeStr: ", err)
	}
	s.ClosingTime, _ = time.Parse(layout, closingTimeStr)
	log.Printf("o: %v c:%v\n", s.OpeningTime, s.ClosingTime)

	//Formatting timestr to store only the time
	s.OpeningTime, _ = time.Parse(time.TimeOnly, s.OpeningTime.Format(time.TimeOnly))
	s.ClosingTime, _ = time.Parse(time.TimeOnly, s.ClosingTime.Format(time.TimeOnly))

	resp.WriteHeaderAndEntity(http.StatusOK, s)
}

func (s Station) updateStation(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("station-id")
	utils.ParseID(resp, id)

	//Decoding request body
	s = Station{}
	decoder := json.NewDecoder(req.Request.Body)
	decoder.DisallowUnknownFields()
	_ = decoder.Decode(&s)

	Id, _ := strconv.ParseInt(id, 0, 0)
	s.ID = int(Id)

	//Querying Database
	statement, _ := DB.Prepare("update station set NAME=?, OPENING_TIME=?, CLOSING_TIME=? where id=?")
	_, err := statement.Exec(s.Name, s.OpeningTime, s.ClosingTime, id)
	if err != nil {
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, "Could Not Execute Update Command")
		return
	}

	//Writing response
	resp.WriteHeaderAndEntity(http.StatusOK, s)
}

func (s Station) deleteStation(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("station-id")
	//Validating id
	err := utils.ParseID(res, id)
	if err != nil {
		return
	}
	var ID string
	//Checking if station exist in the database
	err = DB.QueryRow("select ID from station where id=?", id).Scan(&ID)
	if err != nil {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusNotFound, "SORRY, STATION DOES'NT EXIST!!")
		return
	}
	_, err = DB.Exec("delete from station where id=?", id)
	if err != nil {
		log.Fatal("FAILED TO EXECUTE DELETE COMMAND!")
		return
	} else {
		res.AddHeader("Content-Type", "text/plain")
		res.Write([]byte(fmt.Sprintf("Content of station with ID %v deleted successfully!", id)))
	}
}

func main() {
	//Connect to database
	var err error
	DB, err = sql.Open("sqlite3", "./railapi.db")
	if err != nil {
		log.Println("Drivers creation failed!!")
	}

	//Create tables
	dbutils.Initialize(DB)

	//Creating a go-restful container and registering train resource
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})

	t := &Train{}
	s := &Station{}

	t.Register(wsContainer)
	s.Register(wsContainer)

	log.Println("Server start to listening on localhost:8001")
	srv := &http.Server{Addr: ":8001", Handler: wsContainer}
	log.Fatal(srv.ListenAndServe())
}
