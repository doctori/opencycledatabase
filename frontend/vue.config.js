module.exports = {
  devServer: {
      public: "ocd.io",
      headers: { "Access-Control-Allow-Origin": "*" },
      proxy: {
        '^/api': {
          target: 'http://localhost:8081',
          pathRewrite: {'^/api':''}
        }
      }
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
