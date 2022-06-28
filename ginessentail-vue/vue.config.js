const webpack = require('webpack')
const Timestamp = new Date().getTime()
module.exports = {
  productionSourceMap: false,
  publicPath: './',
  configureWebpack: config => {
    if (process.env.NODE_ENV === 'production') {
      return {
        plugins: [
          new webpack.ProvidePlugin({
            $: 'jquery',
            jQuery: 'jquery',
            'windows.jQuery': 'jquery'
          })
        ],
        output: {
          filename: `js/[name].${process.env.NODE_ENV}.${Timestamp}.js`,
          chunkFilename: `js/[name].${process.env.NODE_ENV}.${Timestamp}.js`
        }
      }
    } else {
      return {
        plugins: [
          new webpack.ProvidePlugin({
            $: 'jquery',
            jQuery: 'jquery',
            'windows.jQuery': 'jquery'
          })
        ]
      }
    }
  },
  chainWebpack: config => {
    const svgRule = config.module.rule('svg')
    svgRule.uses.clear()
    svgRule
      .use('svg-sprite-loader')
      .loader('svg-sprite-loader')
      .options({
        symbolId: 'icon-[name]'
      })
  },
  devServer: {
    https: false
  }
}
