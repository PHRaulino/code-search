const data = {mfe: {teste: {poc1: {url: 't1.com', name: 't1'}}, teste2:{poc2:{url: 't2.com', name: 't2'}}}};

const array = Object.values(data.mfe)
  .map((obj) => Object.values(obj))
  .flat()
  .filter((obj) => obj.hasOwnProperty('url') && obj.hasOwnProperty('name'));

console.log(array);
// Saída: [{url: 't1.com', name: 't1'}, {url: 't2.com', name: 't2'}]
