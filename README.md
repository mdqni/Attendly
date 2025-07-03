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

---

## 🧠 About Attendly

**Attendly** — это система для учёта посещаемости, построенная на микросервисной архитектуре.  
Каждое занятие генерирует уникальный QR-код, который студенты сканируют для отметки.  
Все события проходят через Kafka, данные хранятся в PostgreSQL и Redis.

---

## 📌 Features

- 🔐 JWT-аутентификация + RBAC-права доступа
- 🧾 Генерация временных QR-кодов (через Redis)
- 🧍 Отметка студента по QR + GPS
- 🧠 Микросервисная архитектура (gRPC + grpc-gateway)
- 💬 Асинхронная коммуникация через Kafka
- 📊 Интеграция с аналитикой и уведомлениями
- 🐳 Полностью контейнеризировано (Docker Compose)

---

## 🧱 Architecture

```text
                  +-------------+
                  |   Frontend  |
                  +------+------+
                         |
               REST (via grpc-gateway)
                         ↓
+---------+    +---------+    +----------+    +--------------+
|  Auth   |    |  User   |    |  Lesson  |    |  Attendance  |
| Service |<-> | Service |<-> | Service  |<-> |   Service    |
+----+----+    +---------+    +----------+    +--------------+
     |                                          ↑
     ↓                                          |
+----------------+                   +--------------------+
|     QR         |------------------>|  Redis (QR store)  |
|   Service      |                   +--------------------+
+----------------+
     ↓ Kafka
+----------------+
| Notification / |
| Analytics svc  |
+----------------+
````

---

## 🧪 Tech Stack

| Category       | Stack                                      |
| -------------- | ------------------------------------------ |
| Backend        | Go, gRPC, grpc-gateway                     |
| Database       | PostgreSQL                                 |
| Cache / Temp   | Redis                                      |
| Message Broker | Kafka                                      |
| Auth           | JWT, RBAC                                  |
| DevOps         | Docker, Docker Compose                     |
| Monitoring     | Prometheus, Grafana, Jaeger *(в процессе)* |

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

```bash
services/
├── auth/         # JWT, RBAC
├── user/         # Пользователи
├── group/        # Группы студентов
├── lesson/       # Пары / занятия
├── attendance/   # Отметка студентов
├── qr/           # Генерация и валидация QR
shared/           # Общие модули: Redis, Kafka, middleware и т.п.
proto/            # .proto схемы
```

---

## 🚧 In Progress

* [ ] Уведомления по Kafka
* [ ] Поддержка WebSocket для live-отметки
* [ ] CI/CD с GitHub Actions
* [ ] Дашборды в Grafana

---

## 🧑‍💻 Author

* Telegram: [@mdqni](https://t.me/mdqni)
* GitHub: [github.com/mdqni](https://github.com/mdqni)
