# Pełna dokumentacja architektoniczna, techniczna i biznesowa platformy **Integra Enterprise & Open-Core (Go 1.26 Monorepo)**

> **Wersja:** Go 1.26 Monorepo  
> **Licencja hybrydowa:** Apache 2.0 (Open-Core) / Proprietary Commercial (Enterprise)  
> **Przeglądarka schematów pełnoekranowych:** [`../DIAGRAMS-VIEWER.html`](../DIAGRAMS-VIEWER.html)  
> **Jednolity plik zbiorczy:** [`../doc-pl.md`](../doc-pl.md)

---

## 📚 Spis treści dokumentacji modułowej

### [1. Cel projektu, wartość biznesowa i unikalność na rynku Go](01-cel-i-wartosc-biznesowa/README.md)

- [1.1. Cel projektu i wizja architektoniczna](01-cel-i-wartosc-biznesowa/1.1-cel-projektu-i-wizja-architektoniczna.md)
- [1.2. Wartość biznesowa (Business Value)](01-cel-i-wartosc-biznesowa/1.2-wartosc-biznesowa.md)
- [1.3. Unikalność w skali rynku projektów w Go](01-cel-i-wartosc-biznesowa/1.3-unikalnosc-na-rynku-go.md)

### [2. Analiza statystyczna benchmarków i uzasadnienie wydajnościowe](02-analiza-benchmarkow-i-wydajnosc/README.md)

- [2.1. Statystyki opisowe całej platformy](02-analiza-benchmarkow-i-wydajnosc/2.1-statystyki-opisowe-platformy.md)
- [2.2. Statystyki wydajności w podziale na moduły monorepo](02-analiza-benchmarkow-i-wydajnosc/2.2-statystyki-wydajnosci-modulow.md)
- [2.3. Dowody benchmarkowe na unikalność architektoniczną](02-analiza-benchmarkow-i-wydajnosc/2.3-dowody-benchmarkowe-wzorce-i-anomalie.md)
- [2.4. Wizualizacje wyników benchmarków (Mermaid XY Charts)](02-analiza-benchmarkow-i-wydajnosc/2.4-wizualizacje-wykresy-mermaid.md)

### [3. Globalna architektura systemu i wizualizacje Mermaid](03-architektura-globalna/README.md)

- [3.1. Architektura ekosystemu i infrastruktury](03-architektura-globalna/3.1-ekosystem-i-infrastruktura.md)
- [3.2. Architektura zależności modułów Go Workspace (go.work) i podział Open-Core](03-architektura-globalna/3.2-zaleznosci-modulow-go-work-i-open-core.md)
- [3.3. Wewnętrzna architektura czysta / heksagonalna + CQRS mikroserwisów](03-architektura-globalna/3.3-architektura-czysta-i-cqrs-mikroserwisow.md)

### [4. Kluczowe przepływy pracy (System Workflows)](04-kluczowe-przeplywy-pracy/README.md)

- [4.1. Przepływ 1: Wielopoziomowa weryfikacja tokena JWT i stanu bezpieczeństwa (L1 -> L2 -> L3)](04-kluczowe-przeplywy-pracy/4.1-wielopoziomowa-weryfikacja-tokena-jwt.md)
- [4.2. Przepływ 2: Atomowe księgowanie podwójnego zapisu (Double-Entry Ledger Posting & Idempotency)](04-kluczowe-przeplywy-pracy/4.2-atomowe-ksiegowanie-podwojnego-zapisu.md)
- [4.3. Przepływ 3: Rozproszona saga rejestracji użytkownika i synchronizacji replik odczytu](04-kluczowe-przeplywy-pracy/4.3-rozproszona-saga-rejestracji-uzytkownika.md)
- [4.4. Przepływ 4: Potok jakości CI/CD, audyt bezpieczeństwa i bramka regresji benchmarków](04-kluczowe-przeplywy-pracy/4.4-potok-jakosci-cicd-i-bramka-regresji.md)

### [5. Szczegółowa dokumentacja katalogów i plików projektu](05-struktura-katalogow-i-modulow/README.md)

- [5.1. Katalog główny (Root Monorepo), benchmarks/, examples/ oraz scripts/](05-struktura-katalogow-i-modulow/5.1-katalog-glowny-benchmarks-examples-scripts.md)
- [5.2. Moduł .dist/integra-core (wyeksportowana dystrybucja Open-Core)](05-struktura-katalogow-i-modulow/5.2-modul-dist-integra-core.md)
- [5.3. Moduł integra-proto i integra-postgres-adapters](05-struktura-katalogow-i-modulow/5.3-moduly-integra-proto-i-postgres-adapters.md)
- [5.4. Moduł integra-shared-go (biblioteka współdzielona Enterprise)](05-struktura-katalogow-i-modulow/5.4-modul-integra-shared-go.md)
- [5.5. Mikroserwis integra-iam-service (zarządzanie tożsamością, rolami i licencjami)](05-struktura-katalogow-i-modulow/5.5-mikroserwis-integra-iam-service.md)
- [5.6. Mikroserwis integra-user-service (profile użytkowników, konta i kompetencje)](05-struktura-katalogow-i-modulow/5.6-mikroserwis-integra-user-service.md)
- [5.7. Mikroserwis integra-ledger-service (księga główna podwójnego zapisu, Event Sourcing, faktury i wypłaty)](05-struktura-katalogow-i-modulow/5.7-mikroserwis-integra-ledger-service.md)
- [5.8. Moduł integra-infra (infrastruktura chmurowa, Kubernetes, Helm, Terraform, WAF i monitoring)](05-struktura-katalogow-i-modulow/5.8-modul-integra-infra.md)


---

> [!TIP]
> **Interaktywny podgląd na pełnym ekranie**: Dostępna jest dedykowana przeglądarka [DIAGRAMS-VIEWER.html](../DIAGRAMS-VIEWER.html) z obsługą trybu pełnoekranowego (przycisk *Pełny ekran* lub klawisz `F11`), płynnego przybliżania kółkiem myszy (zoom) oraz przesuwania (drag & pan) wszystkich 9 schematów architektonicznych i przepływów.
