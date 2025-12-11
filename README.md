# Система управления бронированиями

Простое API на языке Go с использованием фреймворка Fiber для управления бронированием комнат. Поддерживает регистрацию пользователей, создание комнат администраторами, бронирование, отмену бронирования и просмотр информации.

## Зависимости

- [Go](https://golang.org/dl/)
- [Fiber v2](https://docs.gofiber.io/)
- [Logger middleware](https://docs.gofiber.io/api/middleware/logger)

Для установки зависимостей выполните:

```bash
go mod init <your_module_name>
go get github.com/gofiber/fiber/v2
go get github.com/gofiber/fiber/v2/middleware/logger

Сервер будет запущен на http://localhost:3000.

API Эндпоинты
Регистрация пользователя
Регистрирует нового пользователя в системе.

URL: /register
Метод: GET
Параметры:
fam (string) — Фамилия
ima (string) — Имя
otch (string) — Отчество
phone (string) — Номер телефона
login (string) — Логин (уникальный)
kolvo (int) — Количество персон

Пример:
GET /register?fam=Иванов&ima=Иван&otch=Иванович&phone=123456&login=user1&kolvo=4

Создание администратора
Создаёт нового администратора. Требует уникального логина.

URL: /admin/create
Метод: GET
Параметры:
login (string) — Логин администратора (уникальный)
Пример:
GET /admin/create?login=admin1

Создание комнаты
Создаёт новую комнату. Доступно только администраторам.

URL: /room/create
Метод: GET
Параметры:
login (string) — Логин администратора
num (int) — Номер комнаты
cap (int) — Вместимость комнаты
Пример:
GET /room/create?login=admin1&num=101&cap=10

Удаление пользователя
Удаляет пользователя из системы. Доступно только администраторам.

URL: /user/delete
Метод: GET
Параметры:
login (string) — Логин администратора
user_login (string) — Логин удаляемого пользователя
Пример:
GET /user/delete?login=admin1&user_login=user1

Бронирование комнаты
Позволяет пользователю забронировать свободную комнату на определённую дату. Проверяется вместимость комнаты и статус бронирования.

URL: /room/book
Метод: GET
Параметры:
login (string) — Логин пользователя
room_num (int) — Номер комнаты
date (string) — Дата бронирования (формат: YYYY-MM-DD)
Пример:
GET /room/book?login=user1&room_num=101&date=2025-12-10

Отмена бронирования
Отменяет бронирование комнаты. Доступно владельцу брони или администратору.

URL: /room/unbook
Метод: GET
Параметры:
login (string) — Логин пользователя
room_num (int) — Номер комнаты
Пример:
GET /room/unbook?login=user1&room_num=101

Удаление комнаты
Удаляет комнату из системы. Доступно только администраторам.

URL: /room/delete
Метод: GET
Параметры:
login (string) — Логин администратора
room_num (int) — Номер удаляемой комнаты
Пример:
GET /room/delete?login=admin1&room_num=101

Изменение количества персон
Позволяет пользователю изменить своё значение количества персон.

URL: /user/update_kolvo
Метод: GET
Параметры:
login (string) — Логин пользователя
kolvo (int) — Новое количество персон
Пример:
GET /user/update_kolvo?login=user1&kolvo=6

Изменение даты бронирования
Позволяет пользователю изменить дату своей брони. Доступно владельцу брони или администратору.

URL: /room/update_date
Метод: GET
Параметры:
login (string) — Логин пользователя
room_num (int) — Номер комнаты
new_date (string) — Новая дата (формат: YYYY-MM-DD)
Пример:
GET /room/update_date?login=user1&room_num=101&new_date=2025-12-15

Просмотр профиля
Показывает информацию о пользователе.

URL: /user/profile
Метод: GET
Параметры:
login (string) — Логин пользователя
Пример:
GET /user/profile?login=user1

Просмотр всех комнат
Показывает список всех комнат. Доступно зарегистрированным пользователям и администраторам.

URL: /rooms
Метод: GET
Параметры:
login (string) — Логин пользователя
Пример:
GET /rooms?login=user1

Просмотр бронирований пользователя
Показывает все бронирования конкретного пользователя.

URL: /user/bookings
Метод: GET
Параметры:
login (string) — Логин пользователя
Пример:
GET /user/bookings?login=user1

Просмотр доступных комнат
Показывает список комнат, подходящих пользователю по вместимости и статусу (свободны).

URL: /rooms/available
Метод: GET
Параметры:
login (string) — Логин пользователя
Пример:
GET /rooms/available?login=user1

##Примеры всех запросов
Регистрация пользователя: GET /register?fam=Иванов&ima=Иван&otch=Иванович&phone=123456&login=user1&kolvo=4
Создание администратора: GET /admin/create?login=admin1
Создание комнаты (только админ): GET /room/create?login=admin1&num=101&cap=10
Удаление пользователя (только админ): GET /user/delete?login=admin1&user_login=user1
Бронирование комнаты: GET /room/book?login=user1&room_num=101&date=2025-12-10 ...БРОНИРОВАНИЕ ПО ЛОГИНУ
Отмена бронирования: GET /room/unbook?login=user1&room_num=101
Удаление комнаты (только админ): GET /room/delete?login=admin1&room_num=101
Изменение количества персон: GET /user/update_kolvo?login=user1&kolvo=6
Изменение даты бронирования: GET /room/update_date?login=user1&room_num=101&new_date=2025-12-15
Просмотр профиля: GET /user/profile?login=user1
Просмотр всех комнат: GET /rooms?login=user1
Просмотр бронирований пользователя: GET /user/bookings?login=user1
Просмотр доступных комнат: GET /rooms/available?login=user1

Хранение данных
Данные сохраняются в файл data.json в корне проекта.
```
