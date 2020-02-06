<p align="center">
  <img src="https://hsto.org/webt/ih/ds/fu/ihdsfuqni5apj0my18tnukzztw0.png" alt="Logo" width="128" />
</p>

# Тестовое задание для GoLang-разработчика


## Описание задания

Приложение, предоставляет HTTP API для получения данных о парковках такси в г. Москва. Данные берутся с [этой страницы][dataset_link] ("Актуальная версия").

## Метрики PROMETHEUS

- Метрики доступны по URL http://localhost:8080/metrics
- requests_counters - счетчик запросов к API EndPoint
- request_processing_time_summary_ms - длительность запросов к API EndPoint
- response_status  - счетчик запросов по статусам

[dataset_link]:https://data.gov.ru/opendata/7704786030-taxiparking
