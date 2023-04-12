package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main(){
	
	// --------------------------------
	// Single Connection to timescaledb
	// --------------------------------
	_ = godotenv.Load()
	ctx := context.Background()
	connStr := os.Getenv("DATABASE_CONENECTION_STRING")	// postgres://username:password@host:port/dbname?sslmode=require

	conn, err := pgx.Connect(ctx, connStr)
	// conn, err := pgxpool.New(ctx, connStr)
	if err!= nil {
		fmt.Println("Unable to connect to db, \n", err)
		return
	}
	fmt.Println("Connection successful !")
	defer conn.Close(ctx)
	// defer conn.Close()



	// ------------------------------------------
	// run a simple query to check our connection 
	// ------------------------------------------
	
	var greetings string
	err=conn.QueryRow(ctx,"select 'Hello, Timescale!'").Scan(&greetings)
	if (err!= nil){
		fmt.Printf("QueryRow failed %v\n",err)
		return
	}
	fmt.Println(greetings)
	fmt.Println("------------------------")


	
	// -----------------------
	// Create Relational Table
	// -----------------------

	queryCreateTable := `CREATE TABLE sensors (id SERIAL PRIMARY KEY,
	type VARCHAR(50), location VARCHAR(50));`
	_,err=conn.Exec(ctx, queryCreateTable)
	if err!= nil {
		fmt.Println("Unable to create sensors table: \n",err)
		return
	}
	fmt.Println("Successfully created relational table: Sensors")
	fmt.Println("------------------------")




	// -------------------
	// Generate hypertable
	// -------------------

	queryCreateTable = `CREATE TABLE sensor_data (
        time TIMESTAMPTZ NOT NULL,
        sensor_id INTEGER,
        temperature DOUBLE PRECISION,
        cpu DOUBLE PRECISION,
        FOREIGN KEY (sensor_id) REFERENCES sensors (id));
        `
	queryCreateHyperTable:= `SELECT create_hypertable('sensor_data','time');`
	_,err = conn.Exec(ctx,queryCreateTable+queryCreateHyperTable)
	if err!=nil {
		fmt.Println("UNable to create sensor_data hypertable")
		return
	}
	fmt.Println("Successfully created hypertable 'sensor_data'")
	fmt.Println("------------------------")

}