# Full Architectural, Technical, and Business Documentation: **Integra Enterprise & Open-Core (Go 1.26 Monorepo)**

> **Language**: [Polski (Polish)](../pl/README.md) | **English**  
> **Platform Version**: Go 1.26 Monorepo  
> **Hybrid Licensing**: Apache 2.0 (Open-Core) / Commercial Proprietary (Enterprise)  
> **Interactive Fullscreen Diagrams**: [`../DIAGRAMS-VIEWER.html`](../DIAGRAMS-VIEWER.html)  
> **Monolithic Single File**: [`../doc-en.md`](../doc-en.md)

---

## 📚 Modular Documentation Table of Contents

### [1. Purpose, Business Value, and Market Uniqueness in Go](01-purpose-and-business-value/README.md)
- [1.1. Project Purpose and Architectural Vision](01-purpose-and-business-value/1.1-project-purpose-and-vision.md)
- [1.2. Business Value](01-purpose-and-business-value/1.2-business-value.md)
- [1.3. Market Uniqueness in the Go Ecosystem](01-purpose-and-business-value/1.3-market-uniqueness.md)

### [2. Statistical Benchmark Analysis and Performance Justification](02-benchmark-analysis-and-performance/README.md)
- [2.1. Platform Descriptive Statistics](02-benchmark-analysis-and-performance/2.1-platform-statistics.md)
- [2.2. Performance Statistics by Monorepo Module](02-benchmark-analysis-and-performance/2.2-module-performance.md)
- [2.3. Benchmark Evidence of Architectural Uniqueness](02-benchmark-analysis-and-performance/2.3-architectural-evidence.md)
- [2.4. Benchmark Visualizations (Mermaid XY Charts)](02-benchmark-analysis-and-performance/2.4-benchmark-visualizations.md)

### [3. Global System Architecture and Mermaid Visualizations](03-global-architecture/README.md)
- [3.1. Ecosystem and Infrastructure Architecture](03-global-architecture/3.1-ecosystem-and-infrastructure.md)
- [3.2. Go Workspace Dependencies and Open-Core Split](03-global-architecture/3.2-go-workspace-and-open-core.md)
- [3.3. Internal Clean Architecture + CQRS of Microservices](03-global-architecture/3.3-clean-architecture-and-cqrs.md)

### [4. Key System Workflows](04-key-system-workflows/README.md)
- [4.1. Workflow 1: Multi-Tier JWT Token and Security State Verification](04-key-system-workflows/4.1-jwt-and-security-verification.md)
- [4.2. Workflow 2: Atomic Double-Entry Ledger Posting & Idempotency](04-key-system-workflows/4.2-double-entry-ledger-posting.md)
- [4.3. Workflow 3: Distributed User Registration Saga](04-key-system-workflows/4.3-user-registration-saga.md)
- [4.4. Workflow 4: CI/CD Quality Pipeline and Regression Gate](04-key-system-workflows/4.4-cicd-quality-pipeline.md)

### [5. Detailed Directory and File Documentation](05-directory-and-file-structure/README.md)
- [5.1. Root Monorepo, benchmarks/, examples/, and scripts/](05-directory-and-file-structure/5.1-root-monorepo-and-tooling.md)
- [5.2. Module .dist/integra-core (Exported Open-Core Distribution)](05-directory-and-file-structure/5.2-open-core-distribution.md)
- [5.3. Modules integra-proto and integra-postgres-adapters](05-directory-and-file-structure/5.3-proto-and-postgres-adapters.md)
- [5.4. Module integra-shared-go (Enterprise Shared Core)](05-directory-and-file-structure/5.4-shared-go-library.md)
- [5.5. Microservice integra-iam-service (Identity & Licenses)](05-directory-and-file-structure/5.5-iam-service.md)
- [5.6. Microservice integra-user-service (Profiles & Skills)](05-directory-and-file-structure/5.6-user-service.md)
- [5.7. Microservice integra-ledger-service (General Ledger & Invoicing)](05-directory-and-file-structure/5.7-ledger-service.md)
- [5.8. Module integra-infra (Cloud & Kubernetes)](05-directory-and-file-structure/5.8-cloud-infrastructure.md)


---

> [!TIP]
> **Interactive Fullscreen Diagrams**: A dedicated interactive viewer [DIAGRAMS-VIEWER.html](../DIAGRAMS-VIEWER.html) is available for exploring all system architecture diagrams with full-screen view (F11), smooth wheel zoom, and drag & pan navigation.
