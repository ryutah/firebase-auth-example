import webpack from "webpack"
import config from "./webpack.base.config.babel"

config.plugins = [
  new webpack.optimize.UglifyJsPlugin({
    // if need save license after ugilify
    // output: {
    //   comments: uglifySaveLicense,
    // },
    // if need save comment after ugilify
    // comments: true,
    compress: {
      warnings: false
    },
  }),
]

export default config
