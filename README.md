### Collect system stats from raspberry and save it into a dockerized postgres database :
Reference : https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql


Folder structure :

#### step 1 : 
```
git clone https://github.com/devashishRaj/rip_telemetry.git
```

**then cd into client folder:**



#### step 2 : build executable , 
```
make build-all 
```
#### step 3: transfer build file to your raspy , eg: use scp if it's on same network 
```
scp <path to executable> <username>@<ip address>:<path to  save file on raspberry >
```
#### step 4 : run file on raspberry , ssh or login into raspi

```
<path to file>/<file name>
```

#### step 5 : run main.go on server side where database is setup , cd into server folder , if you have air live relaoad for goalng 
#### cd int server folder and type  " air "

#### OR 

__NOTE__ : make sure postgres is set up properly and viper config file is setup properly for right credentials and network info is present to make connection to database , refer Postgres.MD and https://github.com/spf13/viper

####IMPORTANT PART : https://github.com/spf13/viper#getting-values-from-viper


### viper guide 

**cd into server folder** 
```
makedir -p local/.config
```
cd into config 
```
vim config.json
```

sample json format for server side
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

tip : use zerotier for multiple devices existing on different lans for this project
      for to get a geoloaction visulization in grafana , you article mentioned below :
          https://medium.com/@tomjohnburton/visualising-a-distributed-network-c36871da52af
