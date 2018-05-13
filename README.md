## Задание:
>Написать сервис на Golang, который принимает массив URL-ов товаров https://www.amazon.co.uk,
>для данных URL он должен загрузить наименование, цена, фото товара(только URL, грузить изображение не надо), пример:
>
>[
>   {
>    "url": "https://www.amazon.co.uk/gp/product/1509836071",
>    "meta": {
>         "title": "The Fat-Loss Plan: 100 Quick and Easy Recipes with Workouts",
>         "price": "8.49",
>         "image": "https://images-na.ssl-images-amazon.com/images/I/51kB2nKZ47L._SX382_BO1,204,203,200_.jpg",
>     }
>   },
>...
>]
>
>Сервис необходимо завернуть в Docker. В Docker можно сделать опционально.
>Загрузчик должен иметь конфигурируемое ограничение на количество одновременных запросов в amazon.

## Тестирование:
curl -X POST -d "
[\"https://www.amazon.co.uk/Lean-15-Sustain-Minute-Workouts/dp/1509820221/\",
\"https://www.amazon.co.uk/Lean-15-Minute-Workouts-Strong/dp/1509800697/\",
\"https://www.amazon.co.uk/Lean-15-Minute-Workouts-Healthy/dp/1509800662/\",
\"https://www.amazon.co.uk/gp/product/1509836071\"]"
http://localhost:8082/

в конфигурационном файле настройка порта и количество числа параллельных обработчиков запроса

[SERVER]

port = 8082

workers = 10

В лог выводится информация о том, какая ссылка каким процессом была обработана:

2018/05/13 10:44:32 https://www.amazon.co.uk/Lean-15-Minute-Workouts-Strong/dp/1509800697/ 3

2018/05/13 10:44:32 https://www.amazon.co.uk/Lean-15-Sustain-Minute-Workouts/dp/1509820221/ 9

2018/05/13 10:44:32 https://www.amazon.co.uk/gp/product/1509836071 0

2018/05/13 10:44:32 https://www.amazon.co.uk/Lean-15-Minute-Workouts-Healthy/dp/1509800662/ 4

##### Ответ сервиса:

>[
> {
>  "url": "https://www.amazon.co.uk/Lean-15-Minute-Workouts-Strong/dp/1509800697/",
>  "meta": {
>   "title": "Lean in 15 - The Shape Plan: 15 Minute Meals With Workouts to Build a Strong, Lean Body",
>   "image": "https://images-na.ssl-images-amazon.com/images/I/51duWEaS9uL._SX258_BO1,204,203,200_.jpg",
>   "price": 8.49
>  }
> },
> {
>  "url": "https://www.amazon.co.uk/Lean-15-Sustain-Minute-Workouts/dp/1509820221/",
>  "meta": {
>   "title": "Lean in 15 - The Sustain Plan: 15 Minute Meals and Workouts to Get You Lean for Life",
>   "image": "https://images-na.ssl-images-amazon.com/images/I/51qSkb3aCzL._SX258_BO1,204,203,200_.jpg",
>   "price": 7
>  }
> },
> {
>  "url": "https://www.amazon.co.uk/gp/product/1509836071",
>  "meta": {
>   "title": "The Fat-Loss Plan: 100 Quick and Easy Recipes with Workouts",
>   "image": "https://images-na.ssl-images-amazon.com/images/I/51IsTylYiPL._SX382_BO1,204,203,200_.jpg",
>   "price": 8.49
>  }
> },
> {
>  "url": "https://www.amazon.co.uk/Lean-15-Minute-Workouts-Healthy/dp/1509800662/",
>  "meta": {
>   "title": "Lean in 15 - The Shift Plan: 15 Minute Meals and Workouts to Keep You Lean and Healthy",
>   "image": "https://images-na.ssl-images-amazon.com/images/I/51dVpS3WXzL._SX258_BO1,204,203,200_.jpg",
>   "price": 8
>  }
> }
>]

