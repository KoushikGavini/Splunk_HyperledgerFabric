# Hyperledger Fabric Integration with Splunk

## Objective

* The purpose of this repo is to illustrate splunk transaction logging integration with Hyperledger Fabric. The README will also analysis the pro's and cons of the splunk fabric-logger.

* Table Of Contents
    - [Analysis: Pro's and Con's](#analysis)
    - [Working demo with minifab and car ownership chaincode](#working-demo)
    - [Resources](#resources)

### Analysis: 

* Pro's/What I like, that is already built. 
    - I liked the ease of setting up the splunk fabric-logger with both native hyperledger fabric and minifab. By utilizing docker containers
    - Has holistic monitoring/logging aggregator service from infrastructure to blockchain network, so an engineer does not have to setup grafana, prometheus, and hyperledger explorer. 

* Con's/What I would want and how to build it with splunk. 
    - Certificate Monitoring [**Most desired feature**]
        - One of the key issues and critical bugs that Enterprises are facing when in production is the issue of expiring certs. What is happening is that when certs expires such as the admin certs, the organization gets shut out. I would like if splunk can leverage `fabric-ca-server` and peer logs to notify the organization of an upcoming cert expiration and if possible automated reissuance of a cert. This may be possible because *all communication with fabric-ca-server is done through RESTFUL APIs*. 
        - Here is a [link](https://lists.hyperledger.org/g/fabric/topic/criticial_admin_certificate/71743922?p=,,,20,0,0,0::recentpostdate%2Fsticky,,,20,2,0,71743922) to common critical issue faced by many organizations running Hyperledger Fabric or Enterprise Blockchain Applications. 

        - Below is an illustation of solution architecture with splunk integration in regards to fabric-ca-server. 

        ![Splunk_Fabric_CA](Splunk_Fabric_CA.jpeg)
    - When I used fabric-logger with native hyperledger fabric setup replicating a production network. I would of liked these information on the chaincode commit status when debugging or monitoring, 
        - Check commit status of a chaincode (ie: `peer lifecycle chaincode checkcommitreadiness`)
        - This is useful when orgs are on production scale multi machine env and you want to check which org is rejecting it and why. Especially with 2.x since there is a more decentralized approach to transaction flow. 
    - Monitor of fabric network event listener 
        - I would like a bit more monitoring of the `fabric-network-event-listener` as it is commonly used in production grade application applications such as in integrating off-chain-data analysis and oracle development. 

### Working Demo: 

#### Splunk transaction integration with Hyperledger Fabric working example using minifab and car ownership chaincode

* Note: Please use legacy osxfs file sharing mechanism for MAC. [link](https://github.com/hyperledger-labs/minifabric/issues/141)
* Example using HF 2.2

* Step 1: Run the start script 

> ./start.sh

* Step 2: Run the privatemarbles script


*The privatemarbles script will configure, init, approve, invoke, and query the privatemarbles chaincode.*

> ./privatemarbles.sh

* Step 3: Open port 8000 on local host or VM ip

```text
https://localhost:8000
Username: admin
Password: changeme
```
* Step 4: Terminate the example

> ./stop.sh


### Resources: 

* https://github.com/hyperledger-labs/minifabric
* https://github.com/splunk/fabric-logger
* https://github.com/splunk/fabric-logger/tree/master/examples/minifab
* https://github.com/hyperledger/fabric-samples/tree/main/chaincode/marbles02_private
* https://github.com/splunk/fabric-logger/tree/master/examples/vaccine-demo




