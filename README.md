# Launchpad Agent

![](Launchpad-Agent.png)

A daemon that runs in the webserver that is responsible for:
1. Updating an application's static artifacts located in the webserver's public directory
2. Running a websocket service that can broadcast `AppUpdateEvent` to the connected clients.

## Architecture

## Motivation
A `portal app` can be added, removed or updated anytime. The `Portal Platform` should be able to 
handle these changes at runtime without affecting other `portal apps`.

### Solution
To achieve the above requirement, we hereby introduce `The Launchpad Agent`. This agent will run as a background process in the webserver's container.

This agent listens for `PortalAppCompileCompletedEvent`, downloads the compiled `portal app` artifacts from s3 and extracts them into the webserver's public directory.

To inform connected clients that a new version of a `portal app` has been made available, this agent runs a websocket at ws://portalws. The server reverse proxies requests to ws://portalws to the agent's websocket server running on the same container.

In order to minimize the footprint of this agent on the webserver a programming language that compiles to native have been choosen.

## Package Description
This section describes the top-level packages and their purpose.

|S/No|Package Name|Description                                        |Remarks|
|----|:-----------|---------------------------------------------------|-------|
|1   |admin       |Handles bi-directional communication with the admin server||
|2   |cfg         |Handles application configuration such as parsing of the startup configuration file and populating relevant data structures||
|3   |messaging   |Handles interaction with various broker providers. e.g. Kafka, RabbitMQ, ZeroMQ, ActiveMQ, etc|Currently, Only Kafka is supported|
|4   |servers     |Handles interaction with various WebServer implementations. e.g. Nginx, Apache Httpd, etc|Currently, Only Nginx is supported|

## Getting Started

1. [See Compiling Protocol Buffers](https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers)

2. Generating Protobuf Stubs

    ```bash
    $ protoc -I=./launchpad-schema -I=./launchpad-schema/include --go_out=./launchpad-agent launchpad-schema/metrics.proto
    ```

## Starting Locally

1. Run Kafka
    ```bash
   $ docker-compose up -d broker 
   ```

2. Run Nexus
    ```bash
    $ docker-compose up -d nexus
    ```

## Troubleshooting

### Debugging on Docker
Please refer to [this article](https://blog.jetbrains.com/go/2020/05/06/debugging-a-go-application-inside-a-docker-container/)

### MacOS
Stop "developer tools access needs to take control of another process for debugging to continue" alert

```bash
$ sudo /usr/sbin/DevToolsSecurity --enable
```
