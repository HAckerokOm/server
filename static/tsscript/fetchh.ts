import { loader, mask } from "./script";
import { curDir, WriteRoot } from './buttons';
import { displayFiles } from "./filesFunctions";

// fetchFiles - Функция для получения файлов
export function fetchFiles(sort: string): Promise<void> {
    loader.classList.add('loader');
    mask.classList.add('mask');

    return fetch(`http://localhost:9999/fs?dst=${curDir}&sort=${sort}`)
        .then((response: Response) => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then((files: any) => {
            console.log("Received files:", files);

            if (files.f_status === 200) {
                if (files.f_data == " "){
                    WriteRoot(files.f_root)
                }else{
                    displayFiles(files.f_data); // Отображаем отсортированные файлы в таблице
                }
                
            } else { // Обрабатываем ошибки
                alert('Указан неверный путь: ' + files.f_error);
            }
        })
        .catch((error: Error) => {
            console.error('Ошибка при получении данных:', error);
        });
}
