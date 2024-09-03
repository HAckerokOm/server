<?php

    // Берем данные аунтификации из соответсвющего файла
    include 'Auth.php';

    // Проверяем, был ли отправлен POST-запрос
    if ($_SERVER['REQUEST_METHOD'] === 'POST') {

        //Получаем данные из POST-запроса
        header('Content-Type: application/json'); 
        $jsonData = file_get_contents('php://input');
        $data = json_decode($jsonData, true); 

        $root = $data['f_root'];
        $size = $data['f_size'];
        $elapsedTime = $data['f_requesttime'];

        // Подключаемся к базе данных
        $conn = new mysqli($host, $user , $password, $dbname);

        // Проверяем подключение
        if ($conn->connect_error) {
            die(json_encode(['status' => 'error', 'message' => 'Ошибка подключения: ' . $conn->connect_error]));
        }

        // Подготавливаем запрос INSERT
        $sql = "INSERT INTO Test (root, size, elapsed_time, date_created) VALUES (?, ?, ?, ?)";
        $stmt = $conn->prepare($sql);

        // Связываем параметры с запросом
        $currentDateTime = date('Y-m-d H:i:s');
        $stmt->bind_param("sids", $root, $size, $elapsedTime, $currentDateTime);

        // Выполняем запрос
        if ($stmt->execute()) {
            echo json_encode(['status' => 'success', 'message' => 'Данные успешно добавлены в базу данных!']);
        } else {
            echo json_encode(['status' => 'error', 'message' => 'Ошибка при добавлении данных: ' . $stmt->error]);
        }

        // Закрываем подготовленный запрос
        $stmt->close();

        // Закрываем подключение
        $conn->close();
    } 
    else {
        echo json_encode(['status' => 'error', 'message' => 'Неверный метод запроса']);
    }