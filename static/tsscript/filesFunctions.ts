import { mask, loader } from "./script";
import { fetchFiles } from "./fetchh";
import { updateCurrentPath, curDir, changeCurDir } from './buttons';


export let parDir: string | null = null; // Переменная для сохранения родительской директории

// displayFiles - Функция отображения файлов в таблице
export function displayFiles(files: any[]): void {
    const tableBody = document.querySelector('.item-list tbody'); // Находим тело таблицы с классом .item-list tbody
    
    tableBody!.innerHTML = '';
    
    // Проходимся по списку файлов
    files.forEach(file => { 
        // Создаем новую строку таблицы
        const row = document.createElement('tr'); 
        // Добавляем класс itemdir, если это директория, иначе itemfile
        row.classList.add(file.f_type === 'directory' ? 'itemdir' : 'itemfile'); 

        // Создаем ячейки для имени файла и типа файла
        ['f_type', 'f_name'].forEach(key => { 
            const cell = document.createElement('td');
            cell.innerHTML = `<span class="${key}">${file[key]}</span>`;
            row.appendChild(cell);
        });

        // Создаем ячейку для размера файла
        const sizeCell = document.createElement('td');
        sizeCell.className = 'size-cell';

        // Создаем элементы для отображения размера и типа размера файла
        const sizeSpan = document.createElement('span');
        // используем toFixed для округления после запятой
        sizeSpan.textContent = `${file.f_size.toFixed(1)} `; 
        sizeSpan.className = 'f_size';
        sizeCell.appendChild(sizeSpan);

        const typeSpan = document.createElement('span');
        typeSpan.textContent = `${file.f_typesize}`;
        typeSpan.className = 'f_typesize';
        sizeCell.appendChild(typeSpan);

        row.appendChild(sizeCell);
        
        tableBody!.appendChild(row);

        // Добавляем обработчик клика для директорий
        if (file.f_type === "Директория") {
            row.classList.add("itemdir");
            row.classList.add("directstyle");
            // Вызываем функцию перехода в поддиректорию при клике на директорию
            row.addEventListener('click', function() { 
                goSubDir(file.f_name); 
            });
        }
         // Убираем маску и загрузочный индикатор после отображения файлов
        loader.classList.remove("loader");
        mask.classList.remove("mask"); 
    });
}

// goSubDir - Функция для перехода в поддиректорию
function goSubDir(path: string): void {
    // Сохраняем текущий путь в переменную parDir
    parDir = curDir; 
    if (curDir === "/"){
        changeCurDir(`/${path}`)
    }else{
        // Обновляем текущий путь, добавив указанный путь
        changeCurDir(`${curDir}/${path}`); 
    }
    // Обновляем текущий путь в DOM
    updateCurrentPath(); 
    // Вызываем fetchFiles для обновления списка файлов в новой директории
    fetchFiles(""); 
}
