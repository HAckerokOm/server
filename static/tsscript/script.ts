import { fetchFiles } from "./fetchh";
import * as buttons from './buttons';
import "./../style/style.css"

export let mask: HTMLElement;
export let loader: HTMLElement;

// Добавляем слушатель события загрузки DOM
document.addEventListener('DOMContentLoaded', () => {
    // Получаем элементы с id loaderr и maskk
    loader = document.getElementById('loaderr')!;
    mask = document.getElementById('maskk')!;

    // Вызываем функцию SortBut()
    buttons.SortBut();

    // buttons.StatsButton()
    // Вызываем функцию updateCurrentPath() для обновления текущего пути
    buttons.updateCurrentPath();

    // Вызываем fetchFiles() для получения файлов
    fetchFiles("");

    // Вызываем функцию BackBut() для добавления обработчика кнопки назад
    buttons.BackBut();
});
