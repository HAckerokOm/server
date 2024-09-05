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
$sql = "SELECT  id, root, size, elapsed_time, date_created FROM Test";

// Получение всех строк результата
$stmt = $pdo->prepare($sql);
$stmt->execute();

$rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

// Закрытие соединения
$pdo = null;

// Подготовка данных для графика
$sizes = array_column($rows, 'size');
$times = array_column($rows, 'elapsed_time');


echo "<h2>Данные из таблицы:</h2>";
echo "<table border='1'>";
echo "<tr><th>ID</th><th>Путь</th><th>Размер файла</th><th>Время запроса</th><th>Дата</th></tr>";

foreach ($rows as $row) {
    echo "<tr>";
    foreach ($row as $value) {
        echo "<td>" . htmlspecialchars($value) . "</td>";
    }
    echo "</tr>";
}
echo "</table>";


// Создание массива и сортировка
$dataPoints = [];
for ($i = 0; $i < count($sizes); $i++) {
    $dataPoints[] = [
        'x' => $sizes[$i],
        'y' => $times[$i]
    ];
}
usort($dataPoints, function($a, $b) {
    return $a['x'] - $b['x'];
  });
$dataPoints = json_encode($dataPoints);



echo "<!DOCTYPE html>
<html lang='ru'>
<head>
    <meta charset='UTF-8'>
    <title>График зависимости размера от времени запроса</title>
    <script src='https://cdn.jsdelivr.net/npm/chart.js'></script>
</head>
<body>

<h2>График зависимости размера от времени запроса</h2>
<canvas id='myChart' width='400' height='200'></canvas>

<script>
    const ctx = document.getElementById('myChart').getContext('2d');
    const dataPoints = $dataPoints;

    const myChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Размер', 'Время запроса'],
            datasets: [{
                label: 'Зависимость',
                data: dataPoints,
                borderColor: 'rgba(75, 192, 192, 1)',
                borderWidth: 1,
                fill: false
            }]
        },
        options: {
            scales: {
                y: {
                    type:'logarithmic',
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Время (секунды)'
                    }
                },
                x: {
                    type:'logarithmic',
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Размер (байты)'
                    }
                }
            }
        }
    });
</script>

</body>
</html>";
?>