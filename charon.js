fetch('https://www.exemplo.com/')
  .then(response => response.text())
  .then(html => {
    const parser = new DOMParser();
    const htmlDocument = parser.parseFromString(html, 'text/html');
    const valor = htmlDocument.querySelector('#id-do-elemento').textContent;
    console.log(valor);
  })
  .catch(error => {
    console.error('Erro ao fazer a request:', error);
  });
