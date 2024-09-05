import { fetchFiles } from "./fetchh";

// Переменная для дефолтного пути
let rootDir: string | null = null;

// Определяем текущую директорию
export let curDir: string = ""; // Корневая директория

const asc: string = "asc"; // Константа для сортировки по возрастанию
const desc: string = "desc"; // Константа для сортировки по убыванию

// WriteRoot - Функция для установки дефолтного пути
export function WriteRoot(defaultRoot: string) {
    // Записываем дефолтный путь в rootDir
    rootDir = defaultRoot;
    curDir = rootDir;
    fetchFiles("")
    console.log(`Дефолтный путь установлен: ${rootDir}`);
}

// SortBut - Функция кнопки сортировки
export function SortBut() {
    // Добавляем обработчик клика на элемент с id="sortbut"
    document.getElementById("sortbut")!.addEventListener('click', SortPrompt);
}

//Функция для выбора сортировки файлов
function SortPrompt() {
    // Показываем пользователю окно ввода для сортировки
    const sort = prompt('Введите параметр сортировки (asc или desc):', ''); 
    // Проверяем, не пустая ли строка ввода
    if (sort === asc || sort === desc) { 
        fetchFiles(sort)
    } else {
        alert("Вы ввели неправильное значение")
    }
}

// BackBut - функция кнопки назад
export function BackBut(){
    // Получаем элемент с id="backbut"
    const backButton = document.getElementById("backbut")!; 
    // При нажатии кнопки вызывается функция returnToPreviousDirectory
    backButton.addEventListener('click', returnToPreviousDirectory); 
}

// returnToPreviousDirectory - Функция возврата в предыдущую директорию
function returnToPreviousDirectory() {
    // Проверяем, равен ли текущий путь корневой директории
    if(curDir === rootDir) { 
        // Если да, показываем сообщение и прекращаем выполнение функции
        alert("Вы в корневой дирректории"); 
        return;
    }
    // Разделяем текущий путь на компоненты
    curDir = curDir.split('/').slice(0, -1).join('/'); 
    // Если результат пустой, устанавливаем корневую директорию
    if (curDir === "") {
        curDir = "/";
    }
    // Обновляем текущий путь в DOM
    updateCurrentPath(); 
    // Вызываем функцию fetchFiles для обновления списка файлов
    fetchFiles(""); 
}

// changeCurDir - Функция изменения текущей директории
export function changeCurDir(NewcurDir: string){
    curDir = NewcurDir;
}

// updateCurrentPath - Функция для обновления текущего пути в DOM
export function updateCurrentPath() {
    // Получаем элемент с id="Path" и устанавливаем его текстовое содержимое
    document.getElementById("Path")!.textContent = `Путь: ${curDir}`;
}

