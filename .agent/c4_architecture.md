# C4 Architecture - GoTicket

Ten dokument wizualizuje architekturę systemu na poziomie kontenerów (serwisów).

## Level 1: System Context
Użytkownik kupuje bilet. System komunikuje się z Płatnościami i Emailami.

## Level 2: Container Diagram (Target State: Microservices)

```mermaid
graph TD
    User((User))
    PaymentGW[External Payment Provider]
    
    subgraph "GoTicket System"
        Gateway[API Gateway / BFF]
        Auth[Auth Service]
        Catalog[Catalog Service]
        Booking[Booking Service]
        Payment[Payment Service]
        Notification[Notification Service]
        
        DB_Core[(Postgres - Core)]
        Redis[(Redis - Cache/Locks)]
        RabbitMQ[RabbitMQ]
    end

    User -->|HTTP/REST| Gateway
    Gateway -->|gRPC| Auth
    Gateway -->|gRPC| Catalog
    Gateway -->|gRPC| Booking
    
    Booking -->|SQL| DB_Core
    Booking -->|Redis Protocol| Redis
    Catalog -->|SQL| DB_Core
    
    Booking -->|gRPC| Payment
    Payment -->|Webhook| PaymentGW
    
    Booking -.->|Publish Event| RabbitMQ
    RabbitMQ -.->|Consume Event| Notification
```

## Level 2: Container Diagram (Initial State: Modular Monolith)
Na początku (Phase 1-2) wszystko będzie w jednym binarnym pliku, ale z zachowaniem separacji pakietów (modułów).

```mermaid
graph TD
    User((User))
    
    subgraph "GoTicket App (Modular Monolith)"
        HTTP[HTTP Handler / Router]
        
        subgraph "Modules (Internal Packages)"
            Mod_Auth[Module: Auth]
            Mod_Catalog[Module: Catalog]
            Mod_Booking[Module: Booking]
        end
        
        DB[(Postgres)]
    end

    User -->|HTTP| HTTP
    HTTP --> Mod_Auth
    HTTP --> Mod_Catalog
    HTTP --> Mod_Booking
    
    Mod_Booking --> DB
    Mod_Catalog --> DB
```
