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
  ]
}
