
import { loader, mask } from "./script.js";
import { curDir, WriteRoot } from './buttons.js';
import { displayFiles } from "./filesFunctions.js";

//fetchFiles - Функция для получения файлов
export function fetchFiles(sort) {
    loader.classList.add('loader');
    mask.classList.add('mask');

    return fetch(`http://localhost:9999/fs?dst=${curDir}&sort=${sort}`) // Создаем URL для запроса к API
        .then(response => { // Проверяем, является ли ответ успешным (код состояния 200-299)
            if (!response.ok) {
                 // Если ответ не успешный, выбрасываем ошибку с информацией о статусе
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json(); // Если ответ успешный, преобразуем его в JSON
        })
        .then(files => { // Вызываем функцию fetchFiles с переданным параметром сортировки
            console.log("Received files:", files);

            if (files.f_status === 200) {
                displayFiles(files.f_data); // Отображаем отсортированные файлы в таблице
                WriteRoot(files.f_root);
            } else { //обрабатываем ошибки
                alert('Указан неверный путь: ' + files.f_error);
            }
        })
        .catch(error => { //обрабатываем ошибки
            console.error('Ошибка при получении данных:', error);
        });
}
