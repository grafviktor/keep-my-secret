# Keep My Secret #

Готовое решение для хранения личных данных пользователей в зашифрованном виде. Поддерживаются следующие типы данных:

* Кредитная карта
* Двоичные данные (файлы, картинки, документы и прочее)
* Заметки в виде произвольного текста
* Данные входа в виде пароль / логин

Приложение состоит из 2 частей - серверной и клиентской части, которые общаются друг с другом по защищенному протоколу `HTTPS`.

## Запуск приложения ##



## Сервер ##

Сервер принимает и обрабатывает REST-запросы от клиентов, а также взаимодействует с хранилищем данных. Сервер работает только с использованием протокола HTTPS.

### Авторизация ###

Если пользователь предоставляет верную комбинацию логина и пароля, в ответ сервер отправляет пару JWT-токенов:

* `Access token` на основе них принимается решение о предоставлении доступа к закрытым ресурсам. Этот токен имеет малый срок жизни и содержит дату окончания действия.

* `Refresh tokens`, удостоверяют запрос на получение нового Access Token'а. Например, когда истек срок действия Access Token, клиент может отправить новый запрос на получение токена на сервер авторизации. Чтобы такой запрос был успешно выполнен, он сопровождается Refresh токеном. РRefresh Token'ы имеет долгий период действия и недоступен для браузерного JavaScript.

### Приватность данных ###


### API ###

Поддерживаются два основных типа REST-запросов:

1. запросы отвечающие за аутентификацию пользователей;
2. запросы отвечающие за сохранение и извелечение данных пользователя.

Данные ответов и запросов передаются в формате JSON. Бинарные файлы передаются в части запроса с типом `multipart/form-data`.

#### Аутентификация пользователей ####

| URL                        | HTTP Method | Параметры          | Описание                        |
|----------------------------|-------------|--------------------|---------------------------------|
| /api/v1/user/register      | POST        | username, password | регистрация нового пользователя |
| /api/v1/user/login         | POST        | username, password | авторизация пользователя        |
| /api/v1/user/logout        | POST        | -                  | завершение сессии               |
| /api/v1/user/token-refresh | GET         | -                  | обновление токена доступа       |

#### Сохранение и получение объектов данных пользователя ####

| URL                        | HTTP Method | Параметры               | Описание                        |
|----------------------------|-------------|-------------------------|---------------------------------|
| /api/v1/secrets/           | GET         | -                       | получение всех сохраненных объектов пользователя |
| /api/v1/secrets/           | POST        | см.пример1              | сохранения нового объекта       |
| /api/v1/secrets/           | PUT         | см.пример2              | обновление (замена) существующего объекта|
| /api/v1/secrets/{id}       | DELETE      | -                       | удаление объекта                |
| /api/v1/secrets/file/{id}  | GET         | -                       | получение бинарного файла       |

пример1
```json
{
  "type":"card",
  "title":"Bank card",
  "cardholder_name":"Mr. Tony Tester",
  "card_number":"1234 5678 9012 3456",
  "expiration":"2023-09-03",
  "security_code":"999",
  "note":"Карта, где деньги лежат"
}
```

пример2
```json
{
  "id": 33,
  "type":"card",
  "title":"Bank card",
  "cardholder_name":"Mr. Tony Tester",
  "card_number":"1234 5678 9012 3456",
  "expiration":"2023-09-03",
  "security_code":"999",
  "note":"Карта, где деньги лежат"
}
```

#### Дополнительно ####

| URL              | HTTP Method | Параметры          | Описание                        |
|------------------|-------------|--------------------|---------------------------------|
| /api/v1/version  | GET         | -                  | регистрация нового пользователя |

#### Ответы сервера ####

## Клиент ##

Клиенсткая часть представляет собой SPA-приложение, которое несет в себе минимум логики (Dumb Client). Браузерный клиент легко может замененым любым другим приложением, которое может работать поверх `HTTPS`-протокола.

## Test Coverage Report ##

[![codecov](https://codecov.io/gh/grafviktor/keep-my-secret/branch/master/graph/badge.svg?token=wrIL0tyQ5q)](https://codecov.io/gh/grafviktor/keep-my-secret)

## Имеющиеся проблемы ##

* данные между сервером и клиентом передаются в несжатой форме;
* все криптографические операции с данными производятся на стороне сервера, что ведет к проблемам масштабируемости
* ... и другие неисчислимые возможности по улучшению