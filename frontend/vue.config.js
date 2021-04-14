module.exports = {
  devServer: {
      proxy: 'http://localhost:8081',
  },

  chainWebpack: config => {
    config
        .plugin('html')
        .tap(args => {
            args[0].title = "Open Cycle Database";
            return args;
        })
  },

  transpileDependencies: [
    'vuetify'
  ],

  pluginOptions: {
    i18n: {
      locale: 'en',
      fallbackLocale: 'fr',
      localeDir: 'locales',
      enableInSFC: true
    }
  }
}
