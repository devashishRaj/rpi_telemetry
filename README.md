### Collect system stats from raspberry and save it into a dockerized postgres database :
Reference : https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql


Folder structure :
```
#### step 1 : 
```
git clone https://github.com/devashishRaj/rip_telemetry.git
```

then cd into client folder:

#### step2 : build executable , 
```
make build-all 
```

#### step3: transfer build file to your raspy , eg: use scp if it's on same network 

```
scp <path to executable> <username>@<ip address>:<path to  save file on raspberry >
```
#### step4: run file on raspberry , ssh or login into raspi

```
<path to file>/<file name>
```

#### step 5 : run main.go on server side where database is setup , cd into server folder , if you have air live relaoad for goalng 
#### cd int server folder and type  " air "

OR 

__NOTE__ : make sure postgres is set up properly and viper config file is setup porperly for right credentials and network info is present to make connection to database , refer Postgres.MD and https://github.com/spf13/viper

#### viper guide 

**cd into server folder** 

```
makedir -p local/.config

```
cd into config 

```
vim config.json
```

sample json format : 

```
{
    "postgresDB": {
        "host": "localhost",
        "port": "5432",
        "user": "xyz",
        "password": "xyz",
        "dbname": "xyz" ,
        "sslmode": "disable"
    }
}

```
