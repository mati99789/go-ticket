# Projekt: "GoTicket" - High-Scale Event Ticketing System

## Cel
Stworzenie systemu rezerwacji biletów na masowe wydarzenia (koncerty, mecze), który jest w stanie obsłużyć tysiące zapytań na sekundę ("The Taylor Swift Problem").

To idealny poligon dla Go, ponieważ wymaga:
1.  **Ekstremalnej wydajności** (High Concurrency).
2.  **Spójności danych** (Distributed Locking, Database Transactions).
3.  **Architektury Mikroserwisowej** (Separacja domen).

## Architektura (High Level)
System będzie ewoluował od Modularnego Monolitu do pełnych Mikroserwisów:
1.  **Gateway Service (BFF)**: Public REST API.
2.  **Auth Service**: JWT, Identity.
3.  **Catalog/Events Service**: Heavy Read operations.
4.  **Booking Service**: Core logic, transactional processing.
5.  **Payment Service**: Integracja płatności.
6.  **Notification Service**: Async processing (RabbitMQ).

## Decyzje Projektowe
1.  **AWS**: Free Tier confirmed. Użyjemy OpenTofu/Terraform.
2.  **Frontend**: Brak. Testujemy przez CLI/Postman/Integration Tests. Skupienie 100% Backend & DevOps.

*(Szczegółowy stack technologiczny znajduje się w pliku `tech_stack.md`)*
