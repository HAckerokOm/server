let curDir = "/home/yaroslav/work"; // Определяем текущую директорию
let parDir; // Переменная для сохранения родительской директории

document.addEventListener('DOMContentLoaded', function(){ // Добавляем слушатель события загрузки DOM
    Elget() // Вызываем функцию Elget()
    updateCurrentPath(); // Вызываем функцию updateCurrentPath() для обновления текущего пути
    fetchFiles('').then(files => { // Загружаем файлы без сортировки
        displayFiles(files.f_data);
    }).catch(error => { //обработка ошибки 
        console.error('Ошибка при загрузке данных:', error);
    });
    BackBut() // Вызываем функцию BackBut() для добавления обработчика кнопки назад
});
//fetchFiles - Функция для получения файлов
function fetchFiles(sort) {
    return fetch(`http://localhost:9999/fs?dst=${curDir}&sort=${sort}`) // Создаем URL для запроса к API
        .then(response => { // Проверяем, является ли ответ успешным (код состояния 200-299)
            if (!response.ok) { 
                // Если ответ не успешный, выбрасываем ошибку с информацией о статусе
                throw new Error(`HTTP error! status: ${response.status}`); 
            }
            return response.json();  // Если ответ успешный, преобразуем его в JSON
        });
}
// displayFiles - Функция отображения файлов в таблице
function displayFiles(files) {
    const tableBody = document.querySelector('.item-list tbody');// Находим тело таблицы с классом .item-list tbody
    tableBody.innerHTML = ''; // Очищаем содержимое тела таблицы
    
    files.forEach(file => {// проходимся по списку файлов
        const row = document.createElement('tr'); // Создаем новую строку таблицы
        row.classList.add(file.f_type === 'directory'? 'itemdir' : 'itemfile');// Добавляем класс itemdir, если это директория, иначе itemfile

        ['f_type', 'f_name', 'f_size','f_typesize'].forEach(key => {  // Цикл для создания ячеек в строке
            const cell = document.createElement('td');// Создаем новый элемент ячейки
            cell.innerHTML = `<span class="${key}">${file[key]}</span>`; // Формируем содержимое ячейки с применением стилей
            row.appendChild(cell);  // Добавляем ячейку в строку
        });

        tableBody.appendChild(row);// Добавляем строку в тело таблицы

        if (file.f_type == "Директория") { //если это директория
             // Добавляем эффект наведения и обработчик клика для директорий
            row.classList.add("itemdir")
            row.addEventListener('click', function() {
                directoryTraversal(file.f_name)// Обработчик клика для перехода в поддиректорию
            
           
            });
        }
    });
}
// Elget - Функция для обработки сортировки файлов
function Elget(){
    document.getElementById("3").addEventListener('click', function() { // Добавляем обработчик клика на элемент с id="3"
        const sort = prompt('Введите параметр сортировки (asc или desc):', '');  // Показываем пользователю окно ввода для сортировки
        if (sort) { // Проверяем, не пустая ли строка ввода
            fetchFiles(sort).then(files => { // Вызываем функцию fetchFiles с переданным параметром сортировки
                displayFiles(files.f_data);// Отображаем отсортированные файлы в таблице
            }).catch(error => { //обрабатываем ошибки
                console.error('Ошибка при получении данных:', error);
            });
        }
    })}
// BackBut - функция кнопки назад
function BackBut(){
    let backButton = document.getElementById("2") // Получаем элемент с id="2" и добавляем к нему обработчик клика
        backButton.addEventListener('click', returnToPreviousDirectory) // При нажатии кнопки вызывается функция returnToPreviousDirectory
}
// directoryTraversal - Функция для перехода в поддиректорию
function directoryTraversal(path){
    parDir = curDir; // Сохраняем текущий путь в переменную parDir
    curDir = `${curDir}/${path}`; // Обновляем текущий путь, добавив указанный путь
    updateCurrentPath(); // Обновляем текущий путь в DOM
    fetchFiles('').then(files => { // Загружаем файлы без сортировки
        displayFiles(files.f_data);  // Отображаем полученные файлы в таблице
    }).catch(error => { //обработка ошибки
        console.error('Ошибка при загрузке данных:', error);
    });
}
function returnToPreviousDirectory() {
    if(curDir === "/home") { // Проверяем, равен ли текущий путь корневой директории
        alert("Вы в корневой дирректории")  // Если да, показываем сообщение и прекращаем выполнение функции
        return
    }
    curDir = curDir.split('/').slice(0, -1).join('/'); // Разделяем текущий путь на компоненты
    updateCurrentPath(); // Обновляем текущий путь в DOM
    fetchFiles('').then(files => { // Загружаем файлы из новой директории
        displayFiles(files.f_data); // Отображаем полученные файлы в таблице
    }).catch(error => { //обработка ошибки
        console.error('Ошибка при загрузке данных:', error);
    });
}
// updateCurrentPath - Функция для обновления текущего пути в DOM
function updateCurrentPath() {
    // Получаем элемент с id="4" и устанавливаем его текстовое содержимое
    document.getElementById("4").textContent = `Путь: ${curDir}`;
}
