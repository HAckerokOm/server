import { loader, mask } from "./script";
import { curDir, WriteRoot } from './buttons';
import { displayFiles } from "./filesFunctions";

// fetchFiles - Функция для получения файлов
export function fetchFiles(sort: string): Promise<void> {
    // Добавление классов для отображения индикатора загрузки и маски
    loader.classList.add('loader');
    mask.classList.add('mask');
    // Выполнение запроса к серверу
    return fetch(`/fs?dst=${curDir}&sort=${sort}`)
        .then((response: Response) => {
            // Проверка статуса ответа
            if (!response.ok) { // Проверка статуса ответа
                throw new Error(`HTTP-ошибка! статус: ${response.status}`);
            } 
            // Если ответ успешен, преобразование ответа в JSON
            return response.json();
        })
        .then((files: any) => { 
            // Проверка статуса ответа от сервера
            if (files.f_status === 200) {  
                if (files.f_data == " "){ 
                    WriteRoot(files.f_root)
                }else{
                    // Отображаем  файлы в таблице
                    displayFiles(files.f_data); 
                }
                
            } else { 
                // Обрабатываем ошибки
                alert('Указан неверный путь: ' + files.f_error);
            }
        })
        // Обработка любых ошибок, возникших во время выполнения запроса
        .catch((error: Error) => { 
            console.error('Ошибка при получении данных:', error);
        });
}
