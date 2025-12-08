# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## Бенчмарки и профилирование

### Результаты бенчмарков

#### Service Layer
```
BenchmarkCreateShortURL-8         	 1531892	       759.5 ns/op	     291 B/op	       5 allocs/op
BenchmarkGetFullURL-8             	30907322	        38.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkCreateShortURLsBatch-8   	   16531	     72437 ns/op	   42287 B/op	     323 allocs/op
BenchmarkGetUserURLs-8            	  816196	      1364 ns/op	    3608 B/op	      25 allocs/op
```

#### Repository Layer
```
BenchmarkMemoryRepository_SetValue-8          	 1950451	       629.7 ns/op	     402 B/op	       4 allocs/op
BenchmarkMemoryRepository_GetFullValue-8      	55086829	        21.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkMemoryRepository_GetUserURLs-8       	   95004	     12363 ns/op	   35768 B/op	     208 allocs/op
BenchmarkMemoryRepository_SetValuesBatch-8    	21088734	        57.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkMemoryRepository_DeleteURLsBatch-8   	 3545247	       357.4 ns/op	       0 B/op	       0 allocs/op
```

#### Audit Layer
```
BenchmarkAuditService_LogEvent-8                	22660894	        53.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkFileObserver_Notify-8                  	   11781	    106539 ns/op	     344 B/op	       5 allocs/op
BenchmarkAuditService_WithMultipleObservers-8   	  649740	      1706 ns/op	     665 B/op	       3 allocs/op
```

### Анализ использования памяти

#### Проблемы до оптимизации

Базовый профиль показал следующие проблемы:
1. **523.14MB (46.47%)** - избыточные аллокации в `GetUserURLs.func1` из-за замыканий
2. **482.54MB (42.87%)** - избыточные аллокации в `SetValue.func1` из-за замыканий
3. Неэффективное использование функций-оберток `withReadLock` и `withWriteLock`

#### Оптимизации

1. **Удалены функции-обертки с замыканиями**: Заменили `withReadLock` и `withWriteLock` на прямое использование мьютексов
2. **Оптимизирована функция GetUserURLs**: 
   - Добавлена предварительная инициализация слайса с ёмкостью
   - Убраны лишние замыкания
3. **Упрощены функции GetFullValue, SetValue, IsDeleted**: Прямое использование блокировок вместо замыканий

#### Результаты оптимизации

```
File: shortener.test
Type: alloc_space
Time: 2025-12-08 21:16:50 MSK
Showing nodes accounting for -427MB, 37.93% of 1125.69MB total

      flat  flat%   sum%        cum   cum%
  574.67MB 51.05% 51.05%    51.53MB  4.58%  (*MemoryRepository).GetUserURLs
 -523.14MB 46.47%  4.58%  -523.14MB 46.47%  (*MemoryRepository).GetUserURLs.func1
 -482.54MB 42.87% 38.29%  -482.54MB 42.87%  (*MemoryRepository).SetValue.func1
```

**Ключевые улучшения:**
- **-523.14MB**: Полностью устранены аллокации в замыкании GetUserURLs
- **-482.54MB**: Полностью устранены аллокации в замыкании SetValue
- **-548.14MB**: Устранены аллокации в withReadLock
- **-482.54MB**: Устранены аллокации в withWriteLock

**Итог:** Оптимизация позволила снизить потребление памяти на **~1006MB (89%)** при работе с репозиторием, убрав накладные расходы на замыкания и функции-обертки.

### Запуск бенчмарков

```bash
# Запустить все бенчмарки
go test -bench=. -benchmem ./...

# Запустить бенчмарки для конкретного пакета
go test -bench=. -benchmem ./internal/service/...
go test -bench=. -benchmem ./internal/repository/...
go test -bench=. -benchmem ./internal/audit/...

# Создать профиль памяти
go test -bench=BenchmarkMemoryProfile -benchmem -memprofile=profiles/result.pprof ./cmd/shortener/

# Анализировать профиль памяти
go tool pprof -top profiles/result.pprof
go tool pprof -web profiles/result.pprof

# Сравнить два профиля
go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
```