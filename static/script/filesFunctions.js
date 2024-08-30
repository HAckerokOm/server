
import {mask,loader} from "./script.js"
import {fetchFiles} from "./fetchh.js"
import {updateCurrentPath, curDir, changeCurDir}from './buttons.js'

export let parDir; // Переменная для сохранения родительской директории


// displayFiles - Функция отображения файлов в таблице
export function displayFiles(files) {
    const tableBody = document.querySelector('.item-list tbody');  // Находим тело таблицы с классом .item-list tbody
    tableBody.innerHTML = '';
    
    files.forEach(file => { // Проходимся по списку файлов
        const row = document.createElement('tr'); // Создаем новую строку таблицы
        row.classList.add(file.f_type === 'directory' ? 'itemdir' : 'itemfile'); // Добавляем класс itemdir, если это директория, иначе itemfile

        ['f_type', 'f_name'].forEach(key => { // Создаем ячейки для имени файла и типа файла
            const cell = document.createElement('td');
            cell.innerHTML = `<span class="${key}">${file[key]}</span>`;
            row.appendChild(cell);
        });

        // Создаем ячейку для размера файла
        const sizeCell = document.createElement('td');
        sizeCell.className = 'size-cell';

        // Создаем элементы для отображения размера и типа размера файла
        const sizeSpan = document.createElement('span');
        sizeSpan.textContent = `${file.f_size.toFixed(1)} `; // используем toFixed для округления после запятой
        sizeSpan.className = 'f_size';
        sizeCell.appendChild(sizeSpan);

        const typeSpan = document.createElement('span');
        typeSpan.textContent = `${file.f_typesize}`;
        typeSpan.className = 'f_typesize';
        sizeCell.appendChild(typeSpan);

        row.appendChild(sizeCell);
        
        tableBody.appendChild(row);

        // Добавляем обработчик клика для директорий
        if (file.f_type == "Директория") {
            row.classList.add("itemdir");
            row.classList.add("directstyle");
            row.addEventListener('click', function() {
                goSubDir(file.f_name); // Вызываем функцию перехода в поддиректорию при клике на директорию
            });
        }
         // Убираем маску и загрузочный индикатор после отображения файлов
        loader.classList.remove("loader");
        mask.classList.remove("mask"); 
    });
}
// goSubDir - Функция для перехода в поддиректорию
function goSubDir(path){
    parDir = curDir; // Сохраняем текущий путь в переменную parDir
    changeCurDir(`${curDir}/${path}`); // Обновляем текущий путь, добавив указанный путь
    updateCurrentPath(); // Обновляем текущий путь в DOM
    fetchFiles() // Вызываем fetchFiles для обновления списка файлов в новой директории
}