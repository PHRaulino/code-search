import { HttpClient } from '@angular/common/http';

async function checkUrl(httpClient: HttpClient, url: string): Promise<boolean> {
  try {
    // Faz uma solicitação HTTP GET para a URL
    await httpClient.get(url).toPromise();
    // Se a solicitação for bem-sucedida, retorna verdadeiro
    return true;
  } catch {
    // Se a solicitação falhar, retorna falso
    return false;
  }
}

async function bootstrap() {
  const httpClient = new HttpClient();
  const url = 'https://www.example.com';

  // Verifica se a URL funciona
  const isUrlWorking = await checkUrl(httpClient, url);

  console.log(`A URL ${url} está funcionando? ${isUrlWorking}`);
}

bootstrap();
