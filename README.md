# Tasks API

## Description

This app helps track maintenance work done during the day and is made for two types of users: Managers and Technicians. Technicians can do tasks and only look at, add, or change the tasks they did. Each task includes a detailed description (up to 2500 characters) and the date it was done. 

Managers can see all the tasks done by technicians, remove any tasks, and get notified when a task is completed. This way, managers always know what maintenance work is being done.

### How can I run this project?

First you need to build the images from the API and kafka consumer, run `make build-containers`.

After that you can run `make start` that will initiate the following containers:

- [APP] Rest API for tasks & users;
- [APP] Kafka consumer for the `managerNotification` topic;
- [Infrastructure] MySQL database that stores users and tasks;
- [Infrastructure] Zookeeper Kafka metadata manager;
- [Infrastructure] The Kafka;
- [Infrastructure] [Control Center](http://localhost:9021) to get an interface to help manage the Kafka's traffic.

Image sizes for the apps should be around this:

| Image name                       | Tag    | Size   |
|----------------------------------|--------|--------|
| cristovaoolegario/tasks-consumer | latest | 15.7MB |
| cristovaoolegario/tasks-api      | latest | 32.5MB |

You can check this locally after you build the containers by running ` docker images -a | grep cristovaoolegario/tasks`
