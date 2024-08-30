import {fetchFiles} from "./fetchh.js"
import * as buttons from './buttons.js'
export let mask;
export let loader;


document.addEventListener('DOMContentLoaded', function(){ // Добавляем слушатель события загрузки DOM

    loader = document.getElementById('loaderr')
    mask = document.getElementById('maskk')
    
    buttons.SortBut() // Вызываем функцию SortBut()
    buttons.updateCurrentPath(); // Вызываем функцию updateCurrentPath() для обновления текущего пути
    fetchFiles()
    buttons.BackBut() // Вызываем функцию BackBut() для добавления обработчика кнопки назад
});