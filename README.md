# BankSystem
## Запуск
1) Для работы приложения требуется:
- go >= v1.24
- docker >= v27
- GNU Make >= v4.3
- Свободные порты 8080 и 55432

2) Для запуска программы выполняем команду:
```
make run
```

## Тестирование
1) Все доступные методы можно загрузить в http клиент (insomnia/postman), используя файл [Insomnia.json](Insomnia.json)
2) В БД уже загружен тестовый пользователь с login `1_user@mail.com` password `123456`, с помощью этого пользователя можно протестировать все методы.
3) Для тестирования всего флоу пользотеля, необходимо:
    - `/account/create` создать аккаунт
    - `/account/deposit` пополнить депозит
    - `/card/create` создать карту и привязать ее к акаунту
    - `/card/payment` теперь можно производить оплату через карту 

## Структура API:
/auth/register → Регистрация нового пользователя
/auth/login → Авторизация через email/пароль
/user/profile → Получение данных профиля
/account/create → Создание аккаунта
/account/deposit → Пополнение баланса
/account/withdraw → Списание средств
/card/create → Создание новой карты
/card/payment → Оплата по карте
/transfer/create → Перевод между аккаунтами

## Таблица эндпоинтов API
|Метод|Путь             |Назначение                           |Группа  |Требует авторизации| Описание                                                     |Описание                              |
|-----|-----------------|-------------------------------------|--------|-------------------|--------------------------------------------------------------|------------------------------------|
|POST |/auth/register   |Регистрация нового пользователя      |auth    |❌ Не требуется     | Создает нового пользователя с email                          | username и паролем.                |
|POST |/auth/login      |Авторизация                          |auth    |❌ Не требуется     | Возвращает JWT-токен после успешной проверки учетных данных. |                                    |
|GET  |/user/profile    |Получить данные текущего пользователя|user    |✅ Да               | Возвращает информацию о пользователе из базы данных.         |                                    |
|POST |/account/create  |Создание аккаунта                    |account |✅ Да               | Создает новый банковский аккаунт для пользователя.           |                                    |
|POST |/account/deposit |Пополнение баланса аккаунта          |account |✅ Да               | Увеличивает баланс указанного аккаунта.                      |                                    |
|POST |/account/withdraw|Списание средств с аккаунта          |account |✅ Да               | Уменьшает баланс указанного аккаунта.                        |                                    |
|GET  |/account/all     |Получить все аккаунты пользователя   |account |✅ Да               | Возвращает список всех аккаунтов                             | связанных с пользователем.         |
|POST |/card/create     |Создать новую карту                  |card    |✅ Да               | Привязывает карту к аккаунту.                                |                                    |
|POST |/card/payment    |Оплата по карте                      |card    |✅ Да               | Выполняет оплату и уведомляет пользователя по email          | проверяя CVV и срок действия карты.|
|POST |/transfer/create |Перевод между аккаунтами             |transfer|✅ Да               | Переводит средства с одного аккаунта на другой.              |                                    |
