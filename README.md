## MoonPay Educational Wallet Microservice

- Each service is independent and has its **own database**.
- Uses **event-driven patterns** to handle asynchronous operations.
- Focuses on **learning concepts** like Sagas, message passing, and service orchestration.

---

## Tech Stack

- **Language:** Go
- **Communication:** gRPC, REST (HTTP/JSON), WebSocket
- **Database:** MongoDB / PostgreSQL (can be configured per service)
- **Messaging:** RabbitMQ / Kafka (optional, for event-driven workflows)
- **Authentication:** JWT

---

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sidiik/MoonPay.git
   cd MoonPay
   ```
