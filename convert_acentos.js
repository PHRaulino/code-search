function converterAcentos(str) {
  return str.normalize('NFC')
    .replace(/[áàãâä]/g, '\u00E1') // á
    .replace(/[éèêë]/g, '\u00E9') // é
    .replace(/[íìîï]/g, '\u00ED') // í
    .replace(/[óòõôö]/g, '\u00F3') // ó
    .replace(/[úùûü]/g, '\u00FA') // ú
    .replace(/[ç]/g, '\u00E7') // ç
    .replace(/[ñ]/g, '\u00F1'); // ñ
}

let minhaString = "Mamãe, água mole em pedra dura, tanto bate até que fura";
let minhaStringConvertida = converterAcentos(minhaString);
console.log(minhaStringConvertida); // Output: "Mama\u0303e, a\u0301gua mole em pedra dura, tanto bate ate\u0301 que fura"
