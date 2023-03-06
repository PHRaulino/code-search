const fetch = require('node-fetch');

const url = 'https://google.com.br';

fetch(url, { agent: false, rejectUnauthorized: false })
  .then(response => {
    console.log(response.status);
    return response.text();
  })
  .then(text => {
    console.log(text);
  })
  .catch(error => {
    console.error(error);
  });
