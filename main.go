package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {

	// --------------------------------
	// Single Connection to timescaledb
	// --------------------------------
	_ = godotenv.Load()
	ctx := context.Background()
	connStr := os.Getenv("DATABASE_CONENECTION_STRING") // postgres://username:password@host:port/dbname?sslmode=require

	conn, err := pgx.Connect(ctx, connStr)
	// conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
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
	err = conn.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greetings)
	if err != nil {
		fmt.Printf("QueryRow failed %v\n", err)
		return
	}
	fmt.Println(greetings)
	fmt.Println("------------------------")

	// -----------------------
	// Create Relational Table
	// -----------------------

	queryCreateTable := `CREATE TABLE sensors (id SERIAL PRIMARY KEY,
	type VARCHAR(50), location VARCHAR(50));`
	_, err = conn.Exec(ctx, queryCreateTable)
	if err != nil {
		fmt.Println("Unable to create sensors table: \n", err)
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
	queryCreateHyperTable := `SELECT create_hypertable('sensor_data','time');`
	_, err = conn.Exec(ctx, queryCreateTable+queryCreateHyperTable)
	if err != nil {
		fmt.Println("UNable to create sensor_data hypertable")
		return
	}
	fmt.Println("Successfully created hypertable 'sensor_data'")
	fmt.Println("------------------------")

	// ------------------------------
	// insert rows of data into table
	// ------------------------------

	// single row insert
	sensorTypes := []string{"a", "b", "c", "d"}
	sensorLocations := []string{"floor", "ceiling", "floor", "ceiling"}

	for i := range sensorTypes {
		queryInsertMetadata := `INSERT INTO sensors (type,location) VALUES ($1,$2);`

		_, err = conn.Exec(ctx, queryInsertMetadata, sensorTypes[i], sensorLocations[i])
		if err != nil {
			fmt.Println("Unable to insert data into the table: sensors \n", err)
			return
		}
		fmt.Printf("Inserted sensor (%s,%s) into database\n", sensorTypes[i], sensorLocations[i])
	}
	fmt.Println("------------------------")
	fmt.Println("Successfully Inserted all sensors into database")
	fmt.Println("------------------------")

	// multiple row insert

	//generate random data
	queryDataGeneration := `
	SELECT generate_series(now() - interval '24 hour', now(), interval '5 minute') AS time,
	floor(random()*(3)+1)::int as sensor_id,
	random()*100 AS temperature,
	random() AS cpu
	`

	rows, err := conn.Query(ctx, queryDataGeneration)
	if err != nil {
		fmt.Println("Unable to generate sensor data\n", err)
		return
	}
	defer rows.Close()
	fmt.Println(rows)

	fmt.Println("Successfully generated sensor_data")

	type result struct {
		Time        time.Time
		SensorId    int
		Temperature float64
		CPU         float64
	}

	var results []result
	for rows.Next() {
		var r result
		err = rows.Scan(&r.Time, &r.SensorId, &r.Temperature, &r.CPU)
		if err != nil {
			fmt.Println("Unable to scan: \n", err)
			return
		}
		results = append(results, r)
	}

	// Any errors encountered by rows.Next or rows.Scan methods are returned here
	if rows.Err() != nil {
		fmt.Println("rows error \n", rows.Err())
		return
	}

	fmt.Println("Contents of result slice")
	for i := range results {
		var r result
		r = results[i]
		fmt.Printf("Time: %s  | ID: %d  | Temperature: %f  | CPU: %f \n", &r.Time, r.SensorId, r.Temperature, r.CPU)
	}
	fmt.Println("------------------------")

	queryInsertTimeseriesData := `
	INSERT INTO sensor_data (time, sensor_id, temperature,cpu) VALUES ($1,$2,$3,$4);`

	for i := range results {
		var r result
		r = results[i]
		_, err := conn.Exec(ctx, queryInsertTimeseriesData, r.Time, r.SensorId, r.Temperature, r.CPU)
		if err != nil {
			fmt.Println("Unable to insert sample data into timescale", err)
			return
		}
		defer rows.Close()
	}
	fmt.Println("Sucessfully inserted samples into sensor_data hypertable")
	fmt.Println("---------------------------")

	// ---------------
	// Execute a query
	// ---------------

	queryTimebucketFiveMin := `
	SELECT time_bucket('5 minutes',time) AS five_min, avg(cpu) FROM sensor_data
	JOIN sensors ON sensors.id = sensor_data.sensor_id WHERE sensors.location = $1
	AND sensors.type = $2 GROUP BY five_min
	ORDER BY five_min DESC;
	`
	rows, err = conn.Query(ctx, queryTimebucketFiveMin, "floor", "a")
	if err != nil {
		fmt.Println("Unable to execute query \n", err)
		return
	}
	defer rows.Close()
	fmt.Println("Successfully executed query")

	type result2 struct {
		Bucket time.Time
		Avg    float64
	}

	// Print rows returned and fill up results slice for later use
	var results2 []result2
	for rows.Next() {
		var r result2
		err = rows.Scan(&r.Bucket, &r.Avg)
		if err != nil {
			fmt.Println("Unable to scan \n", err)
			return
		}
		results2 = append(results2, r)
		fmt.Printf("Time bucket: %s | Avg: %f\n", &r.Bucket, r.Avg)
	}

	// Any errors encountered by rows.Next or rows.Scan are returned here
	if rows.Err() != nil {
		fmt.Printf("rows Error: %v\n", rows.Err())
		return
	}

}
