# Сервис создан для нагрузки на cpu или выделения памяти по требованию для автотестирвоания приложения задействующего метрики

Все взаимодейтсвие происходит по endpoint'ам / и /memory

Первый необходим для загрузки cpu , второй для выделеиния памяти 


Примеры запросов 

#Выделяем память , каждый memory percent это ~ 2350 
```
curl -X POST  localhost:3000/memory  -H "Content-Type: application/json" —data '{"memoryPercent":112500}'
```

```
curl -X POST  localhost:3000 —header "Content-Type: application/json" —data '{"cpuPercent":1}'

```

Для того , чтобы мы могли могли обращаться к нему из пода приложения в хельм чартах создается сервис 

Service Account необходим для подвязывания токена 