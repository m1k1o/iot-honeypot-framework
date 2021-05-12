module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  devServer: {
    disableHostCheck: true,
    proxy: process.env.API_SERVER ? {
      '^/api': {
        target: process.env.API_SERVER,
      },
    } : undefined,
  },
}
