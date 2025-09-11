# Taskkr

### Task Management Microservice

This microservice provides operations(CRUD) over task resource


### Running the service

```bash
#create .env file using the .env.wxample
docker-compose up --build -d
# the containers for service and db will start
```

APIs docs - http://localhost:8080/swagger/index.html

### Design descisions 

This service only handles the tasks entities(CRUD operations) as per requirements.
Auth is not included because its better to have a centralized separate auth service rather than each service
having its own, or could be handled at API Gateway(if its part of infra).
For service to service communications if service mesh is deployed mtls can be used, and specifically for auth using
tokens, a sidecar can be used which handles the auth.

This service has a Postgres DB on which it persists the data.

This service can be scaled horizontally as per the load dynmically using HPA on k8s, but need to keep database scalability and perfomrance in check as well, adding replicas for reads would help, also partitioning the data will be useful at larger scales.

### Connecting to other microserviecs

Services could connect either via rest or grpc, grpc would be faster and more advisable for service to service communication.
In case the internal communication is a lot having a service mesh would be nice it offers lot of features like traffic management, security, observability etc.
For async communications between services we could use a message queue(rabbitmq, kafka) which would enable services to pass events and messages.

### Improvements

Add a proper logger, providing metrics, handling concurrent write requests maintaining consistency and availability, adding more detailed checks for validations. Add tests unit and e2e.