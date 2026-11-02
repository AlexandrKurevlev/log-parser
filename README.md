# Log Parser

## Описание
Учебный проект по основам Go. Инструмент для анализа логов распределённых систем.

## Возможности
- Сопоставление записей по request_id
- Обнаружение ошибок и построение хронологии
- Конкурентная обработка
- JSON-отчёты

## Установка и запуск
```bash
go build -o log-processor main.go
./log-processor --input-dir /var/log/services/
```

## Примеры использования
```bash
go run ./main.go -input-dir "./logs" -output-file "results.json"
```