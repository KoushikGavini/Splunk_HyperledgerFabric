#!/bin/bash
# Script to invoke chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer2.org0.example.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/org0.example.com/peers/peer2.org0.example.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=org0-example-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/org0.example.com/users/Admin@org0.example.com/msp
export ORDERER_ADDRESS=orderer3.example.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/ca.crt
peer chaincode invoke -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA \
  --tls -C mychannel -n privatemarbles  \
  --peerAddresses peer1.org0.example.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/org0.example.com/peers/peer1.org0.example.com/tls/ca.crt \
  --peerAddresses peer1.org1.example.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt \
  --transient '{"marble_owner":"eyJuYW1lIjoibWFyYmxlMiIsIm93bmVyIjoiamVycnkifQo="}' \
  -c '{"Args":["transferMarble"]}'
