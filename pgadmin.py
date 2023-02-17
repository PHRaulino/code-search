version: '3.8'

services:
  pgadmin:
    image: dpage/pgadmin4
    container_name: pgadmin
    restart: always
    environment:
      PGADMIN_DEFAULT_EMAIL: 'seu-email@exemplo.com'
      PGADMIN_DEFAULT_PASSWORD: ''
    ports:
      - "5050:80"
