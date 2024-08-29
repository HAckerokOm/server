import{fetchFiles} from './filesFunctions.js'
import * as buttons from './buttons.js'
export let mask;
export let loader;


document.addEventListener('DOMContentLoaded', function(){ // Добавляем слушатель события загрузки DOM
    loader = document.getElementById('6')
    mask = document.getElementById('5')
    
    buttons.ElementGet() // Вызываем функцию ElementGet()
    buttons.updateCurrentPath(); // Вызываем функцию updateCurrentPath() для обновления текущего пути
    fetchFiles()
    buttons.BackBut() // Вызываем функцию BackBut() для добавления обработчика кнопки назад
});


