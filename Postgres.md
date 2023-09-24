#### step 1 : pull docker image from docker hub , you can use alpine too but it's very minimal just under 100 mb thus it be missing important functionalities.

```
docker pull postgres
```

#### step 2 : - To obtain the list of existing Docker Images, run the following command.
```
docker images
```

#### step 3 : setting up PostgreSQL on Docker
```
docker run --name postgres -e POSTGRES_USER=<username> -e POSTGRES_PASSWORD=<password> -p 5432:5432 -v <path on host system 
to persist data or backup>:/var/lib/postgresql/data -d postgres
```

In the command given above, 

- **PostgreSQL** is the name of the Docker Container.
- **-e POSTGRES_USER** is the parameter that sets a unique username to the Postgres database.
- **-e POSTGRES_PASSWORD** is the parameter that allows you to set the password of the Postgres database.
- **-p 5432:5432** is the parameter that establishes a connection between the Host Port and Docker Container Port. In this case, both ports are given as 5432, which indicates requests sent to the Host Ports will automatically redirect to the Docker Container Port. In addition, 5432 is also the same port where PostgreSQL will be accepting requests from the client.
- **-v** is the parameter that synchronizes the Postgres data with the local folder. This ensures that Postgres data will be safely present within the Home Directory even if the Docker Container is terminated.(To back up the data, we mounted the _/var/lib/postgresql/data_ directory to a folder on host(using pwd) the  directory of the host machine of the _postgres_ container.
- **-d** is the parameter that runs the Docker Container in the detached mode, i.e., in the background. If you accidentally close or terminate the Command Prompt, the Docker Container will still run in the background.
- **Postgres** is the name of the Docker image that was previously downloaded to run the Docker Container.

  ### OR

use docker compose : 

```

version: '3'
services:
  postgres:
    image: postgres:latest
    container_name: postgres
    environment:
      POSTGRES_DB: xyz
      POSTGRES_USER: xyz
      POSTGRES_PASSWORD: xyz 
    ports:
      - "5432:5432"
    volumes:
      - <path to schema>/schema.sql:/docker-entrypoint-initdb.d/schema.sql
      - <path to a fodler to backup data>:/var/lib/postgresql/data

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    depends_on:
      - postgres
    environment:
      GF_DATABASE_TYPE: postgres
      GF_DATABASE_HOST: postgres
      GF_DATABASE_NAME: xyz
      GF_DATABASE_USER: xyz
      GF_DATABASE_PASSWORD: xyz
    volumes:
      - <path to bakcup grafana data>:/var/lib/grafana


```
then use : save this file in porject folder and in same folder run 

```
docker compose start 
```
** NOTE : Grafana container may stop on first run restart it via cli or desktop app 


#### step 4: check the status of the newly created containers

```
docker ps -a
```

output

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS                    NAMES
<...>  postgres   "docker-entrypoint.s…"   20 seconds ago   Up 19 seconds   0.0.0.0:5432->5432/tcp   postgres
```

#### step 5 : Running an Interactive Shell in a Docker Container

```
docker exec -it postgresrpi bash
```

#### step 6 : connect to particular database via a user :

```
psql -d <db name> -U <username> -W
```
The above command includes three flags:

- `-d` - specifies the name of the database to connect to
- `-U` - specifies the name of the user to connect as
- `-W` - forces psql to ask for the user password before connecting to the database

some queries for grafana :

#### schema.sql file 

```
-- Create schema
create schema telemetry;
-- Create the devices table
CREATE TABLE telemetry.devices (
    MacAddress VARCHAR(25) PRIMARY KEY,
    PrivateIP VARCHAR(25),
    PublicIP VARCHAR(25),
    Hostname VARCHAR(25),
    OSType VARCHAR(25),
    TotalMemory INT
);

-- Create the rpi4b_metrics table
CREATE TABLE telemetry.rpi4b_metrics (
    MacAddress VARCHAR(25),
    CPUuserLoad DOUBLE PRECISION,
    MemoryUsage INT,
    Temperature REAL,
    TotalProcesses INT,
    TimeStamp TIMESTAMP UNIQUE,
    CONSTRAINT fk_MacAddress FOREIGN KEY (MacAddress) REFERENCES telemetry.devices(MacAddress)
);

-- Create the rpi_temp_alert table
CREATE TABLE telemetry.rpi_temp_alert (
    MacAddress VARCHAR(25),
    CPUuserLoad DOUBLE PRECISION,
    MemoryUsage INT,
    PrivateIP VARCHAR(25),
    Temperature REAL,
    TotalProcesses INT,
    TimeStamp TIMESTAMP,
    CONSTRAINT fk_MacAddress FOREIGN KEY (MacAddress) REFERENCES telemetry.devices(MacAddress)
);


```
query examples for grafana
```
select  telemetry.rpi4b_metrics.temperature , telemetry.rpi4b_metrics.timestamp as time  ,  telemetry.devices.privateip
FROM telemetry.rpi4b_metrics
FULL JOIN telemetry.devices
ON telemetry.rpi4b_metrics.MacAddress = telemetry.devices.MacAddress ;

```
#### NOTE : here 'macAd' is a variable , present in settings menu (look for gear icon) in created dashboard this can used for per device data
```
SELECT
  telemetry.rpi4b_metrics.MemoryUsage , 
  telemetry.rpi4b_metrics.timestamp as time,
  telemetry.devices.MacAddress
FROM telemetry.rpi4b_metrics
FULL JOIN telemetry.devices
ON telemetry.rpi4b_metrics.MacAddress = telemetry.devices.MacAddress
WHERE
  telemetry.devices.MacAddress = '$macAd'

```


 


