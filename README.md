# 📚 Attendly  
Scalable microservice-based attendance tracking system  

---

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
