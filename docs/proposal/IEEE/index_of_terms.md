# Index of Terms

| Term | Definition |
| :--- | :--- |
| **ACID** | Atomicity, Consistency, Isolation, Durability; a set of properties of database transactions intended to guarantee data validity despite errors, power failures, and other mishaps. |
| **Anonymous-first** | A design philosophy adopted by Cthulhu that prioritizes user anonymity, allowing full access to core file sharing features without requiring user account creation or personal data collection. |
| **API Gateway** | A centralized server that acts as the single entry point for the Cthulhu system, responsible for routing client requests to appropriate microservices, handling request validation, and managing initial security checks. |
| **ClamAV** | An open-source antivirus engine integrated into Cthulhu's Security Service to perform server-side scanning of uploaded files for malicious content such as viruses, trojans, and ransomware. |
| **CQRS** | Command Query Responsibility Segregation; a pattern that separates read and update operations for a data store, often used in microservices to optimize performance and scalability. |
| **Event-Driven Architecture** | A software architecture paradigm promoting the production, detection, consumption of, and reaction to events, utilized by Cthulhu via RabbitMQ to decouple services. |
| **File Sharing Platform** | A digital service enabling users to upload, store, and distribute files. Cthulhu operates as a privacy-focused alternative to mainstream services like WeShare, OneDrive, and Google Drive, distinguishing itself by eliminating mandatory account creation and enforcing temporal data life cycles. |
| **Golang (Go)** | A statically typed, compiled programming language designed at Google, used for building Cthulhu's high-performance backend microservices. |
| **Microservices** | An architectural style where the application is structured as a collection of loosely coupled, independently deployable services, each focused on a specific business capability (e.g., File Manager, Auth, Security). |
| **Next.js** | A React framework used to build Cthulhu's frontend, enabling a responsive and modern user interface for file uploads and management. |
| **RabbitMQ** | An open-source message-broker software that facilitates asynchronous inter-service communication within Cthulhu, ensuring scalable and reliable message passing. |
| **SAGA Pattern** | A design pattern used to manage data consistency across microservices in distributed transactions, ensuring that if one step fails, compensating transactions are triggered to rollback changes. |
| **SHA-256** | Secure Hash Algorithm 256-bit; a cryptographic hash function used by Cthulhu to generate unique, anonymized fingerprints for uploaded files, separate from their content. |
| **Unrestricted File Upload (UFU)** | A critical web security vulnerability where an application allows users to upload files without sufficient validation of type, size, or content, which Cthulhu mitigates through strict validation policies. |