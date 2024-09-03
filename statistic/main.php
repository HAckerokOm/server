<?php
// Настройки подключения к базе данных
include 'Auth.php';
// Подключение к базе данных
try {
    $pdo = new PDO("mysql:host=$host;dbname=$dbname", $user, $password);
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
} catch (PDOException $e) {
    echo "Ошибка подключения к базе данных: " . $e->getMessage();
    exit();
}

// Формирование SQL-запроса
$sql = "SELECT * FROM Test";

// Выполнение запроса
try {
    $stmt = $pdo->query($sql);
} catch (PDOException $e) {
    echo "Ошибка выполнения запроса: " . $e->getMessage();
    exit();
}

// Получение всех строк результата
$rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

// Закрытие соединения
$pdo = null;

// Вывод данных
echo "<h2>Данные из таблицы:</h2>";
echo "<table border='1'>";
foreach ($rows as $row) {
    echo "<tr>";
    foreach ($row as $value) {
        echo "<td>" . htmlspecialchars($value) . "</td>";
    }
    echo "</tr>";
}
echo "</table>";

?>