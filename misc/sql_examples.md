####
Group by 

```
calculate the average CPU user load for each unique "macaddress":

SELECT macaddress, AVG(cpuuserload) AS avg_cpu_user_load
FROM telemetry.rpi4b_metrics
GROUP BY macaddress;

 average temperature for each unique "macaddress"

SELECT macaddress, AVG(temperature) AS avg_temperature
FROM telemetry.rpi4b_metrics
GROUP BY macaddress;

maximum memory usage for each unique "macaddress"

SELECT macaddress, MAX(memoryusage) AS max_memory_usage
FROM telemetry.rpi4b_metrics
GROUP BY macaddress;


```

ORDER BY , LIMIT 

```
 first 5 rows  ordered by "timestamp" column:

SELECT *
FROM telemetry.rpi4b_metrics
ORDER BY timestamp
LIMIT 5;

retrieve the last 10 rows 

SELECT *
FROM telemetry.rpi4b_metrics
ORDER BY timestamp DESC
LIMIT 10;

Retrieve 10 rows , starting from the 11th row, ordered by "memoryusage" in descending order:

SELECT *
FROM telemetry.rpi4b_metrics
ORDER BY memoryusage DESC
LIMIT 10 OFFSET 10;

```

Order By multiple rows :

```
retrieve the last 10 rows , order them  "timestamp" in descending order and "temperature" in ascending order:

SELECT *
FROM telemetry.rpi4b_metrics
ORDER BY timestamp DESC, temperature ASC
LIMIT 10;

retrieve the last 10 rows , order them by "timestamp" in descending order and "cpuuserload" in descending order:

SELECT *
FROM telemetry.rpi4b_metrics
ORDER BY timestamp DESC, cpuuserload DESC
LIMIT 10;

```
Use group by and order by :
```
group  by the "macaddress"  , calculate the average temperature for each unique id  , order the result set by the average temperature in descending order:

SELECT macaddress, AVG(temperature) AS avg_temperature
FROM telemetry.rpi4b_metrics
GROUP BY macaddress
ORDER BY avg_temperature DESC;


group the rows by "macaddress" and "timestamp"  ,calculate the maximum temperature for each unique combination of device and timestamp. order the result set by "macaddress" in ascending order and "timestamp" in descending order:

SELECT macaddress, timestamp, MAX(temperature) AS max_temperature
FROM telemetry.rpi4b_metrics
GROUP BY macaddress, timestamp
ORDER BY macaddress ASC, timestamp DESC;

retrieve the minimum temperature for each unique combination of device and timestamp, ordered by "macaddress" in descending order and "timestamp" in ascending order:

SELECT macaddress, timestamp, MIN(temperature) AS min_temperature
FROM telemetry.rpi4b_metrics
GROUP BY macaddress, timestamp
ORDER BY macaddress DESC, timestamp ASC;

```
find devices which have not send metric data in last 5min


```
with active as (
select DISTINCT macaddress
from telemetry.rpi4b_metrics
where timestamp > NOW() - INTERVAL '10 minute';)

Select macaddress
from telemetry.devices


```


