#!/bin/bash
# Script to create channel block 0 and then create channel
cp $FABRIC_CFG_PATH/core.yaml /vars/core.yaml
cd /vars
export FABRIC_CFG_PATH=/vars
configtxgen -profile OrgChannel \
  -outputCreateChannelTx mychannel.tx -channelID mychannel

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer2.org0.example.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/org0.example.com/peers/peer2.org0.example.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=org0-example-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/org0.example.com/users/Admin@org0.example.com/msp
export ORDERER_ADDRESS=orderer3.example.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/ca.crt
peer channel create -c mychannel -f mychannel.tx -o $ORDERER_ADDRESS \
  --cafile $ORDERER_TLS_CA --tls
