#!/bin/bash

# copies chaincode configuration to vars folder and chaincode sourcecode
sudo cp privateassets_collection_config.json vars
sudo cp -r privateassets vars/chaincode

# installs chaincode
echo "installing chaincode"
./minifab install -n privateassets -r true

# approves, commits and initalizes the chaincode
echo "Approving, commiting and initializing chaincode"
./minifab approve,commit,initialize -p ''

# inits asset chaincode
echo "init asset chaincode........"
ASSET=$( echo '{"name":"asset1","color":"blue","size":35,"owner":"tom","price":99}' | base64 | tr -d \\n )
./minifab invoke -p '"initAsset"' -t '{"asset":"'$ASSET'"}'

ASSET=$( echo '{"name":"asset2","color":"red","size":50,"owner":"tom","price":102}' | base64 | tr -d \\n )
./minifab invoke -p '"initAsset"' -t '{"asset":"'$ASSET'"}'

ASSET=$( echo '{"name":"asset3","color":"blue","size":70,"owner":"tom","price":103}' | base64 | tr -d \\n )
./minifab invoke -p '"initAsset"' -t '{"asset":"'$ASSET'"}'

# transfers asset 
ASSET_OWNER=$( echo '{"name":"asset2","owner":"jerry"}' | base64 | tr -d \\n )
./minifab invoke -p '"transferAsset"' -t '{"asset_owner":"'$ASSET_OWNER'"}'

# querys chaincode
echo "querying asset"
./minifab query -p '"readAsset","asset1"' -t ''
./minifab query -p '"readAssetPrivateDetails","asset1"' -t ''
./minifab query -p '"getAssetsByRange","asset1","asset4"' -t ''