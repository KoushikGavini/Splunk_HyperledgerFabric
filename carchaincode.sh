#!/bin/bash

# copies chaincode configuration to vars folder and chaincode sourcecode
sudo cp -r car_chaincode vars/chaincode
sudo cp car_chaincode_collection_config.json vars

# installs car chaincode
echo "installing chaincode"
./minifab install -n car_chaincode -r true

# approves, commits and initalizes the chaincode
echo "Approving, commiting and initializing chaincode"
./minifab approve,commit,initialize -p ''

# inits cars chaincode
echo "init car chaincode........"
CAR=$( echo '{"name":"car1","color":"blue","tiresize":35,"owner":"jake","price":99}' | base64 | tr -d \\n )
./minifab invoke -p '"initCar"' -t '{"CAR":"'$CAR'"}'

CAR=$( echo '{"name":"car2","color":"red","tiresize":50,"owner":"jake","price":102}' | base64 | tr -d \\n )
./minifab invoke -p '"initCar"' -t '{"CAR":"'$CAR'"}'

CAR=$( echo '{"name":"car3","color":"blue","tiresize":70,"owner":"jake","price":103}' | base64 | tr -d \\n )
./minifab invoke -p '"initCar"' -t '{"CAR":"'$CAR'"}'

# transfers ownsership of car 
CAR_OWNER=$( echo '{"name":"car2","owner":"jerry"}' | base64 | tr -d \\n )
./minifab invoke -p '"transferCar"' -t '{"car_owner":"'$CAR_OWNER'"}'

# querys chaincode
echo "querying cars"
./minifab query -p '"readCar","car1"' -t ''
./minifab query -p '"readCarPrivateDetails","car1"' -t ''
./minifab query -p '"getCarByRange","car1","car4"' -t ''
