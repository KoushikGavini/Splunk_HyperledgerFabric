#!/bin/bash
# Script to instantiate chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer1.org0.example.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/org0.example.com/peers/peer1.org0.example.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=org0-example-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/org0.example.com/users/Admin@org0.example.com/msp
export ORDERER_ADDRESS=orderer2.example.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/ca.crt
peer chaincode query -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA \
  --tls -C mychannel -n privatemarbles  \
  --peerAddresses peer1.org0.example.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/org0.example.com/peers/peer1.org0.example.com/tls/ca.crt \
  -c '{"Args":["getMarblesByRange","marble1","marble4"]}'
