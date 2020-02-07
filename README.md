<p align="center">
  <img src="https://hsto.org/webt/ih/ds/fu/ihdsfuqni5apj0my18tnukzztw0.png" alt="Logo" width="128" />
</p>

# Тестовое задание для GoLang-разработчика


## Описание задания

Приложение, предоставляет HTTP API для получения данных о парковках такси в г. Москва. Данные берутся с [этой страницы][dataset_link] ("Актуальная версия").

## Запуск приложения
  ```
  $ make build
  $ docker-compose up -d
  ```
  По умолчанию приложение запускается на 8080 порту. После запуска можно перейти по ссылке http://localhost:8080/

## Параметры запуска приложения
- `-redis` - Cтрока подключения к БД. По умолчанию `redis:6379`
- `-db` Номер БД Redis
- `-password` Пароль от БД, По умолчанию нет пароля
- `-port` Порт на котором запускается WEB приложение
- `-source` Источник данных. По умолчанию [ссылка][dataset_link],так же возможно указать файл ( в репозитории лежит файл с данными `/data/`).

## Endpoins
 -`/api/load` - прозводит загрузку данных с источника

 -`/api/taxiparkings/{id:[0-9]+}` - Возвращает данные о парковке по её Global ID

  -`/api/taxiparkings/localid/{id:[0-9]+}` - Возвращает данные о парковках по ID ( поддерживается `?limit={num}&offset={num}`)

  -`/api/taxiparkings/mode/{mode}` - Возвращает данные о парковках по времени работ. Например `24-hours`. ( поддерживается `?limit={num}&offset={num}`)

## Метрики PROMETHEUS

- Метрики доступны по URL http://localhost:8080/metrics
- requests_counters - счетчик запросов к API EndPoint
- request_processing_time_summary_ms - длительность запросов к API EndPoint
- response_status  - счетчик запросов по статусам

[dataset_link]:https://data.gov.ru/opendata/7704786030-taxiparking
