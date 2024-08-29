
import { fetchFiles } from "./filesFunctions.js";

export let curDir = "/home/yaroslav/work"; // Определяем текущую директорию

const asc = "asc"
const desc = "desc"

// ElementGet - Функция для обработки сортировки файлов
export function ElementGet(){
    document.getElementById("3").addEventListener('click', function() { // Добавляем обработчик клика на элемент с id="3"
        const sort = prompt('Введите параметр сортировки (asc или desc):', '');  // Показываем пользователю окно ввода для сортировки
        if (sort == asc || sort == desc ) { // Проверяем, не пустая ли строка ввода
            fetchFiles(sort)

        }else alert("вы ввели не правильное значение")
    })}

// BackBut - функция кнопки назад
export function BackBut(){
    let backButton = document.getElementById("2") // Получаем элемент с id="2" и добавляем к нему обработчик клика
        backButton.addEventListener('click', returnToPreviousDirectory) // При нажатии кнопки вызывается функция returnToPreviousDirectory
}

function returnToPreviousDirectory() {
    if(curDir === "/home") { // Проверяем, равен ли текущий путь корневой директории
        
        alert("Вы в корневой дирректории")  // Если да, показываем сообщение и прекращаем выполнение функции
        return
    }
    curDir = curDir.split('/').slice(0, -1).join('/'); // Разделяем текущий путь на компоненты
    updateCurrentPath(); // Обновляем текущий путь в DOM
    fetchFiles()
}

export function changeCurDir(NewcurDir){
    curDir = NewcurDir
}
// updateCurrentPath - Функция для обновления текущего пути в DOM
export function updateCurrentPath() {
    // Получаем элемент с id="4" и устанавливаем его текстовое содержимое
    document.getElementById("4").textContent = `Путь: ${curDir}`;
}