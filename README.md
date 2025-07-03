<h1 align="center">📚 Attendly</h1>
<h3 align="center">Scalable microservice-based attendance tracking system</h3>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" />
  <img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" />
  <img src="https://img.shields.io/badge/gRPC-0078D7?style=for-the-badge&logo=grpc&logoColor=white" />
  <img src="https://img.shields.io/badge/Kafka-000000?style=for-the-badge&logo=apache-kafka&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
</p>

## 🧠 About Attendly  
**Attendly** is a scalable attendance tracking system built with a microservice architecture.  
Each class generates a unique QR code that students scan to mark their attendance.  
All events go through Kafka, with data stored in **PostgreSQL** and **Redis**.

---

## 📌 Features  
- 🔐 JWT authentication + RBAC access control  
- 🧾 Temporary QR code generation (stored in Redis)  
- 🧍 Student check-in via QR + GPS  
- 🧠 Microservices architecture (gRPC + grpc-gateway)  
- 💬 Asynchronous communication through Kafka  
- 📊 Integration with analytics and notification services  
- 🐳 Fully containerized using Docker Compose  

---

## 🧱 Architecture  

```text
                 +-------------+
                 |   Frontend  |
                 +------+------+
                        |
              REST (via grpc-gateway)
                        ↓
+---------+   +---------+   +----------+   +--------------+
|  Auth   |   |  User   |   |  Lesson  |   |  Attendance  |
| Service |<->| Service |<->| Service  |<->|   Service    |
+----+----+   +---------+   +----------+   +--------------+
     |                                       ↑
     ↓                                       |
+----------------+               +--------------------+
|     QR         |-------------->|  Redis (QR store)  |
|   Service      |               +--------------------+
+----------------+
     ↓ Kafka
+----------------+
| Notification / |
| Analytics svc  |
+----------------+
````

---

## 🧪 Tech Stack

| Category       | Stack                                     |
| -------------- | ------------------------------------------|
| Backend        | Go, gRPC, grpc-gateway                    |
| Database       | PostgreSQL                                |
| Cache / Temp   | Redis                                     |
| Message Broker | Kafka                                     |
| Auth           | JWT, RBAC                                 |
| DevOps         | Docker, Docker Compose                    |
| Monitoring     | Prometheus, Grafana, Jaeger (In progress) |

---

## 🛠️ Running Locally

```bash
# Clone the repo
git clone https://github.com/mdqni/Attendly.git
cd Attendly

# Start everything
docker-compose up --build
```

---

## 📁 Service Structure

```text
services/
├── auth/         # JWT, RBAC
├── user/         # Users
├── group/        # Student groups
├── lesson/       # Lessons
├── attendance/   # Student attendance tracking
├── qr/           # QR code generation & validation
shared/           # Common modules: Redis, Kafka, middleware, etc.
proto/            # .proto definitions
```

---

## 🚧 In Progress

[] Notifications via Kafka
[] WebSocket support for real-time check-ins
[] CI/CD with GitHub Actions
[] Grafana dashboards for monitoring

---

## 🧑‍💻 Author

* Telegram: [@mdqni](https://t.me/mdqni)
* GitHub: [github.com/mdqni](https://github.com/mdqni)
