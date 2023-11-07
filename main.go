package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"math"
	"runtime"
	"time"
)

type conf struct {
	CpuPercent    float64 `json:"cpuPercent"`
	MemoryPercent int     `yaml:"memoryPercent"`
	Off           bool    `json:"off"`
	wasStopped    bool
	context       context.Context
	cancel        context.CancelFunc
}

func main() {
	//инициализируем mem для чтения количетства текущей аллоцированной памяти
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("Before get(): Alloc = %v MiB\n", mem.Alloc/1024/1024)
	// fiber используется для быстроты написания
	app := fiber.New() // create a new Fiber instance
	b := conf{}
	resOfRes := make([][]int, 0, 100)
	//создали контекст
	b.context, b.cancel = context.WithCancel(context.Background())

	app.Post("/", func(c *fiber.Ctx) error {
		//прокинули струткуру с контекстом
		handler(c, &b)
		return nil
	})

	app.Post("/memory", func(c *fiber.Ctx) error {
		//прокинули струткуру с контекстом
		c.BodyParser(&b)
		handlerMemory(resOfRes, mem, &b)
		return nil
	})
	app.Listen(":3000")

}

func handlerMemory(resOfRes [][]int, mem runtime.MemStats, b *conf) {
	for i := 0; i < 100; i++ {
		res := getLastElem(b)

		runtime.GC() // запускаем для очистки ранее выделенной памяти
		runtime.ReadMemStats(&mem)
		fmt.Printf("After getAll(): Alloc = %v MiB, slice: %v\n", mem.Alloc/1024/1024, res)

		resOfRes = append(resOfRes, res)
	}
}

func getLastElem(b *conf) []int {
	sl := make([]int, 0, b.MemoryPercent)

	for i := 0; i < b.MemoryPercent; i++ {
		sl = append(sl, i)
	}

	// return last element

	return sl[(b.MemoryPercent - 1):]
}

func handler(c *fiber.Ctx, b *conf) {

	c.BodyParser(&b)
	log.Println(c.Body())
	log.Println(b)
	if (b.CpuPercent > 100 || b.CpuPercent < 1) && !b.Off {
		log.Println(fmt.Errorf("неправильный процент cpu"))
		return
	}
	//если контекст не стопали , то стопаем
	if !b.wasStopped {
		b.cancel()
	}

	if b.Off {
		b.Off = false
		b.wasStopped = true
		return
	}
	fmt.Println("---------------------------------------------------------------------------------------------")
	time.Sleep(5 * time.Second)
	b.context, b.cancel = context.WithCancel(context.TODO())
	b.wasStopped = false
	go RunCPULoad(b.context, 1, b.CpuPercent)
}

func RunCPULoad(ctx context.Context, coresCount int, percentage float64) {
	runtime.GOMAXPROCS(coresCount)
	// 1 unit = 100 ms may be the best
	go func() {
		runtime.LockOSThread()
		// endless loop
		for {
			select {
			case <-ctx.Done():
				log.Printf("завершил контекст для %f", percentage)
				runtime.UnlockOSThread()
				return
			default:
				load := percentage / 100.0 * 0.2 // вычисляем нагрузку в 200 миллиядрах * 0.2
				// тут может быть for
				start := time.Now() // запоминаем время начала выполнения цикла
				for i := 0; i < 1000000; i++ {
					math.Sqrt(float64(i)) // выполняем некоторую операцию, чтобы нагрузить процессор
				}
				elapsed := time.Since(start)                                // вычисляем время выполнения цикла
				sleepTime := time.Duration(float64(elapsed)/load) - elapsed // вычисляем время, которое нужно "спать", чтобы держать заданную нагрузку
				if sleepTime > 0 {
					time.Sleep(sleepTime) // "спим", чтобы держать заданную нагрузку
				} else {
					fmt.Println("Load is too high!") // если заданная нагрузка слишком высока, выводим сообщение об ошибке
					break
				}
				select {
				case <-ctx.Done():
					log.Printf("завершил контекст для %f", percentage)
					runtime.UnlockOSThread()
					return
				default:
					continue
				}
			}
		}
	}()
	return
}
