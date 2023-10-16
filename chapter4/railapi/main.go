package main

import (
	"database/sql"
	"encoding/json"
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

type Train struct {
	ID              int
	DriverName      string
	OperatingStatus bool
}

type Station struct {
	ID            int
	Name          string
	OperatingTime time.Time
	ClosingTime   time.Time
}

type Schedule struct {
	ID          int
	TrainID     int
	StationID   int
	ArrivalTime time.Time
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

	t.Register(wsContainer)

	log.Println("Server start to listening on localhost:8000")
	srv := &http.Server{Addr: ":8000", Handler: wsContainer}
	log.Fatal(srv.ListenAndServe())
}
