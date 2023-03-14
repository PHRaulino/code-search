.catch((reasons: any[]) => {
    const failedUrls = reasons.map((reason, index) => {
      return { url: urls[index], reason };
    });
    console.error('Erro ao carregar as seguintes URLs:');
    console.error(failedUrls);
  });
