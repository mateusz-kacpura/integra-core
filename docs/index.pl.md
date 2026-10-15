# Pełna dokumentacja architektoniczna, techniczna i biznesowa platformy **Integra Enterprise & Open-Core (Go 1.26 Monorepo)**

> **Wersja:** Go 1.26 Monorepo  
> **Licencja hybrydowa:** Apache 2.0 (Open-Core) / Proprietary Commercial (Enterprise)  
> **Przeglądarka schematów pełnoekranowych:** [`DIAGRAMS-VIEWER.html`](DIAGRAMS-VIEWER.html)  
> **Jednolity plik zbiorczy:** [`doc.pl.md`](doc.pl.md)

---

## 📚 Spis treści dokumentacji modułowej

### [1. Cel projektu, wartość biznesowa i unikalność na rynku Go](01-purpose-and-business-value/index.md)

- [1.1. Cel projektu i wizja architektoniczna](01-purpose-and-business-value/1.1-project-purpose-and-vision.md)
- [1.2. Wartość biznesowa (Business Value)](01-purpose-and-business-value/1.2-business-value.md)
- [1.3. Unikalność w skali rynku projektów w Go](01-purpose-and-business-value/1.3-market-uniqueness.md)

### [2. Analiza statystyczna benchmarków i uzasadnienie wydajnościowe](02-benchmark-analysis-and-performance/index.md)

- [2.1. Statystyki opisowe całej platformy](02-benchmark-analysis-and-performance/2.1-platform-statistics.md)
- [2.2. Statystyki wydajności w podziale na moduły monorepo](02-benchmark-analysis-and-performance/2.2-module-performance.md)
- [2.3. Dowody benchmarkowe na unikalność architektoniczną](02-benchmark-analysis-and-performance/2.3-architectural-evidence.md)
- [2.4. Wizualizacje wyników benchmarków (Mermaid XY Charts)](02-benchmark-analysis-and-performance/2.4-benchmark-visualizations.md)

### [3. Globalna architektura systemu i wizualizacje Mermaid](03-global-architecture/index.md)

- [3.1. Architektura ekosystemu i infrastruktury](03-global-architecture/3.1-ecosystem-and-infrastructure.md)
- [3.2. Architektura zależności modułów Go Workspace (go.work) i podział Open-Core](03-global-architecture/3.2-go-workspace-and-open-core.md)
- [3.3. Wewnętrzna architektura czysta / heksagonalna + CQRS mikroserwisów](03-global-architecture/3.3-clean-architecture-and-cqrs.md)

### [4. Kluczowe przepływy pracy (System Workflows)](04-key-system-workflows/index.md)

- [4.1. Przepływ 1: Wielopoziomowa weryfikacja tokena JWT i stanu bezpieczeństwa (L1 -> L2 -> L3)](04-key-system-workflows/4.1-jwt-and-security-verification.md)
- [4.2. Przepływ 2: Atomowe księgowanie podwójnego zapisu (Double-Entry Ledger Posting & Idempotency)](04-key-system-workflows/4.2-double-entry-ledger-posting.md)
- [4.3. Przepływ 3: Rozproszona saga rejestracji użytkownika i synchronizacji replik odczytu](04-key-system-workflows/4.3-user-registration-saga.md)
- [4.4. Przepływ 4: Potok jakości CI/CD, audyt bezpieczeństwa i bramka regresji benchmarków](04-key-system-workflows/4.4-cicd-quality-pipeline.md)

### [5. Szczegółowa dokumentacja katalogów i plików projektu](05-directory-and-file-structure/index.md)

- [5.1. Katalog główny (Root Monorepo), benchmarks/, examples/ oraz scripts/](05-directory-and-file-structure/5.1-root-monorepo-and-tooling.md)
- [5.2. Moduł .dist/integra-core (wyeksportowana dystrybucja Open-Core)](05-directory-and-file-structure/5.2-open-core-distribution.md)
- [5.3. Moduł integra-proto i integra-postgres-adapters](05-directory-and-file-structure/5.3-proto-and-postgres-adapters.md)
- [5.4. Moduł integra-shared-go (biblioteka współdzielona Enterprise)](05-directory-and-file-structure/5.4-shared-go-library.md)
- [5.5. Mikroserwis integra-iam-service (zarządzanie tożsamością, rolami i licencjami)](05-directory-and-file-structure/5.5-iam-service.md)
- [5.6. Mikroserwis integra-user-service (profile użytkowników, konta i kompetencje)](05-directory-and-file-structure/5.6-user-service.md)
- [5.7. Mikroserwis integra-ledger-service (księga główna podwójnego zapisu, Event Sourcing, faktury i wypłaty)](05-directory-and-file-structure/5.7-ledger-service.md)
- [5.8. Moduł integra-infra (infrastruktura chmurowa, Kubernetes, Helm, Terraform, WAF i monitoring)](05-directory-and-file-structure/5.8-cloud-infrastructure.md)


---

> [!TIP]
> **Interaktywny podgląd na pełnym ekranie**: Dostępna jest dedykowana przeglądarka [DIAGRAMS-VIEWER.html](DIAGRAMS-VIEWER.html) z obsługą trybu pełnoekranowego (przycisk *Pełny ekran* lub klawisz `F11`), płynnego przybliżania kółkiem myszy (zoom) oraz przesuwania (drag & pan) wszystkich 9 schematów architektonicznych i przepływów.
