### Step 1: Install Docker and TimescaleDB

To get started, you'll need to install Docker and TimescaleDB on your local machine. Here are the steps to do that:

1. Install Docker on your machine. You can download Docker from the official website: https://www.docker.com/get-started

2. Once Docker is installed, run the following command to download the TimescaleDB Docker image:

```
docker pull timescale/timescaledb:latest-pg15
```

### Step 2: Start the TimescaleDB Container

After downloading the TimescaleDB Docker image, you can start the container with the following command:

```
docker run --name timescaledb -e POSTGRES_PASSWORD=password -d -p 5432:5432 --restart always timescale/timescaledb:latest-pg15
```

This command starts a TimescaleDB container named timescaledb with the password password, and maps port 5432 to the host machine's port 5432.

### Step 3: Connect to the TimescaleDB Container

After starting the container, you can connect to it with the following command:

```
docker exec -it timescaledb psql -U postgres
```

This command opens a PostgreSQL shell inside the container, where you can create tables and run SQL commands.

### Step 4: Create a PostgreSQL Database & Table

1. To create a PostgreSQL database, you can run the following SQL command:

```
CREATE DATABASE timescale_test;
```

2. To create a PostgreSQL table, First you have to move to the following database
   ```
   \c timescale_test
   ```
3. Now you can run the following SQL command:

```
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);
```

This command creates a users table with columns for id, name, email, and password. The id column is an auto-incrementing primary key, and the email column has a unique constraint to ensure that each email address is only used once.

### Step 5: Create a TimescaleDB Hypertable

To create a TimescaleDB hypertable, you can run the following SQL command:

```
CREATE TABLE sensor_data (
    time TIMESTAMPTZ NOT NULL,
    sensor_id TEXT NOT NULL,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    PRIMARY KEY (time, sensor_id)
);
```

```
SELECT create_hypertable('sensor_data', 'time');
```

This command creates a sensor_data table with columns for time, sensor_id, temperature, and humidity. The time and sensor_id columns are used as the primary key, and the create_hypertable function is used to create a TimescaleDB hypertable for the sensor_data table.

## Connecting to the project

Your database is created. Now use this database is now ready to use. Inside of your project, declare a connection string variable where write the string in the below format:

### For Timescale cloud:

Copy the 'service URL' from the web dashboard.

### For Localhost (Using from docker)

write the string in the given format:

```
postgres://username:password@host:port/dbname?sslmode=disable
```

<br>
Now , You can connect and run this project on your machine.
