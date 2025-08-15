## MoonPay Educational Wallet Microservice

- Each service is independent and has its **own database**.
- Uses **event-driven patterns** to handle asynchronous operations.
- Focuses on **learning concepts** like Sagas, message passing, and service orchestration.

## Contributors

- [sidiik](https://github.com/sidiik) - Sidiiq Omar (sidiik)

## Tech Stack

- **Language:** Go
- **Communication:** gRPC, REST (HTTP/JSON), WebSocket
- **Database:** MongoDB / PostgreSQL (can be configured per service)
- **Messaging:** RabbitMQ / Kafka (optional, for event-driven workflows)
- **Authentication:** JWT

## Getting Started

1. **Install Dependencies:**

   - Ensure you have Go installed.
   - Install Tilt for local development:
     ```bash
     1. curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
     2. tilt version
     ```
   - Install Docker for container management.
   - Minikube or Kubernetes cluster for local development.
   - kubectl command-line tool for interacting with your cluster.

2. **Clone the repository:**

   ```bash
   1. git clone https://github.com/sidiik/MoonPay.git
   2. cd MoonPay
   3. tilt up - to start the services
   4. tilt down - to stop the services
   ```
