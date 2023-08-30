step 1 : pull docker image from docker hub , you can use alpine too but it's very minimal just under 100 mb thus it be missing important functionalities.

```
docker pull postgres
```

step 2 : - To obtain the list of existing Docker Images, run the following command.
```
docker images
```

step 3 : setting up PostgreSQL on Docker
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

step 4: check the status of the newly created PostgreSQL container.
```
docker ps -a
```

output

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS                    NAMES
<...>  postgres   "docker-entrypoint.s…"   20 seconds ago   Up 19 seconds   0.0.0.0:5432->5432/tcp   postgres
```

step 5 : Running an Interactive Shell in a Docker Container

```
docker exec -it postgresrpi bash
```

step 6 : connect to particular database via a user :

```
psql -d iot -U <username> -W
```
The above command includes three flags:

- `-d` - specifies the name of the database to connect to
- `-U` - specifies the name of the user to connect as
- `-W` - forces psql to ask for the user password before connecting to the database

step 7 : 
```
CREATE SCHEMA IF NOT EXISTS telemetry;
```

step 8 :
```
create table telemetry.devices( HardwareID VARCHAR(255)  , primary key(HardwareID) );


CREATE TABLE telemetry.rpi4b_metrics (
    HardwareID VARCHAR(255),
    CPUuserLoad DOUBLE PRECISION,
    CPUidle DOUBLE PRECISION,
    TotalMemory BIGINT,
    FreeMemory BIGINT,
    IP VARCHAR(255),
    Temperature real,
    TimeStamp timestamp , constraint fk_HardwareID FOREIGN KEY (HardwareID) REFERENCES telemetry.devices(HardwareID));


CREATE TABLE telemetry.rpi_temp_alert (
    HardwareID VARCHAR(255),
    CPUuserLoad DOUBLE PRECISION,
    CPUidle DOUBLE PRECISION,
    TotalMemory BIGINT,
    FreeMemory BIGINT,
    IP VARCHAR(255),
    Temperature real,
    TimeStamp timestamp , constraint fk_HardwareID FOREIGN KEY (HardwareID) REFERENCES telemetry.devices(HardwareID));


```
