# Приложение для обработки изображений на GO
Описание:
Это веб-приложение, написанное на языке Go с использованием фреймворка Gin. Оно позволяет пользователям загружать изображения, применять к ним различные фильтры (например, черно-белый, сепия, размытие, резкость) и скачивать отфильтрованные изображения.

Основные функции:
- Загрузка изображений на сервер
- Применение фильтров к загруженным изображениям (граyscale, sepia, blur, sharpen)
- Скачивание отфильтрованных изображений
- Просмотр списка загруженных изображений
- Удаление загруженных изображений

Используемые технологии:
- Go
- Gin Web Framework
- Библиотека обработки изображений "github.com/disintegration/imaging"
- Docker

Дополнительно:
- CI/CD
- Dockerfile

# PhotoEditorAPI

## Building the Docker Image

To build the Docker image, run the following command in the project directory:

```sh
docker build -t photoeditorapi .
```

## Running the Docker Container

To run the Docker container, use the following command:

```sh
docker run -p 8081:8081 -v $(pwd)/uploads:/app/uploads photoeditorapi
```

This will map the `uploads` directory on your host machine to the `uploads` directory in the container, ensuring that the photos are accessible.
