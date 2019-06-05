module.exports = {
    publicPath: '/dynup',
    devServer: {
        proxy: {
            '^/dynup/api': {
                target: 'http://127.0.0.1:8090',
                changeOrigin: true
            },
        }
    }
}
