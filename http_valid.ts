import { enableProdMode } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { AppModule } from './app/app.module';
import { environment } from './environments/environment';
import { HttpClient } from '@angular/common/http';

if (environment.production) {
  enableProdMode();
}

// realiza uma solicitação HTTP para verificar se a URL está disponível
const http = new HttpClient(null);
http.get('http://exemplo.com').subscribe(
  response => {
    console.log('A URL está acessível');
    platformBrowserDynamic().bootstrapModule(AppModule)
      .catch(err => console.error(err));
  },
  error => {
    console.log('A URL não está acessível');
  }
);
