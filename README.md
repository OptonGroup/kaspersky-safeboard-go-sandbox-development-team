Worker Pool (Go)

Простой пул воркеров на Go без сторонних зависимостей.

API

Пакет: `github.com/OptonGroup/kaspersky-safeboard-go-sandbox-development-team/pool`

- type Pool interface:
  - Submit(func()) error — добавить задачу на выполнение. Возвращает ErrQueueFull, если очередь заполнена, или ErrStopped, если пул остановлен.
  - Stop() error — корректно остановить пул (идемпотентно): закрывает очередь и дожидается завершения воркеров.

- Конструктор: NewPool(workers, queueSize int, opts ...Option) (Pool, error)
  - workers >= 0, queueSize >= 0, иначе ErrInvalidConfig.
  - Опции:
    - WithOnTaskDone(func()) — хук, вызывается один раз на задачу (успех/паника).

Гарантии

- Безопасное выполнение задач: паники внутри задач перехватываются, логируются, пул продолжает работу.
- onTaskDone вызывается всегда, даже если задача паникнула; паника внутри хука игнорируется.
- Stop идемпотентен, дожидается завершения уже взятых задач; новые Submit после Stop запрещены.
- Никаких утечек горутин: воркеры завершаются после закрытия канала задач.

Пример

См. `pool/example_test.go`.

```
p, _ := pool.NewPool(4, 16, pool.WithOnTaskDone(func(){ /* метрика */ }))
defer p.Stop()

_ = p.Submit(func(){ /* работа */ })
_ = p.Submit(func(){ /* ещё работа */ })
```

Разработка

- Тесты: go test -race -count=1 ./...
- CI: GitHub Actions запускает go vet и тесты с -race, собирает покрытие для ./tests.


