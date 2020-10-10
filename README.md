### Run the program
1. install go 1.15 on your computer
2. run `go run cmd/main.go`

### Send HTTP request
- the default port is 8088 
- temperature monitor url is `/tmp` 
- quality monitor url is `/qlt` 

### Example
send quality monitor request to `http://127.0.0.1:8088/qlt`  
Or you could use go tests in cmd/main_test.go  
**TestTemperature and TestQuality**    
these test will fire 1000 requests for 10 devices in parallel.    

### Story

##### IoT data collection

<p>sensors are being installed in an oil factory for data collection and processing
design a http server that could collect these big volume and high frequency data for the following two purposes:

1. real time anomaly detection
2. report generation


Sensor data format as follows:

Temperature sensor data format: `id,time,value;`

eg, T1,2020-01-30 19:30:40,27.65;

Oil quality sensor data format: `id,time,index1:value1,index2:value2;`

eg, Q1,2020-01-30 19:30:40,AB:38.9,AE:221323,CE:0.00001;

index names are abbreviated, including AB(Acid?)，AE(stickiness?), CE(water percentage?)

PS: T stands for temperature sensors, Q stands for quality sensors
</p>
 
##### Anomaly Detection

1. if one temperature sensor detects a temperature gap above 5 degrees within one day, send out a warning showing temperature being too high/low

2. if one quality sensor detects an index being 10% and above higher than last record twice straight, give a warning  


##### Report Generation

1. Daily average temperature  


##### Example

Input:
```
T1,2020-01-30 19:00:00,25;

T1,2020-01-30 19:00:01,22;

T1,2020-01-30 19:00:02,28;

Q1,2020-01-30 19:30:10,AB:37.8,AE:100,CE:0.01;

Q1,2020-01-30 19:30:20,AB:39.8,AE:100,CE:0.01;

Q1,2020-01-30 19:30:25,AB:39.9,AE:100,CE:0.01;

Q1,2020-01-30 19:30:32,AB:48.9,AE:101,CE:0.011;

Q1,2020-03-30 19:30:40,AB:58.9,AE:103,CE:0.012;
```


Output:

warnings:

`T1,2020-01-30 19:00:02,28; Temperature too high`

`Q1,2020-03-30 19:30:40,AB:58.9 AB too high`

Report:

`Temperature：2020-01-30 25.0`


##### Requirements:

1. fulfill input output requirement
2. good OO design
3. good extensibility
4. support high concurrency


