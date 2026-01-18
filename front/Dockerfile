FROM node:20-alpine AS build 
WORKDIR /app 

# Установка зависимостей 
COPY package.json yarn.lock ./ 
RUN yarn install --frozen-lockfile 
COPY . . 
RUN yarn build 

FROM nginx:alpine 
WORKDIR /usr/share/nginx/html 

# Удаление дефолтной страницы nginx 
RUN rm -rf ./* 


COPY --from=build /app/dist ./
COPY nginx.conf /etc/nginx/conf.d/default.conf 
COPY entrypoint.sh /entrypoint.sh 
RUN chmod +x /entrypoint.sh 

EXPOSE 80 443
CMD ["/entrypoint.sh"]