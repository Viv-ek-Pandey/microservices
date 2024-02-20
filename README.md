This project demonstrates a set of 6 Go microservices communicating with different protocols: HTTP, RPC, gRPC, and AMQP. Docker is used for containerization and deployments, while external services like MongoDB, PostgreSQL, RabbitMQ, and MailHog are leveraged for data storage, messaging, and email capabilities.
Technologies Used

    Programming Language: Go
    Microservice Frameworks: (depending on service)
        HTTP: Echo framework
        RPC: Go Micro
        gRPC: grpc-go
        AMQP: amqp
    Containerization: Docker
    External Services:
        MongoDB: Data storage
        PostgreSQL: Data storage
        RabbitMQ: Message broker
        MailHog: Email server

Architecture

The project follows a microservices architecture where each service focuses on a specific domain and communicates with others through defined interfaces. The chosen communication protocols provide flexibility and cater to different integration needs:

    HTTP: Public APIs and RESTful communication.
    RPC: Remote procedure calls for internal service interactions.
    gRPC: High-performance, bi-directional communication between services.
    AMQP: Asynchronous messaging for distributed tasks and event-driven workflows.

Docker Usage

Each microservice resides in its own Docker container, enabling independent development, deployment, and scaling. Additionally, Docker Compose can be used to manage the orchestration of all services and external dependencies.
External Services

    MongoDB: NoSQL database for storing unstructured data.
    PostgreSQL: Relational database for structured data.
    RabbitMQ: Message broker for asynchronous communication and event handling.
    MailHog: Local email server for testing and development purposes.
