# pong-online

## Requirements
* Kubernetes 1.9
* Agones installed on the cluster ([see](https://github.com/GoogleCloudPlatform/agones/blob/master/docs/installing_agones.md))

## Installation
```
git clone https://github.com/zaymat/pong-online
```
First, build the docker image 
```
cd ./server
docker build -t pong .
```
and then create the gameserver on Kubernetes
```
cd ./kubernetes
kubectl apply -f server.yaml
```

## Goal
This is a proof of concept of using Agones plugin for Kubernetes.
We will create a pong server and a pong client to use the Google Gameserver plugin for Kubernetes : Agones [https://github.com/GoogleCloudPlatform/agones]

## Agones
Few weeks ago, Google announced his new plugin for Kubernetes, developed with Ubisoft : Agones. This plugin allows to manage and scale gameservers, and I found interesting trying to developed a homemade gameserver and running it on Kubernetes.

## Server
In addition to running the game, the server implement Agones SDK whick perform healthcheck and shutdown automation.

## Client
The client is written in Go and graphics are managed by the SDL library. There is no logic in the client, only message handling and window drawing.

## State of progress

The game is playable with client and server, and the server implement the Agones SDK and run on Kubernetes.

## Improvements

* The GUI of the client is very poor. It would be cool to improve it a little bit
* Work on ball speed and racket speed
* Developed a matchmaking server, which scale the kubernetes gameservers
* Create JS client (more is better)
