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

}