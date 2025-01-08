Приложение, позволяющее просматривать доступные отели, выбирать понравившиеся номера и бронировать их на назначенный период. Имеется возможность регистрации в системе как в роли обычного клиента, так и в роли владельца отеля. Есть имитация оплаты через собственную систему с ответом через вебхук. Сервис бронирования посылает сообщение о заказе в кафку, сервис нотификации его обрабатывает. Общение между приложением и клиентом - REST, общение между микросервисами - gRPC. Работа велась в команде. Все сервисы подняты в docker контейнерах с использованием docker-compose. Обеспечивается наблюдаемость сервисов с использованием jaeger, prometheus и логгера zap. Клиенты и отельеры доолжны пройти аутентификацию и авторизацию.
Стэк: Go, postgres, kafka, docker, docker-compose, REST API, gRPC, jaeger, prometheus, git

Основные микросервисы:

1) Booking Svc (Сервис бронирования):
Основной сервис, управляющий процессом бронирования.
Отвечает за создание бронирований.
Взаимодействует с базой данных Booking Data для хранения информации о бронированиях.
Отправляет события в Queue (Kafka) для последующей обработки в Notification Svc.
2) Hotel Svc (Сервис отеля):
Управляет информацией об отелях, включая данные о номерах и ценах.
Отвечает за добавление и обновление отелей.
Хранит данные в базе Hotels Data.
Предоставляет информацию Booking Svc для обработки запросов на бронирование.
3) Notification Svc (Сервис уведомлений):
Отвечает за отправку уведомлений клиентам и отельерам о статусе бронирования.
Получает события из очереди Queue (Kafka) и отправляет уведомления через Delivery System.
4) Auth Svc (Сервис авторизации):
Отвечает за регистрацию клиентов и отельеров в системе.
Выдает JWT токены.
Хранит информацию о пользователях в БД.
5) Payment System (Система оплаты):
Отвечает за обработку платежей.
Интегрирована с Booking Svc для проведения транзакций по бронированию отелей.
Имитация оплаты.
Предоставляет API для создания заявки на оплату, возвращает результат через webhook.
