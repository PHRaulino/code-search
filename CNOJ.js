function gerarCnpj() {
  var n1 = Math.floor(Math.random() * 9) + 1; // O primeiro dígito não pode ser zero
  var n2 = Math.floor(Math.random() * 10);
  var n3 = Math.floor(Math.random() * 10);
  var n4 = Math.floor(Math.random() * 10);
  var n5 = Math.floor(Math.random() * 10);
  var n6 = Math.floor(Math.random() * 10);
  var n7 = Math.floor(Math.random() * 10);
  var n8 = Math.floor(Math.random() * 10);
  var n9 = 0;
  var n10 = 0;
  var n11 = 0;
  var n12 = 1;

  // Cálculo do primeiro dígito verificador
  var soma1 = n1 * 5 + n2 * 4 + n3 * 3 + n4 * 2 + n5 * 9 + n6 * 8 + n7 * 7 + n8 * 6 + n9 * 5 + n10 * 4 + n11 * 3 + n12 * 2;
  var digito1 = 11 - (soma1 % 11);
  if (digito1 > 9) {
    digito1 = 0;
  }

  // Cálculo do segundo dígito verificador
  var soma2 = n1 * 6 + n2 * 5 + n3 * 4 + n4 * 3 + n5 * 2 + n6 * 9 + n7 * 8 + n8 * 7 + n9 * 6 + n10 * 5 + n11 * 4 + n12 * 3 + digito1 * 2;
  var digito2 = 11 - (soma2 % 11);
  if (digito2 > 9) {
    digito2 = 0;
  }

  // Retorna o CNPJ formatado com pontos, barra e traço
  return '' + n1 + n2 + '.' + n3 + n4 + n5 + '.' + n6 + n7 + n8 + '/' + n9 + n10 + n11 + n12 + '-' + digito1 + digito2;
}

console.log(gerarCnpj());
