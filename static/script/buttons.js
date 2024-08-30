
import {fetchFiles} from "./fetchh.js"

let rootDir; // Переменная для дефолтного пути

export let curDir = "/home/yaroslav/work"; // Определяем текущую директорию

const asc = "asc" // Константа для сортировки по возрастанию
const desc = "desc" // Константа для сортировки по убыванию

// WriteRoot - Функция для установки дефолтного пути
export function WriteRoot(defaultRoot) {
    // Записываем дефолтный путь в rootDir
    rootDir = defaultRoot;
    console.log(`Дефолтный путь установлен: ${rootDir}`);
}

// SortBut - Функция кнопки сортировки
export function SortBut() {
    // Добавляем обработчик клика на элемент с id="sortbut"
    document.getElementById("sortbut").addEventListener('click', SortPrompt);
}

//Функция для выбора сортировки файлов
function SortPrompt() {
    const sort = prompt('Введите параметр сортировки (asc или desc):', ''); // Показываем пользователю окно ввода для сортировки
    if (sort == asc || sort == desc ) { // Проверяем, не пустая ли строка ввода
        fetchFiles(sort)
    } else {
        alert("Вы ввели неправильное значение")
    }
}
// BackBut - функция кнопки назад
export function BackBut(){
    let backButton = document.getElementById("backbut") // Получаем элемент с id="backbut" и добавляем к нему обработчик клика
        backButton.addEventListener('click', returnToPreviousDirectory) // При нажатии кнопки вызывается функция returnToPreviousDirectory
}
// returnToPreviousDirectory - Функция возврата в предыдущую директорию
function returnToPreviousDirectory() {
    if(curDir === rootDir) { // Проверяем, равен ли текущий путь корневой директории
        
        alert("Вы в корневой дирректории")  // Если да, показываем сообщение и прекращаем выполнение функции
        return
    }
    curDir = curDir.split('/').slice(0, -1).join('/'); // Разделяем текущий путь на компоненты
    // Если результат пустой, устанавливаем корневую директорию
    if (curDir == ""){
        curDir = "/"
    }
    updateCurrentPath(); // Обновляем текущий путь в DOM
    fetchFiles() // Вызываем функцию fetchFiles для обновления списка файлов
}

// changeCurDir - Функция изменения текущей директории
export function changeCurDir(NewcurDir){
    curDir = NewcurDir
}
// updateCurrentPath - Функция для обновления текущего пути в DOM
export function updateCurrentPath() {
    // Получаем элемент с id="Path" и устанавливаем его текстовое содержимое
    document.getElementById("Path").textContent = `Путь: ${curDir}`;
}
