# Integra Platform – Dokumentacja Techniczna i Architektoniczna / Technical Documentation

Witamy w centrum dokumentacji platformy **Integra Enterprise & Open-Core (Go 1.26 Monorepo)**.  
Dokumentacja jest dostępna w dwóch językach (język polski oraz język angielski) oraz w formacie interaktywnego portalu internetowego wdrażanego automatycznie na **GitHub Pages**.

---

## 🌐 Wybierz wersję językową / Choose Language

| Język / Language | Dokumentacja modularna (rozdziały) | Plik jednolity (zbiorczy) |
| :--- | :--- | :--- |
| 🇵🇱 **Polski (Polish)** | [**Przejdź do rozdziałów (docs/pl/)**](pl/README.md) | [**`docs/doc-pl.md`**](doc-pl.md) |
| 🇬🇧 **English (Angielski)** | [**Browse Chapters (docs/en/)**](en/README.md) | [**`docs/doc-en.md`**](doc-en.md) |

---

## ⚡ Interaktywny portal internetowy & schematy

- **Portal Dokumentacji (Web Portal)**: [`docs/index.html`](index.html) *(wdrażany przez GitHub Actions na GitHub Pages z natywnym przełącznikiem PL/EN, wyszukiwarką pełnotekstową i podglądem na żywo)*.
- **Interaktywna przeglądarka schematów (Fullscreen)**: [**`docs/DIAGRAMS-VIEWER.html`**](DIAGRAMS-VIEWER.html) *(9 pełnoekranowych schematów Mermaid z zoomem kółkiem myszy, przeciąganiem i eksportem do SVG)*.

---

## 📚 Spis treści / Table of Contents

### 1. Cel projektu i wartość biznesowa / Purpose & Business Value
- [1.1. Cel projektu i wizja architektoniczna](pl/01-cel-i-wartosc-biznesowa/1.1-cel-projektu-i-wizja-architektoniczna.md) | [EN](en/01-purpose-and-business-value/1.1-project-purpose-and-vision.md)
- [1.2. Wartość biznesowa (Business Value)](pl/01-cel-i-wartosc-biznesowa/1.2-wartosc-biznesowa.md) | [EN](en/01-purpose-and-business-value/1.2-business-value.md)
- [1.3. Unikalność w skali rynku Go](pl/01-cel-i-wartosc-biznesowa/1.3-unikalnosc-na-rynku-go.md) | [EN](en/01-purpose-and-business-value/1.3-market-uniqueness.md)

### 2. Analiza statystyczna benchmarków / Benchmark Analysis
- [2.1. Statystyki opisowe całej platformy](pl/02-analiza-benchmarkow-i-wydajnosc/2.1-statystyki-opisowe-platformy.md) | [EN](en/02-benchmark-analysis-and-performance/2.1-platform-statistics.md)
- [2.2. Statystyki wydajności w podziale na moduły](pl/02-analiza-benchmarkow-i-wydajnosc/2.2-statystyki-wydajnosci-modulow.md) | [EN](en/02-benchmark-analysis-and-performance/2.2-module-performance.md)
- [2.3. Dowody benchmarkowe na unikalność architektoniczną](pl/02-analiza-benchmarkow-i-wydajnosc/2.3-dowody-benchmarkowe-wzorce-i-anomalie.md) | [EN](en/02-benchmark-analysis-and-performance/2.3-architectural-evidence.md)
- [2.4. Wizualizacje wyników benchmarków (Mermaid XY Charts)](pl/02-analiza-benchmarkow-i-wydajnosc/2.4-wizualizacje-wykresy-mermaid.md) | [EN](en/02-benchmark-analysis-and-performance/2.4-benchmark-visualizations.md)

### 3. Globalna architektura systemu / Global Architecture
- [3.1. Architektura ekosystemu i infrastruktury](pl/03-architektura-globalna/3.1-ekosystem-i-infrastruktura.md) | [EN](en/03-global-architecture/3.1-ecosystem-and-infrastructure.md)
- [3.2. Architektura zależności modułów Go Workspace (go.work) i podział Open-Core](pl/03-architektura-globalna/3.2-zaleznosci-modulow-go-work-i-open-core.md) | [EN](en/03-global-architecture/3.2-go-workspace-and-open-core.md)
- [3.3. Wewnętrzna architektura czysta / heksagonalna + CQRS mikroserwisów](pl/03-architektura-globalna/3.3-architektura-czysta-i-cqrs-mikroserwisow.md) | [EN](en/03-global-architecture/3.3-clean-architecture-and-cqrs.md)

### 4. Kluczowe przepływy pracy / Key System Workflows
- [4.1. Przepływ 1: Weryfikacja tokena JWT i stanu bezpieczeństwa (L1-L3)](pl/04-kluczowe-przeplywy-pracy/4.1-wielopoziomowa-weryfikacja-tokena-jwt.md) | [EN](en/04-key-system-workflows/4.1-jwt-and-security-verification.md)
- [4.2. Przepływ 2: Atomowe księgowanie podwójnego zapisu](pl/04-kluczowe-przeplywy-pracy/4.2-atomowe-ksiegowanie-podwojnego-zapisu.md) | [EN](en/04-key-system-workflows/4.2-double-entry-ledger-posting.md)
- [4.3. Przepływ 3: Rozproszona saga rejestracji użytkownika](pl/04-kluczowe-przeplywy-pracy/4.3-rozproszona-saga-rejestracji-uzytkownika.md) | [EN](en/04-key-system-workflows/4.3-user-registration-saga.md)
- [4.4. Przepływ 4: Potok jakości CI/CD, audyt i bramka regresji](pl/04-kluczowe-przeplywy-pracy/4.4-potok-jakosci-cicd-i-bramka-regresji.md) | [EN](en/04-key-system-workflows/4.4-cicd-quality-pipeline.md)

### 5. Dokumentacja katalogów i plików / Modules & Files
- [5.1. Katalog główny monorepo, benchmarks, examples, scripts](pl/05-struktura-katalogow-i-modulow/5.1-katalog-glowny-benchmarks-examples-scripts.md) | [EN](en/05-directory-and-file-structure/5.1-root-monorepo-and-tooling.md)
- [5.2. Moduł .dist/integra-core (dystrybucja Open-Core)](pl/05-struktura-katalogow-i-modulow/5.2-modul-dist-integra-core.md) | [EN](en/05-directory-and-file-structure/5.2-open-core-distribution.md)
- [5.3. Moduł integra-proto i integra-postgres-adapters](pl/05-struktura-katalogow-i-modulow/5.3-moduly-integra-proto-i-postgres-adapters.md) | [EN](en/05-directory-and-file-structure/5.3-proto-and-postgres-adapters.md)
- [5.4. Moduł integra-shared-go (biblioteka Enterprise)](pl/05-struktura-katalogow-i-modulow/5.4-modul-integra-shared-go.md) | [EN](en/05-directory-and-file-structure/5.4-shared-go-library.md)
- [5.5. Mikroserwis integra-iam-service (tożsamość i licencje)](pl/05-struktura-katalogow-i-modulow/5.5-mikroserwis-integra-iam-service.md) | [EN](en/05-directory-and-file-structure/5.5-iam-service.md)
- [5.6. Mikroserwis integra-user-service (profile i kompetencje)](pl/05-struktura-katalogow-i-modulow/5.6-mikroserwis-integra-user-service.md) | [EN](en/05-directory-and-file-structure/5.6-user-service.md)
- [5.7. Mikroserwis integra-ledger-service (księgowość i faktury)](pl/05-struktura-katalogow-i-modulow/5.7-mikroserwis-integra-ledger-service.md) | [EN](en/05-directory-and-file-structure/5.7-ledger-service.md)
- [5.8. Moduł integra-infra (chmura, Kubernetes, Helm, Terraform)](pl/05-struktura-katalogow-i-modulow/5.8-modul-integra-infra.md) | [EN](en/05-directory-and-file-structure/5.8-cloud-infrastructure.md)


---

## 🛠️ Działy operacyjne i inżynierskie
- [**Przegląd architektury (Architecture Overview)**](architecture/overview.md)
- [**Konwencje nazewnictwa pakietów (Package Naming Guidelines)**](architecture/package-naming.md)
- [**Przewodnik operacyjny i wdrożeniowy (Operations & Deployment Guide)**](operations/deployment.md)
- [**Przykłady kodu (Examples Directory)**](../examples/README.md)
- [**Katalog raportów audytowych i benchmarkowych (Audit Reports)**](reports/README.md)
