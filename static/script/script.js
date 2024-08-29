import{fetchFiles} from './filesFunctions.js'
import * as buttoms from './buttoms.js'
export let mask;
export let loader;


document.addEventListener('DOMContentLoaded', function(){ // Добавляем слушатель события загрузки DOM
    loader = document.getElementById('6')
    mask = document.getElementById('5')
    
    
    buttoms.Elget() // Вызываем функцию Elget()
    buttoms.updateCurrentPath(); // Вызываем функцию updateCurrentPath() для обновления текущего пути
    fetchFiles()
    buttoms.BackBut() // Вызываем функцию BackBut() для добавления обработчика кнопки назад
});


