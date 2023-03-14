import { enableProdMode } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { AppModule } from './app/app.module';
import { environment } from './environments/environment';
import { loadRemoteEntry } from './remote-entry';

if (environment.production) {
  enableProdMode();
}

const availableUrls = ['http://exemplo.com', 'http://outraurl.com', 'http://maisumaurl.com'];

const loadRemoteEntries = availableUrls.map(url => loadRemoteEntry(url));

Promise.allSettled(loadRemoteEntries).then(results => {
  const rejectedUrls = [];

  results.forEach((result, index) => {
    if (result.status === 'rejected') {
      const url = availableUrls[index];
      const reason = result.reason;
      rejectedUrls.push({ url, reason });
    }
  });

  if (rejectedUrls.length > 0) {
    console.error('Erro ao carregar as seguintes URLs:');
    rejectedUrls.forEach(rejectedUrl => {
      console.error(`${rejectedUrl.url}: ${rejectedUrl.reason}`);
    });
  } else {
    platformBrowserDynamic().bootstrapModule(AppModule)
      .catch(err => console.error(err));
  }
}).catch(error => {
  console.error(error);
});
