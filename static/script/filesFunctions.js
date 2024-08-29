
import {mask} from "./script.js"
import {loader} from "./script.js"
import {updateCurrentPath, curDir, changeCurDir}from './buttoms.js'

let parDir; // Переменная для сохранения родительской директории
//fetchFiles - Функция для получения файлов
export function fetchFiles(sort) {
    loader.classList.add('loader')
    mask.classList.add('mask')
   
    return fetch(`http://localhost:9999/fs?dst=${curDir}&sort=${sort}`) // Создаем URL для запроса к API
        .then(response => { // Проверяем, является ли ответ успешным (код состояния 200-299)
            if (!response.ok) { 
                // Если ответ не успешный, выбрасываем ошибку с информацией о статусе
                throw new Error(`HTTP error! status: ${response.status}`); 
            }
            return response.json();  // Если ответ успешный, преобразуем его в JSON
        }).then(files => { // Вызываем функцию fetchFiles с переданным параметром сортировки
            displayFiles(files.f_data);// Отображаем отсортированные файлы в таблице
        }).catch(error => { //обрабатываем ошибки
            console.error('Ошибка при получении данных:', error);
        });
}
// displayFiles - Функция отображения файлов в таблице
export function displayFiles(files) {
    const tableBody = document.querySelector('.item-list tbody');// Находим тело таблицы с классом .item-list tbody
    tableBody.innerHTML = ''; // Очищаем содержимое тела таблицы
    
    files.forEach(file => {// проходимся по списку файлов
        const row = document.createElement('tr'); // Создаем новую строку таблицы
        row.classList.add(file.f_type === 'directory'? 'itemdir' : 'itemfile');// Добавляем класс itemdir, если это директория, иначе itemfile

        ['f_type', 'f_name'].forEach(key => {  // Цикл для создания ячеек в строке
            const cell = document.createElement('td');// Создаем новый элемент ячейки
            cell.innerHTML = `<span class="${key}">${file[key]}</span>`; // Формируем содержимое ячейки с применением стилей
            row.appendChild(cell);  // Добавляем ячейку в строку
        });
        // Создаем отдельную ячейку для размера и типа размера
        const sizeCell = document.createElement('td');
        sizeCell.className = 'size-cell';

        // Добавляем текст размера
        const sizeSpan = document.createElement('span');
        sizeSpan.textContent = `${file.f_size}`;
        sizeSpan.className = 'f_size';
        sizeCell.appendChild(sizeSpan);

        // Добавляем текст типа размера
        const typeSpan = document.createElement('span');
        typeSpan.textContent = `${file.f_typesize}`;
        typeSpan.className = 'f_typesize';
        sizeCell.appendChild(typeSpan);

        row.appendChild(sizeCell);
        
        tableBody.appendChild(row);// Добавляем строку в тело таблицы

        if (file.f_type == "Директория") { //если это директория
             // Добавляем эффект наведения и обработчик клика для директорий
            row.classList.add("itemdir")
            row.classList.add("directstyle")
            row.addEventListener('click', function() {
                directoryTraversal(file.f_name)// Обработчик клика для перехода в поддиректорию
            
           
            });
        }loader.classList.remove("loader")
        mask.classList.remove("mask") 
        

    });
}
// directoryTraversal - Функция для перехода в поддиректорию
function directoryTraversal(path){
    parDir = curDir; // Сохраняем текущий путь в переменную parDir
    changeCurDir(`${curDir}/${path}`); // Обновляем текущий путь, добавив указанный путь
    updateCurrentPath(); // Обновляем текущий путь в DOM
    fetchFiles()
}