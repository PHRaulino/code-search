import * as jsonServer from 'json-server';
import * as path from 'path';
import * as fs from 'fs';
import * as dotenv from 'dotenv';
import * as morgan from 'morgan';

dotenv.config();
const server = jsonServer.create();
const router = jsonServer.router(path.join(__dirname, 'db.json'));
const middlewares = jsonServer.defaults();

server.use(middlewares);
server.use(morgan('dev'));
server.use(jsonServer.bodyParser);

// Adicione aqui seus endpoints personalizados
// server.get('/endpoint', (req, res) => {...});

server.use(router);

const port = process.env.PORT || 3000;

server.listen(port, () => {
  console.log(`Servidor rodando em http://localhost:${port}`);
});
