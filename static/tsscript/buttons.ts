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
    const sort = prompt('Введите параметр сортировки (asc или desc):', ''); // Показываем пользователю окно ввода для сортировки
    if (sort === asc || sort === desc) { // Проверяем, не пустая ли строка ввода
        fetchFiles(sort)
    } else {
        alert("Вы ввели неправильное значение")
    }
}

// BackBut - функция кнопки назад
export function BackBut(){
    const backButton = document.getElementById("backbut")!; // Получаем элемент с id="backbut"
    backButton.addEventListener('click', returnToPreviousDirectory); // При нажатии кнопки вызывается функция returnToPreviousDirectory
}

// returnToPreviousDirectory - Функция возврата в предыдущую директорию
function returnToPreviousDirectory() {
    if(curDir === rootDir) { // Проверяем, равен ли текущий путь корневой директории
        
        alert("Вы в корневой дирректории"); // Если да, показываем сообщение и прекращаем выполнение функции
        return;
    }
    curDir = curDir.split('/').slice(0, -1).join('/'); // Разделяем текущий путь на компоненты
    // Если результат пустой, устанавливаем корневую директорию
    if (curDir === "") {
        curDir = "/";
    }
    updateCurrentPath(); // Обновляем текущий путь в DOM
    fetchFiles(""); // Вызываем функцию fetchFiles для обновления списка файлов
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

