module.exports = {
  devServer: {
    host: '0.0.0.0',
    port: 80,
    proxy: {
      '/api': {
        target: 'http://api-gateway:8080',
        changeOrigin: true,
      }
    },
    watchFiles: {
      paths: ['src/**/*', 'public/**/*'],
      options: {
        usePolling: true
      }
    }
  },
  transpileDependencies: true
}
