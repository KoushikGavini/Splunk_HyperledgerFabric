#!/bin/bash

# copies chaincode configuration to vars folder
sudo cp privatemarbles_collection_config.json vars

# installs chaincode
echo "installing chaincode"
./minifab install -n privatemarbles -r true

# approves, commits and initalizes the chaincode
echo "Approving, commiting and initializing chaincode"
./minifab approve,commit,initialize -p ''

# inits marbles chaincode
echo "init marble chaincode........"
MARBLE=$( echo '{"name":"marble1","color":"blue","size":35,"owner":"tom","price":99}' | base64 | tr -d \\n )
./minifab invoke -p '"initMarble"' -t '{"marble":"'$MARBLE'"}'

MARBLE=$( echo '{"name":"marble2","color":"red","size":50,"owner":"tom","price":102}' | base64 | tr -d \\n )
./minifab invoke -p '"initMarble"' -t '{"marble":"'$MARBLE'"}'

MARBLE=$( echo '{"name":"marble3","color":"blue","size":70,"owner":"tom","price":103}' | base64 | tr -d \\n )
./minifab invoke -p '"initMarble"' -t '{"marble":"'$MARBLE'"}'

# transfers marbles 
MARBLE_OWNER=$( echo '{"name":"marble2","owner":"jerry"}' | base64 | tr -d \\n )
./minifab invoke -p '"transferMarble"' -t '{"marble_owner":"'$MARBLE_OWNER'"}'

# querys chaincode
echo "querying marbles"
./minifab query -p '"readMarble","marble1"' -t ''
./minifab query -p '"readMarblePrivateDetails","marble1"' -t ''
./minifab query -p '"getMarblesByRange","marble1","marble4"' -t ''
