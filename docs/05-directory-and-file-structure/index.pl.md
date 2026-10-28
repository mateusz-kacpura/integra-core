# 5. Szczegółowa dokumentacja katalogów i plików projektu

[◀ Poprzedni: Kluczowe przepływy pracy](../04-key-system-workflows/index.md) | [▲ Główny spis treści](../index.md)

---

## 📑 Podrozdziały w tym dziale

- [**5.1. Katalog główny (Root Monorepo), benchmarks/, examples/ oraz scripts/**](5.1-root-monorepo-and-tooling.md)
- [**5.2. Moduł .dist/integra-core (wyeksportowana dystrybucja Open-Core)**](5.2-open-core-distribution.md)
- [**5.3. Moduł integra-proto i integra-postgres-adapters**](5.3-proto-and-postgres-adapters.md)
- [**5.4. Moduł integra-shared-go (biblioteka współdzielona Enterprise)**](5.4-shared-go-library.md)
- [**5.5. Mikroserwis integra-iam-service (zarządzanie tożsamością, rolami i licencjami)**](5.5-iam-service.md)
- [**5.6. Mikroserwis integra-user-service (profile użytkowników, konta i kompetencje)**](5.6-user-service.md)
- [**5.7. Mikroserwis integra-ledger-service (księga główna podwójnego zapisu, Event Sourcing, faktury i wypłaty)**](5.7-ledger-service.md)
- [**5.8. Moduł integra-infra (infrastruktura chmurowa, Kubernetes, Helm, Terraform, WAF i monitoring)**](5.8-cloud-infrastructure.md)

---

Monorepo obejmuje łącznie **808 plików** zorganizowanych w ściśle odseparowanych modułach. Poniżej znajduje się kompletny opis każdego katalogu i pliku w podziale na komponenty.

---

### 5.1. Katalog główny (Root Monorepo), `benchmarks/`, `examples/` oraz `scripts/`

| Ścieżka pliku / katalogu | Przeznaczenie i szczegółowy opis techniczny |
| --- | --- |
| `go.work`, `go.work.sum` | Definicja przestrzeni roboczej Go 1.26 łącząca 7 modułów lokalnych (`examples`, `integra-iam-service`, `integra-ledger-service`, `integra-postgres-adapters`, `integra-proto`, `integra-shared-go`, `integra-user-service`) wraz z dyrektywami `replace` dla prefiksu `integra-enterprise`.
| `Makefile` | Główny plik automatyzacji zadań deweloperskich: `fmt`, `imports` (`goimports`), `tidy`, `tidy-check`, `vet`, `lint`, `staticcheck`, `test` (z flagami `-race -shuffle=on`), `test-short`, `audit`.
| `LICENSE` | Licencja komercyjna (*Proprietary and Confidential*) dla części Enterprise platformy Integra.
| `.github/workflows/ci.yaml` | Definicja potoku GitHub Actions uruchamiająca dwa zadania (`modules` oraz `infrastructure`) z kontenerami serwisowymi `postgres:18-alpine` (3 instancje: IAM, User, Ledger), `mongo:8` i `valkey/valkey:9`.
| `.gitignore`, `.dockerignore`, `.gitattributes` | Reguły wykluczania artefaktów kompilacji, profili `.pprof`, sekretów `.env`, kluczy `.pem`/`.key` oraz wymuszanie końców linii LF dla plików `.go`, `.sh`, `.sql`, `.yaml`, `.proto` i CRLF dla `.ps1`.
| `benchmarks/BenchmarkCharts.py` | Silnik wizualizacji napisany w czystym Pythonie (zero zewnętrznych zależności), generujący wykresy wektorowe SVG (`latency-tier-pyramid.svg`, `speedup-top-optimizations.svg`, `historical-trend.svg`), bloki Mermaid oraz interaktywny dashboard HTML.
| `benchmarks/BENCHMARKS-DASHBOARD.html` | Samowystarczalny, interaktywny dashboard wydajnościowy w stylu *dark-mode glassmorphism* wykorzystujący Chart.js do prezentacji metryk L1-L3, alokacji pamięci i trendów historycznych.
| `examples/go.mod`, `examples/go.sum` | Moduł Go dla wykonywalnych przykładów demonstracyjnych.
| `examples/jwt-auth/main.go` | Wykonywalny przykład demonstracyjny pokazujący generowanie kluczy asymetrycznych Ed25519, wystawianie tokena JWT i jego weryfikację z użyciem pamięci podręcznej L1 CPU (500 ms TTL).
| `examples/sqlrouting-l1/main.go` | Wykonywalny przykład demonstrujący dynamiczne przepisywanie nazw tabel SQL (`outbox_events` -> `iam.outbox`), odporność na ataki SQL Injection oraz pomiar czasu kompilacji w pamięci L1.
| `scripts/Audit-Project.py` | Główna bramka jakościowa (Quality Gate) uruchamiająca walidację wszystkich modułów lokalnych, `staticcheck` oraz testy z detektorem wyścigów (`-race`).
| `scripts/Check-LocalModules.py` | Narzędzie orkiestrujące wykonywanie poleceń `go` w każdym podmodule monorepo z wyłączonym `GOWORK=off`, symulujące zachowanie niezależnych pakietów przed publikacją.
| `scripts/Collect-Benchmarks.py` | Kolektor statystyk benchmarkowych parsujący wyjście `go test -bench`, obliczający mediany z wielokrotnych przebiegów, zapisujący historię do `BENCHMARK-HISTORY.csv` i blokujący CI w przypadku regresji przekraczającej zadany próg procentowy.
| `scripts/Prepare-OpenCoreSplit.py` | Skrypt ekstrakcyjny budujący niezależne repozytorium Open-Core (`.dist/integra-core`) na licencji Apache 2.0 poprzez wyodrębnienie wskazanych pakietów z `integra-shared-go` i `integra-proto`, przepisanie ścieżek importów i weryfikację testów.
| `scripts/Migrate-ToEnterprisePrefix.py` | Narzędzie automatyzujące migrację ścieżek modułów do przestrzeni nazw `[github.com/mateusz-kacpura/integra-enterprise/](https://github.com/mateusz-kacpura/integra-enterprise/)...`.
| `scripts/Export-HandwrittenFiles.ps1`, `Export-ProjectStructure.ps1`, `Simulate-CI.ps1` | Skrypty PowerShell do eksportu ręcznie pisanych plików źródłowych, generowania drzewa projektu (`project-structure.txt`) oraz lokalnej symulacji potoku CI na systemie Windows.

---

### 5.2. Moduł `.dist/integra-core` (wyeksportowana dystrybucja Open-Core)

Moduł ten stanowi publiczny rdzeń otwartoźródłowy (`module [github.com/mateusz-kacpura/integra-core](https://github.com/mateusz-kacpura/integra-core)`) wygenerowany przez `Prepare-OpenCoreSplit.py`:

* **Konfiguracja i kontrakty (`pkg/config`, `pkg/baseconfig`, `pkg/contracts`)**:
* `pkg/config/config.go`, `config_test.go`, `doc.go`: Bazowa struktura `BaseConfig` ładowana przez `viper` (z prefiksem środowiskowym `INTEGRA_`), walidacja portów HTTP i wymaganych zmiennych (`SERVICE_NAME`, `DATABASE_URL`) oraz stała `DefaultShutdownGrace = 2s`.


* `pkg/baseconfig/config.go`: Re-eksport typów konfiguracyjnych zapobiegający kolizjom nazw pakietów w mikroserwisach.


* `pkg/contracts/read_replicas.go`, `architecture_test.go`: Współdzielone struktury DTO dla replik odczytu (`UserRead`, `OrgRead` z polem `Version int64`) oraz testy AST pilnujące czystości architektonicznej.




* **Szybki JSON i kodowanie (`pkg/jsoncodec`, `pkg/json`)**:
* `pkg/jsoncodec/jsoncodec.go`: Wrapper na ultraszybką bibliotekę JIT/SIMD `[github.com/bytedance/sonic](https://github.com/bytedance/sonic)` (`Marshal`, `Unmarshal`, `Pretouch`, `NewEncoder`, `NewDecoder`).


* `pkg/json/json.go`, `json_test.go`: Warstwa kompatybilności z testami deterministycznej serializacji znaków specjalnych Unicode i strumieni.




* **Logowanie z sanityzacją PII (`pkg/logger`)**:
* `logger.go`, `logger_test.go`: Inicjalizacja `zerolog` (`InitEntrypoint`, `NewBootstrapLogger`) z automatycznym wstrzykiwaniem metadanych operacyjnych (`service`, `component`, `process_type`, `environment`, `version`, `correlation_id`).


* `sanitizer.go`, `sanitizer_test.go`: Wyrażenia regularne i strumieniowy `SanitizingWriter` maskujący karty płatnicze, PESEL, NIP, tokeny JWT oraz adresy e-mail.


* `slog.go`, `slog_test.go`: Zaawansowane handlery dla standardowego `log/slog` (`MultiHandler`, `LevelFilterHandler`, `TraceAndSanitizingHandler`, `NewOperationalAndAuditHandler` rozdzielający logi `INFO/WARN` na `stdout` od błędów `ERROR` na `stderr`).


* `adapter.go`, `adapter_test.go`: Adapter `PortAdapter` łączący interfejsy portów aplikacyjnych z `zerolog` i telemetrią.




* **Paginacja, odpowiedzi HTTP, telemetria, tematy Kafka i UUIDv7**:
* `pkg/pagination/cursor.go`, `grpc.go`, `cursor_bench_test.go`, `cursor_fuzz_test.go`, `cursor_test.go`, `example_test.go`: Jednoalokacyjne kodowanie i dekodowanie kursorów (`CursorEncode`, `CursorDecode` na buforach stosowych `[64]byte` / `[128]byte`), obsługa `CursorWhere` oraz paginacja gRPC.


* `pkg/response/wrapper.go`, `error_codes.go`, `mapper.go`, `validation.go` (+ testy): Ustandaryzowana koperta odpowiedzi REST (`Response{Success, Code, Message, Data, Meta}`), mapowanie 12 kanonicznych kodów błędów na statusy HTTP oraz walidacja pól.


* `pkg/telemetry/propagation.go`, `trace_hook.go`, `tracer.go` (+ testy): Propagacja kontekstu W3C (`traceparent`, `tracestate`) i `X-Correlation-ID` przez middleware HTTP, klienta `RoundTripper` oraz tabelę transakcyjną Outbox (`OutboxHeaders`).


* `pkg/topics/topics.go`, `example_test.go`: Definicje wszystkich **38 kanonicznych tematów Kafka** w ekosystemie Integra oraz funkcja `TopicForAggregate`.


* `pkg/uuidgen/uuidgen.go`, `pkg/uuidx/uuidx.go`: Generatory czasowo-sortowalnych identyfikatorów **UUIDv7 (RFC 9562)** oraz bezalokacyjne konwertery do `pgtype.UUID` (`ToPgtype`, `FromPgtype`).


* `pkg/health/health.go`, `health_test.go`: Endpointy sond Kubernetes (`LivenessHandler`, `ReadinessHandler`) wykonujące kontrole zależności bez blokowania muteksu.


* `proto/`: Kontrakty Protobuf (`accounting/v1/invoice.proto`, `accounting/v1/ledger.proto`, `iam/v1/user.proto`, `user/v1/profile.proto`, `buf.yaml`).





---

### 5.3. Moduł `integra-proto` i `integra-postgres-adapters`

#### `integra-proto/`

* `buf.yaml`, `go.mod`, `go.sum`: Konfiguracja generatora i lintera **Buf v2** z regułami `STANDARD` i detekcją zmian łamiących kompatybilność (`FILE`).


* `iam/v1/user.proto`: Kontrakt `IdentityService` (`ValidateToken`, `GetUser`, `GetUserBatch` do 100 ID, `GetUserRoles`, `GetLicense`, `CheckPermission`).


* `user/v1/profile.proto`: Kontrakt `ProfileService` (`GetProfile`).


* `accounting/v1/ledger.proto`: Kontrakt `LedgerService` (`GetTrueBalance`, `RecordTransaction` z polami `expected_debit_version` i `expected_credit_version`).


* `accounting/v1/invoice.proto`: Kontrakt `InvoiceService` (`GenerateInvoice`, `InvoiceResponse`, `InvoiceLineItem` z kwotami przesyłanymi jako precyzyjne łańcuchy dziesiętne).



#### `integra-postgres-adapters/`

* `db/schema/adapters.sql` oraz `db/queries/`:
* `authority.sql`: Zapytania SQL dla autorytetu ról i licencji.


* `outbox.sql`: Kanoniczne zapytania obsługujące tabelę `outbox_events`.


* `replicas.sql`: Zapytania dla lokalnych replik odczytu (`users_read`, `orgs_read`).


* `state.sql`: Zapytania o stan bezpieczeństwa kont.




* `sqlrouting/table.go`, `table_test.go`, `integration_test.go`, `example_test.go`, `doc.go`: Adapter `sqlrouting.Table` przechwytujący zapytania `sqlc` i podmieniający nazwę tabeli szablonowej na kwalifikowaną tabelę docelową (np. `iam.outbox`, `ledger.users_read`) przy użyciu pamięci podręcznej L1 (`19.11 ns/op`) oraz rygorystycznej walidacji białej listy schematów (`iam`, `ledger`, `users`, `authz`, `audit`, `public`) i identyfikatorów.


* `fixtures.go`, `sqlc.yaml`, `doc.go`: Konfiguracja generatora `sqlc` i dane pomocnicze do testów adaptera.



---

### 5.4. Moduł `integra-shared-go` (biblioteka współdzielona Enterprise)

Oprócz pakietów eksportowanych do `integra-core`, moduł `integra-shared-go` zawiera kluczowe pakiety infrastrukturalne i bezpieczeństwa klasy Enterprise:

* **`pkg/aead/` oraz `pkg/crypto/**`:
* `chacha.go`, `chacha_bench_test.go`, `chacha_fuzz_test.go`, `chacha_test.go`, `aeadtest/keys.go`: Implementacja szyfrowania uwierzytelnionego **XChaCha20-Poly1305** (klucze 256-bitowe, losowe nonce 192-bitowe) wykorzystywana do szyfrowania zawartości pamięci podręcznej w Valkey oraz wrażliwych tokenów.




* **`pkg/audit/`**:
* `security_event.go`, `security_event_test.go`, `doc.go`: Struktury niezmiennych zdarzeń audytu bezpieczeństwa (*immutable security audit events*) publikowanych na temat `security-audit-events`.




* **`pkg/auth/`**:
* `verifier.go`, `verifier_bench_test.go`, `verifier_concurrency_test.go`, `verifier_test.go`: Weryfikator `JWTVerifier` obsługujący klucze publiczne **Ed25519**, polityki wydawcy/odbiorcy (`WithPolicy`), szyfrowany cache L2 w Valkey (`ClaimsCache`) oraz pamięć podręczną L1.


* `l1_cache.go`, `l1_cache_test.go`: 32-shardowa pamięć podręczna L1 w procesie Go z haszowaniem inline FNV-1a, rotacją co 30s, limitem 32 768 wpisów, TTL 500 ms i separacją linii cache `cpu.CacheLinePad`.


* `account_state.go`, `account_recovery.go`, `blacklist.go`, `local.go`, `claims.go` (+ testy integracyjne i fuzzingowe): Atomowe skrypty Lua w Valkey weryfikujące 1-sekundowe dzierżawy stanu konta (`CheckAccountState`, `PublishAccountSnapshot`), czarne listy unieważnionych tokenów JTI oraz odzyskiwanie spójności po awarii Valkey.


* `ginmiddleware/middleware.go`, `middleware_test.go`: Middleware dla frameworka Gin integrujące weryfikację JWT i stan konta z kontekstem żądania HTTP.




* **`pkg/securitystate/`**:
* `authority.go`, `state.go`, `license_guard.go`, `consumer_linux.go`, `consumer_other.go` (+ testy klastrowe): Zarządzanie scentralizowanym autorytetem ról (`RoleAuthority`) i strażnikiem licencji (`LicenseGuard`) w klastrze Valkey z użyciem bezpiecznych tagów haszujących `{authz:roles:ready}` gwarantujących umieszczenie powiązanych kluczy w tym samym slocie klastra Valkey.




* **`pkg/rbac/` oraz `pkg/policy/**`:
* `matcher.go`, `permission.go`, `rbac.go`, `ginmiddleware/middleware.go` (+ testy i fuzzing): Bezalokacyjny silnik autoryzacji RBAC (`atomic.Pointer` dla odczytów bez blokad) obsługujący hierarchiczne uprawnienia trójsegmentowe (`zasób:podzasób:akcja`) z wieloznacznikami `*`.




* **`pkg/database/`, `pkg/pgnumeric/`, `pkg/numericx/`, `pkg/replicas/**`:
* `database/postgres.go`, `migrate.go`, `embedded_migrate.go`: Zarządzanie pulą połączeń `pgxpool.Pool` oraz wbudowanymi migracjami SQL (`embed.FS`).


* `pgnumeric/pgnumeric.go`, `numericx/numericx.go`: Precyzyjne konwersje pomiędzy typami `pgtype.Numeric` a wartościami dziesiętnymi w Go.


* `replicas/upsert.go`: Idempotentny mechanizm aktualizacji lokalnych replik odczytu (`users_read`, `orgs_read`) chroniący przed nadpisaniem nowszego stanu starszym zdarzeniem (`WHERE version < EXCLUDED.version`).




* **`pkg/kafka/` oraz `pkg/outbox/**`:
* `kafka/api.go`, `producer.go`, `consumer.go`, `batch.go`, `metadata.go`, `middleware.go`, `config.go`, `unsupported.go`: Wysokowydajny klient Kafka z obsługą wsadową, automatycznym wstrzykiwaniem metadanych OpenTelemetry i warunkową kompilacją na platformy Linux/inne.


* `outbox/worker.go`, `relay.go`, `tx.go`: Implementacja wzorca **Transactional Outbox** gwarantującego dostarczenie zdarzeń *at-least-once* z tabeli bazodanowej do brokera Kafka w ramach transakcji SQL.




* **`pkg/idempotency/`, `pkg/ratelimit/`, `pkg/lifecycle/`, `pkg/fxproviders/`, `pkg/fxutil/`, `pkg/telemetry/**`:
* Middleware idempotentności żądań HTTP (`pkg/idempotency`), ograniczania liczby żądań (`pkg/ratelimit`), zarządzania cyklem życia grup gorutyn i sygnałów OS (`pkg/lifecycle`), integracji wstrzykiwania zależności `uber/fx` (`pkg/fxproviders`, `pkg/fxutil`) oraz pełnej instrumentacji `OpenTelemetry` dla Gin i ConnectRPC (`pkg/telemetry/gin.go`, `connect.go`, `otel.go`, `runtime.go`).





---

### 5.5. Mikroserwis `integra-iam-service` (zarządzanie tożsamością, uwierzytelnianiem, rolami i licencjami)

Serwis odpowiedzialny za rejestrację, logowanie, obsługę MFA (TOTP), sesje i resetowanie haseł (MongoDB), autorytet ról RBAC oraz licencjonowanie produktów (PostgreSQL + Valkey):

* **Punkty wejścia (`cmd/`)**:
* `cmd/server/main.go`: Główny serwer HTTP (port 8081) i gRPC/ConnectRPC (port 9081).


* `cmd/worker/main.go`: Proces roboczy publikujący zdarzenia z tabeli `iam.outbox` do Kafka, czyszczący stare wpisy (`outbox_purge.go`) i uzgadniający stan bezpieczeństwa w Valkey.


* `cmd/migrate/main.go`: Narzędzie CLI wykonujące migracje schematu `iam` w PostgreSQL.


* `cmd/admin/main.go`, `main_test.go`: Narzędzie CLI do inicjalnego utworzenia konta administratora (*bootstrap admin*).




* **Baza danych i migracje (`db/queries/`, `migrations/`, `sqlc.yaml`)**:
* Zapytania `sqlc`: `iam.sql`, `identity_state.sql`, `monitoring.sql`, `outbox.sql`, `queries.sql`, `registration.sql`, `role_authority.sql`.


* 11 par migracji SQL (`000_uuidv7` do `010_drop_legacy_dead_tables`) definiujących funkcję `uuidv7()`, tabele użytkowników, ról, licencji, wersjonowanie replik organizacji oraz śledzenie W3C w tabeli `outbox`.




* **Domena (`internal/domain/`)**:
* `user.go`, `auth.go`, `credentials.go`, `license.go`, `permissions.go`, `events.go` (+ testy i benchmarki): Czyste modele domenowe użytkownika, poświadczeń Argon2id/Bcrypt, sesji uwierzytelniania, konfiguracji MFA, przypisań licencji i ról.




* **Aplikacja (`internal/app/`)**:
* `auth.go`, `users.go`, `administration.go`, `access_service.go`, `bootstrap_admin.go`, `queries.go`: Fasady serwisów aplikacyjnych i reguły autoryzacji dostępu.


* `commands/` (`auth.go`, `users.go`, `administration.go`, `bootstrap_admin.go`): Handlery komend mutujących stan (rejestracja, logowanie, weryfikacja MFA, zmiana hasła, nadawanie ról i licencji).


* `queries/handlers.go`: Handlery zapytań odczytujących użytkowników, role i licencje.


* `ports/` (`administration.go`, `audit.go`, `logging.go`, `queries.go`, `repository.go`, `security.go`, `use_case_errors.go`): Kontrakty portów wymaganych przez warstwę aplikacji wraz z benchmarkiem `queries_bench_test.go` testującym *Data-Oriented Design*.




* **Infrastruktura (`internal/infrastructure/`)**:
* `postgres/` (`postgres_unit_of_work.go`, `user_postgres.go`, `mfa_postgres.go`, `rbac_postgres.go`, `outbox_postgres.go`, `role_authority.go`, `identity_state.go`, `postgres_types.go`, `queries.go`): Implementacja wzorca **Unit of Work** spinającego repozytoria `Users`, `Mfa`, `Rbac` i `Outbox` we wspólnej transakcji `pgx.Tx` oraz publikującego aktualizacje autorytetu ról do Valkey.


* `mongodb/` (`mongo_auth.go`, `mongo_password_reset.go` + testy i benchmarki): Repozytorium MongoDB przechowujące sesje odświeżania tokenów, konfigurację MFA oraz tokeny resetowania haseł.


* `security/` (`access_tokens.go`, `access_verifier.go`, `keys.go`, `tokens.go`, `verifier.go`, `audit.go`, `runtime.go`): Generowanie i podpisywanie tokenów JWT kluczem prywatnym Ed25519 (tylko serwis IAM w całym ekosystemie posiada klucz prywatny!).


* `messaging/` (`outbox_linux.go`, `outbox_other.go`, `outbox_purge.go`, `topics.go`): Publikacja zdarzeń z tabeli `iam.outbox` i automatyczna retencja przetworzonych wpisów.


* `valkey/valkey.go`, `logging/logger.go`, `metrics.go`: Adaptery klienta Valkey, logowania i metryk Prometheus.




* **Transport (`internal/transport/`)**:
* `rest/` (`router.go`, `auth_handlers.go`, `user_handlers.go`, `security_admin.go`, `access_middleware.go`, `output.go`, `errors.go`, `docs.go`, `swagger_types.go`): Kompletne REST API z dokumentacją Swagger/OpenAPI i szybką serializacją odpowiedzi.


* `grpc/` (`server.go`, `security.go`): Serwer gRPC/ConnectRPC obsługujący metody `IdentityService`.


* `kafka/` (`accounts_linux.go`, `accounts_other.go`): Konsument zdarzeń kont.




* **Bootstrap i testy (`internal/bootstrap/`, `tests/`)**:
* `bootstrap/server.go`, `worker_linux.go`, `lifecycle.go`, `audit_linux.go`: Składanie zależności i zarządzanie cyklem życia procesów.


* `tests/e2e/api_test.go` oraz `tests/integration/` (`kafka_test.go`, `registration_test.go`, `repository_test.go`): Testy End-to-End oraz integracyjne.





---

### 5.6. Mikroserwis `integra-user-service` (profile użytkowników, konta i kompetencje)

Serwis zarządzający profilami użytkowników, umiejętnościami (*skills*) oraz lokalną projekcją stanu bezpieczeństwa kont synchronizowaną z IAM przez Kafka i Valkey:

* **Punkty wejścia (`cmd/`)**:
* `cmd/server/main.go`: Serwer HTTP (port 8082) i gRPC (port 9082).


* `cmd/worker/main.go`: Worker nasłuchujący zdarzeń z Kafka (`user-events`, `user-roles-sync`) oraz publikujący zdarzenia `user-profiles-sync` i `user-skills-sync`.


* `cmd/migrate/main.go`: Migracje bazy `users`.


* `cmd/import-profiles/main.go`: Narzędzie CLI do masowego importu profili użytkowników.




* **Baza danych (`db/queries/`, `migrations/`)**:
* Zapytania: `accounts.sql`, `monitoring.sql`, `outbox.sql`, `users.sql`.


* Migracje (`001_profiles` do `006_account_security_versions`): Tabele profili z optymistycznym blokowaniem (`version`), umiejętności, stanów autoryzacyjnych kont i śledzenia `outbox`.




* **Domena i aplikacja (`internal/domain/`, `internal/app/`)**:
* `domain/`: `account.go`, `profile.go`, `skill.go`, `permissions.go`, `events.go`.


* `app/`: `profiles.go`, `accounts.go`, `identity.go`, `access_service.go`, podkatalogi `commands/` (`accounts.go`, `identity.go`, `profiles.go`), `queries/` (`accounts.go`, `profiles.go`) oraz `ports/` (`access_ports.go`, `audit.go`, `profile.go`).




* **Infrastruktura i transport (`internal/infrastructure/`, `internal/transport/`)**:
* `barrier_state.go`, `barrier_bench_test.go`: Struktura `SecurityBarrierState` z izolacją linii cache L1D (`cpu.CacheLinePad`) dla atomowych liczników odnawiania dzierżawy i wersji źródłowych.


* `account_state.go`, `accounts.go`, `postgres.go`, `token_verifier.go`, `tokens_permissions.go`, `audit.go`, `metrics.go`, `outbox_linux.go`: Repozytoria PostgreSQL, synchronizacja barier bezpieczeństwa i weryfikacja tokenów kluczem publicznym Ed25519.


* `transport/rest/` (`router.go`, `accounts.go`, `dto.go`), `transport/grpc/server.go`, `transport/kafka/` (`consumer_linux.go`, `identity_event.go`): Warstwa komunikacyjna REST, gRPC i konsument zdarzeń tożsamości z Kafka.





---

### 5.7. Mikroserwis `integra-ledger-service` (księga główna podwójnego zapisu, Event Sourcing, faktury i wypłaty)

Najbardziej zaawansowany obliczeniowo mikroserwis w ekosystemie, realizujący operacje księgowe w architekturze **Event Sourcing + CQRS** z pełną izolacją wielodzierżawczą (*Multi-Tenancy*):

* **Punkty wejścia (`cmd/`)**:
* `cmd/server/main.go`: Serwer HTTP (port 8085) i gRPC (port 9085).


* `cmd/worker/main.go`: Worker przetwarzający zadania migawek (`snapshot_jobs`), publikujący zdarzenia z `ledger.outbox` oraz konsumujący zdarzenia replik z Kafka.


* `cmd/migrate/main.go`: Migracje schematu `ledger`.


* `cmd/project-catchup/main.go`: Narzędzie CLI do przebudowy i doganiania projekcji odczytu (*projection catch-up*) ze strumienia zdarzeń `event_store`.




* **Baza danych, migracje i skrypty SQL (`db/queries/`, `migrations/`, `scripts/`)**:
* 11 plików zapytań `sqlc`: `accounts.sql`, `event_store.sql`, `invoices.sql`, `legacy_aliases.sql`, `monitoring.sql`, `outbox.sql`, `payouts.sql`, `posting_commit.sql` (atomowe CTE księgowania), `projection_catchup.sql`, `projections.sql`, `repository.sql`.


* 16 par migracji SQL (`001_init_event_store` do `016_drop_dead_read_tables`): Definicje `event_store`, modeli odczytu, replik `users_read`/`orgs_read`, kont tenantów, aliasów historycznych, śledzenia pochodzenia projekcji (*projection provenance*), skrótów żądań `request_hash_v2` oraz surowych ładunków `outbox_raw_payload`.


* `scripts/audit-legacy-history.sql`, `reconcile-legacy-aliases.sql`: Skrypty audytowe i uzgadniające historię księgową.




* **Domena (`internal/domain/`)**:
* `aggregate.go`: Główny agregat `LedgerAggregate` wykorzystujący alokację arenową (`accountArena`, `invoiceArena`, `walletArena`, `journalArena`) do bezalokacyjnego odtwarzania stanu ze strumienia zdarzeń.


* `money.go`, `decimal.go`, `currency.go`: Niezmienny typ `Money` oparty na 64-bitowych jednostkach mniejszych (`minor int64`) oraz szybki parser łańcuchów dziesiętnych `ParseScaledDecimal` (`35.77 ns/op`, `0 B/op`).


* `account.go`, `invoice.go`, `payout.go`, `payout_state.go`, `snapshot.go`, `events.go`, `permissions.go`, `errors.go`: Logika domenowa kont księgowych, kalkulacji faktur VAT (netto, podatek, brutto), maszyny stanów wypłat (*payout state machine*) oraz migawek kont.




* **Aplikacja, infrastruktura i transport (`internal/app/`, `internal/infrastructure/`, `internal/transport/`)**:
* `app/`: `operations.go`, `access_service.go`, `commands/handlers.go`, `commands/integration.go`, `queries/handlers.go`, `ports/` (`access_ports.go`, `payout.go`, `posting.go`, `publication.go`, `repository.go`, `tenant.go`).


* `infrastructure/tigerbeetle/`:
* `client.go`, `client_test.go`: Inicjalizacja binarnego klienta Go (`github.com/tigerbeetle/tigerbeetle-go`) z adresem klastra, ID klastra i obsługą Direct I/O.
* `id.go`, `id_test.go`: Bezalokacyjna konwersja RFC 4122 UUID $\leftrightarrow$ 128-bitowe liczby całkowite TigerBeetle (`UUIDToUint128`, `Uint128ToUUID`).
* `repository.go`, `repository_test.go`: Implementacja portów `AccountRepository` i `PostingRepository` z bezpośrednią obsługą `client.CreateTransfers` i `client.LookupAccounts`, atomową weryfikacją debetów/kredytów i eliminacją blokad SQL `SELECT ... FOR UPDATE`.


* `infrastructure/postgres/`:
* `request_hash.go`: Binarne ramkowanie żądań na buforze stosowym `[512]byte` i wyliczanie 33-bajtowego skrótu idempotentności SHA-256 (`RequestDigest`, `200.9 ns/op`, `0 B/op`).
* `event_store.go`, `event_codec.go`: Zapis i odczyt zdarzeń audytowych z pulą buforów serializacyjnych.
* `invoice_store.go`, `invoice_codec.go`, `payout_store.go`, `projections.go`, `outbox.go`, `tenant.go`, `audit.go`, `metrics.go`: Asynchroniczny magazyn metadanych biznesowych, faktur, wypłat i relacji wielodzierżawczych.




* `transport/`: Serwery `grpc/server.go`, `rest/router.go`, `rest/secure_router.go`, `rest/handlers.go`, `kafka/router.go` oraz `security/http.go` i `routes.go`.





---

### 5.8. Moduł `integra-infra` (infrastruktura chmurowa, Kubernetes, Helm, Terraform, WAF i monitoring)

Kompletna definicja **Infrastructure as Code (IaC)** oraz narzędzi DevSecOps/SRE:

* **Warstwa brzegowa i WAF (`deploy/edge/`)**:
* `bunkerweb/bunkerweb.conf` oraz `modsecurity/custom-rules.conf`: Konfiguracja firewalla aplikacyjnego **BunkerWeb + ModSecurity** z obsługą `api.integra.com`, `admin.integra.com`, `iot.integra.com`, ochroną AntiBot, limitem `100r/s`, czarnymi listami IP (`ipsum`) oraz automatycznym TLS Let's Encrypt i HTTP/2.


* `gateways/` (`public-gateway.yaml`, `admin-gateway.yaml`, `admin-routes.yaml`, `iot-gateway.yaml`) oraz `docker-compose.edge.yaml`: Odseparowane bramy wejściowe dla ruchu publicznego, administracyjnego i IoT.




* **Helm Chart i Kubernetes (`deploy/helm/integra/`, `deploy/k8s/`)**:
* `Chart.yaml`, `values.yaml`, `templates/` (`microservice.yaml`, `linkerd.yaml`, `network-policy.yaml`, `configmap.yaml`, `secrets.yaml`, `_helpers.tpl`): Uniwersalny szablon wdrożeniowy Kubernetes obsługujący adnotacje siatki usług **Linkerd mTLS**, polityki sieciowe Zero-Trust (`NetworkPolicy`) oraz autoskalowanie poziome **HPA** (2–10 replik przy 70% CPU).


* `deploy/k8s/infra/`: Manifesty dla klastrów `tigerbeetle.yaml` (StatefulSet TigerBeetle z Direct I/O, `IPC_LOCK` i portem 3000), `postgresql.yaml`, `mongodb.yaml`, `valkey.yaml`, `kafka.yaml`, `jaeger.yaml`, `loki.yaml`, `monitoring.yaml`.




* **Terraform AWS (`deploy/terraform/`)**:
* `main.tf`, `variables.tf`, `outputs.tf`, `environments/` (`dev.tfvars`, `staging.tfvars`, `prod.tfvars`) oraz moduły w `modules/`:
* `vpc/main.tf`: Sieć VPC z podsieciami prywatnymi i publicznymi.


* `eks/main.tf`: Zarządzany klaster Kubernetes **Amazon EKS**.


* `rds/main.tf`: Bazy relacyjne **Amazon RDS (PostgreSQL)**.


* `documentdb/main.tf`: Klaster dokumentowy **Amazon DocumentDB (kompatybilny z MongoDB)**.


* `elasticache/main.tf`: Klaster **Amazon ElastiCache (Valkey)**.


* `msk/main.tf`: Zarządzany klaster **Amazon MSK (Apache Kafka)**.






* **Monitoring i obserwowalność (`monitoring/`)**:
* `prometheus.yaml`, `alertmanager.yaml`, `alert_rules.yaml`: Zbieranie metryk i reguły alertowania.


* `loki-config.yaml`, `fluentd.conf`, `jaeger.yaml`: Centralizacja logów i śladów rozproszonych.


* `grafana/dashboards/` (`integra-overview.json`, `iam-service.json`, `user-service.json`) oraz `monitoring/tests/` (`iam.yaml`, `ledger.yaml`, `user.yaml`): Gotowe tablice rozdzielcze Grafana oraz testy reguł monitoringu.




* **Skrypty automatyzacji i testów chaosu oraz bezpieczeństwa (`scripts/`)**:
* 26 skryptów operacyjnych i weryfikacyjnych: m.in. `platform.py`, `validate.sh`, `check-sqlc.py`, `check-proto.py`, `check-yaml.py`, `gen-jwt-keys.sh`/`.ps1`, `profile-ledger.sh`, `test-security.py`, `test-security-recovery.py`, `test-ledger-quarantine.py`, `test-jaeger-persistence.py`.

---

[◀ Poprzedni: Kluczowe przepływy pracy](../04-key-system-workflows/index.md) | [▲ Główny spis treści](../index.md)
