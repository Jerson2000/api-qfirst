## QFirst: Skip the Wait, Embrace the Ease of Booking

**QFirst** is a robust and easy-to-use booking system designed to help businesses manage appointments and reservations efficiently. Built with simplicity and scalability in mind, **QFirst** allows users to schedule, view, and manage bookings with minimal effort. The system is ideal for small and medium-sized enterprises looking for an intuitive tool to streamline their booking processes.

### To-Do List

- [x] **JWT Middleware**: Implemented JWT-based authentication to secure the application.
- [x] **Refresh Token**: Implemented a refresh token system for token renewal.
- [x] **Role-based Permission**: Implemented a role-based permission system for controlling user access. (RBAC with [Casbin](https://github.com/casbin/casbin))
- [x] **CSRF Middleware**: Implemented partial CSRF protection for secure requests.
- [x] **Cache Middleware**: Implemented caching for faster data retrieval.
- [x] **Swagger UI**: Implemented Swagger, allowing users to visualize and test API endpoints directly within the interface.
- [x] **WebSocket**: Implemented WebSocket for real-time chat and other real-time features (pub/sub).
- [x] **Mailer**: Implemented a mailer using [Go-Mail](https://github.com/wneessen/go-mail).
- [x] **OTP Verification**: Implement OTP-based verification for users.
- [x] **File Upload**: Implement file upload.

### Tech Stack

- **Go stdlib (net/http)**
- **Gorm**
- **🦍/mux for routing**
- **🦍/websocket**
- **etc..** (for more info [go.mod](https://github.com/Jerson2000/api-qfirst/blob/master/go.mod))

---
