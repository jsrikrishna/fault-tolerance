
# Load Balancer
 - Periodic Healthechks are performed for the status of the Server
	 - HealthCheck interval can be specified in the [config.json](https://github.com/jsrikrishna/fault-tolerance/blob/master/config/config.json)
	 - `status_counter` indicates the maximum number of heath checks to be performed once a server is down.
	 - `pingInterval` indicates that max amount of time the load-balancer waits for a pong from the server
 - Load Balancing Strategies
	 - `cpumetrics`, `random`, `roundrobin`, `weightedroundrobin` can be specified as the load balancing strategy in [config.json](https://github.com/jsrikrishna/fault-tolerance/blob/master/config/config.json)
	 - `cpumetrics` are collected in load-balancer by invoking the API `/systemResources` exposed by the server for every `pingInterval` time
	 - Information about new servers can be added at runtime by invoking the API `/server` 
		```
		POST  /server
		Host: load-balancer
		Body: 
		{
			"serverName": "server3",
			"address": "localhost:8084",
			"weight": 2
		}
		```
	- If a server goes down, details about it are automically updated in load-balancer and request are not routed to it.
 - Status Completion of Continuous Queries
	 - API `/status` is exposed in the load-balancer to mark the completion of the continuous queries by the servers. Continuous queries of type `/resources` are only remembered by the load-balancer (to configure them via config.json can be implemented).
	 ```
	 POST /status
	 Host: load-balancer
	 Body: 
	 {
		"backend" : "localhost:8081",
		"starttime": "Tue, 06/06/17, 08:19AM",
		"endtime": "Wed, 06/07/17, 10:50AM"
	 }
	 ```

Team Members : Sri Krishna Jaliparthy, Raghu Sai Gudipati, Thejdeep Gudivada
